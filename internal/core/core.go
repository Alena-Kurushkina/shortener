package core

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/generator"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

// Storager defines operations with data storage.
type Storager interface {
	Insert(ctx context.Context, userID uuid.UUID, key, value string) error
	InsertBatch(_ context.Context, userID uuid.UUID, batch []BatchElement) error
	Select(ctx context.Context, key string) (string, error)
	SelectUserAll(ctx context.Context, userID uuid.UUID) ([]BatchElement, error)
	DeleteRecords(ctx context.Context, deleteItems []DeleteItem) error
	Ping(ctx context.Context) error
	SelectStats(ctx context.Context) (stats *Stats, err error)
	Close()
}

// A ShortenerCore aggregates references to storage, config, delete and done channels.
type ShortenerCore struct {
	repo       Storager
	config     *config.Config
	deleteChan chan DeleteItem
	done       chan struct{}
}

func newShortenerCore(storage Storager, cfg *config.Config) *ShortenerCore {
	return &ShortenerCore{
		repo:       storage,
		config:     cfg,
		deleteChan: make(chan DeleteItem, 1024),
		done:       make(chan struct{}),
	}
}

// NewShortenerCore returns new Shortener pointer initialized by repository and config.
func NewShortenerCore(storage Storager, cfg *config.Config) *ShortenerCore {
	shortener := newShortenerCore(storage, cfg)

	go shortener.flushDeleteItems()

	return shortener
}

// DeleteItem represents pair of ids which identify unique record to delete.
type DeleteItem struct {
	IDs    []string
	UserID uuid.UUID
}

// Stats represents counts of stored URLs and users.
type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// A BatchElement represent structure to marshal element of request`s json array.
type BatchElement struct {
	CorrelarionID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url,omitempty"`
}

func (sh *ShortenerCore) flushDeleteItems() {
	ticker := time.NewTicker(1 * time.Second)

	items := make([]DeleteItem, 0, 1024)

	for {
		select {
		case msg := <-sh.deleteChan:
			items = append(items, msg)
		case <-ticker.C:
			sh.deleteRecords(items)

			items = make([]DeleteItem, 0, 1024)
		case <-sh.done:
			sh.deleteRecords(items)
			return
		}
	}
}

func (sh *ShortenerCore) deleteRecords(items []DeleteItem) {
	if len(items) == 0 {
		return
	}
	err := sh.repo.DeleteRecords(context.TODO(), items)
	if err != nil {
		logger.Log.Infof("Can't delete records", err.Error())
		return
	}
	logger.Log.Info("Patch of shortenings was deleted, patch length: " + strconv.Itoa(len(items)))
}

// Shutdown finishes work gracefully.
func (sh *ShortenerCore) Shutdown() {
	close(sh.done)
	sh.repo.Close()
	logger.Log.Info("Shortener core shutdown finished ")
}

// CreateShortening creates short string and saves it to storate.
func (sh *ShortenerCore) CreateShortening(ctx context.Context, userID uuid.UUID, longURL string) (shortURL string, err error) {
	// generate shortening
	shortStr := generator.GenerateRandomString(15)
	err = sh.repo.Insert(ctx, userID, shortStr, longURL)

	var existError *sherr.AlreadyExistError
	if errors.As(err, &existError) {
		shortURL = sh.config.BaseURL + existError.ExistShortStr
		return
	}

	shortURL = sh.config.BaseURL + shortStr
	return
}

// CreateShorteningBatch creates batch of short strings and saves it.
func (sh *ShortenerCore) CreateShorteningBatch(ctx context.Context, userID uuid.UUID, batch []BatchElement) ([]BatchElement, error) {
	// generate shortening
	for k := range batch {
		batch[k].ShortURL = generator.GenerateRandomString(15)
	}

	// write to data storage
	if err := sh.repo.InsertBatch(ctx, userID, batch); err != nil {
		return batch, err
	}
	// make response
	for k, v := range batch {
		batch[k].ShortURL = sh.config.BaseURL + v.ShortURL
	}
	return batch, nil
}

// GetFullString get full string from repository.
func (sh *ShortenerCore) GetFullString(ctx context.Context, shortURL string) (output string, err error) {
	// get long URL from repository
	output, err = sh.repo.Select(ctx, shortURL)
	return
}

// RegisterToDelete saves record id to delete channel.
func (sh *ShortenerCore) RegisterToDelete(ctx context.Context, recordIDs []string, userID uuid.UUID) {
	sh.deleteChan <- DeleteItem{IDs: recordIDs, UserID: userID}
}

// PingDB checks db connection.
func (sh *ShortenerCore) PingDB(ctx context.Context) error {
	return sh.repo.Ping(ctx)
}

// IsTrustedSubnet checks if ip is in trusted subnet.
func (sh *ShortenerCore) IsTrustedSubnet(ipStr string) (bool, error) {
	_, ipNetTr, err := net.ParseCIDR(sh.config.TrustedSubnet)
	if err != nil {
		return false, err
	}
	subnet := net.ParseIP(ipStr).Mask(ipNetTr.Mask)
	if !subnet.Equal(ipNetTr.IP) || sh.config.TrustedSubnet == "" {
		return false, nil
	}
	return true, nil
}

// GetStats gets static info from storage.
func (sh *ShortenerCore) GetStats(ctx context.Context) (*Stats, error) {
	return sh.repo.SelectStats(ctx)
}

// GetAllUserShortenings returns all user shortenings recods by user id.
func (sh *ShortenerCore) GetAllUserShortenings(ctx context.Context, userID uuid.UUID) ([]BatchElement, error) {
	allRecords, err := sh.repo.SelectUserAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(allRecords) == 0 {
		return nil, sherr.ErrNoShortenings
	}

	for k, v := range allRecords {
		allRecords[k].ShortURL = sh.config.BaseURL + v.ShortURL
	}
	return allRecords, nil
}
