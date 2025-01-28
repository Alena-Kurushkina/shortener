package repository

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func generateRandomString(length int) string {
	charset := []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := strings.Builder{}
	result.Grow(length)
	for i:=0; i< length; i++ {
		result.WriteByte(charset[random.Intn(len(charset))])
	}

	return result.String()
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