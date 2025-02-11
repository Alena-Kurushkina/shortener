package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

// A DBRepository store data in database.
type DBRepository struct {
	database      *sql.DB
	selectStmt    *sql.Stmt
	selectAllStmt *sql.Stmt
}

// GetDB creates DBRepository object in first call, then returns it with no recreation.
var GetDB func() (api.Storager, error)

// newDBRepository initializes data storage in database.
func newDBRepository(ctx context.Context, connectionStr string) (dbRep api.Storager, err error) {

	dbInit := func() (api.Storager, error) {
		dbRep := &DBRepository{}

		db, err := sql.Open("pgx", connectionStr)
		if err != nil {
			return nil, err
		}
		logger.Log.Info("DB connection opened")

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer func(){
			if tErr:=tx.Rollback(); tErr!=nil{
				err=tErr
			}
		}()

		tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS shortening(
				id varchar(50) PRIMARY KEY default substring(md5(random()::text),0,20),
				originalURL varchar(500) NOT NULL,
				shortURL varchar(250) NOT NULL,
				userUUID uuid,
				is_deleted bool DEFAULT(false),
				UNIQUE (originalURL)
			);
			CREATE UNIQUE INDEX IF NOT EXISTS short_idx on shortening (shortURL);
		`)

		err = tx.Commit()
		if err != nil {
			return nil, err
		}

		getDeletedFieldQuery, err := db.PrepareContext(ctx, `
			SELECT originalURL, is_deleted 
			FROM shortening 
			WHERE shortURL LIKE $1
		`)
		if err != nil {
			return nil, err
		}

		getShorteningQuery, err := db.PrepareContext(ctx, `
			SELECT originalURL, shortURL
			FROM shortening 
			WHERE shortening.useruuid = $1
		`)
		if err != nil {
			return nil, err
		}

		dbRep.database = db
		dbRep.selectStmt = getDeletedFieldQuery
		dbRep.selectAllStmt = getShorteningQuery

		return dbRep, err
	}

	GetDB = sync.OnceValues(dbInit)

	return GetDB()
}

// Insert saves short URL and original one to storage by user id.
// It returns AlreadyExistError if short URL is already in storage.
func (r DBRepository) Insert(ctx context.Context, userID uuid.UUID, insertedShortURL, insertedOriginalURL string) (err error) {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if tErr:=tx.Rollback(); tErr!=nil{
			err=tErr
		}
	}()

	sqlRow := tx.QueryRowContext(ctx,
		`INSERT INTO shortening (userUUID, originalURL, shortURL) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;`,
		userID,
		insertedOriginalURL,
		insertedShortURL,
	)

	var (
		dbOriginalURL string
		dbShortURL    string
	)
	err = sqlRow.Scan(&dbOriginalURL, &dbShortURL)
	if err != nil {
		return err
	}

	if dbShortURL != insertedShortURL {
		tx.Rollback()
		return sherr.NewAlreadyExistError(insertedOriginalURL, dbShortURL)
	}

	err=tx.Commit()

	return err
}

// DeleteRecords deletes records by their ids from storage.
// It is getting array of DeleteItem on input.
func (r DBRepository) DeleteRecords(ctx context.Context, deleteItems []api.DeleteItem) error {
	param := ""
	for _, v := range deleteItems {
		for _, i := range v.IDs {
			param = param + `(uuid('` + v.UserID.String() + `'), '` + i + `'),`
		}
	}
	param, _ = strings.CutSuffix(param, ",")

	stmt := `UPDATE shortening
		SET is_deleted=true
		FROM (
			VALUES` + param +
		`) AS data(id_user, shortening)
		WHERE shortening.useruuid=data.id_user
			AND shortening.shorturl=data.shortening`

	logger.Log.Info(stmt)

	sqlRes, err := r.database.ExecContext(ctx, stmt)

	if err != nil {
		logger.Log.Errorf("Error while deletion: ", err.Error())
	}

	sf, _ := sqlRes.RowsAffected()
	logger.Log.Infof("Rows affected while deletion: %s", strconv.FormatInt(sf, 10))

	return err
}

// InsertBatch saves array of BatchElement to storage.
func (r DBRepository) InsertBatch(ctx context.Context, userID uuid.UUID, batch []api.BatchElement) (err error) {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(){
		if tErr:=tx.Rollback(); tErr!=nil{
			err=tErr
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO shortening (id, userUUID, originalURL, shortURL) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;
	`)
	if err != nil {
		return err
	}

	for _, v := range batch {
		_, err = stmt.ExecContext(ctx,
			v.CorrelarionID,
			userID,
			v.OriginalURL,
			v.ShortURL,
		)
		if err != nil {
			return err
		}
	}

	err=tx.Commit()

	return err
}

// Close closes all statements and database.
func (r *DBRepository) Close() {
	r.selectStmt.Close()
	r.selectAllStmt.Close()
	r.database.Close()
}

// Ping verifies database connection.
func (r *DBRepository) Ping(ctx context.Context) error {
	return r.database.PingContext(ctx)
}

// Select returns longURL from storage by it shortening.
func (r DBRepository) Select(ctx context.Context, key string) (string, error) {
	row := r.selectStmt.QueryRowContext(ctx,
		key,
	)
	var (
		longURL string
		deleted bool
	)

	err := row.Scan(&longURL, &deleted)
	if err != nil {
		return "", err
	}
	if deleted {
		return "", sherr.ErrDBRecordDeleted
	}

	return longURL, nil
}

// SelectUserAll returns all user's pairs of long URL and shortening from storage.
func (r DBRepository) SelectUserAll(ctx context.Context, id uuid.UUID) (records []api.BatchElement, err error) {
	rows, err := r.selectAllStmt.QueryContext(ctx,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer func (){
		if tErr:=rows.Close(); tErr!=nil{
			err=tErr
		}
	}()

	records = make([]api.BatchElement, 0, 10)

	// пробегаем по всем записям
	for rows.Next() {
		var v api.BatchElement
		err = rows.Scan(&v.OriginalURL, &v.ShortURL)
		if err != nil {
			return nil, err
		}

		records = append(records, v)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, err
}