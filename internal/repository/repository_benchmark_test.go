package repository

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}

func BenchmarkInsertBatch(b *testing.B) {
	cfg := config.InitConfig()

	db, err:=newDBRepository(context.TODO(),cfg.ConnectionStr)
	assert.NoError(b,err)

	uid:=uuid.NewV4()

	batch := make([]api.BatchElement, 0, 1000)
	for k:=range batch {
		batch[k].CorrelarionID=strconv.Itoa(k)
		batch[k].OriginalURL="http://testurl"+strconv.Itoa(k)
		batch[k].ShortURL=generateRandomString(15)
	}

	b.ResetTimer()

	for i:=0; i<b.N; i++{
		db.InsertBatch(context.TODO(), uid, batch)
	}
}

func BenchmarkSelectUserAll(b *testing.B) {
	cfg := config.InitConfig()

	db, err:=newDBRepository(context.TODO(),cfg.ConnectionStr)
	assert.NoError(b,err)

	uid:=uuid.NewV4()

	b.ResetTimer()

	for i:=0; i<b.N; i++{
		db.SelectUserAll(context.TODO(), uid)
	}
}