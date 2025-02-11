package api

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/config"
)

func waitEmpty(sh *Shortener) {
	for {
		if len(sh.deleteChan) == 0 {
			sh.done <- struct{}{}
			return
		}
	}
}

func BenchmarkFlushDeleteItems(b *testing.B) {
	ctrl := gomock.NewController(b)
	m := NewMockStorager(ctrl)
	//m.EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(nil)

	cfg = config.InitConfig()
	sh := newShortenerObject(m, cfg)

	b.ResetTimer()

	b.Run("flush delete items", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for k := 1; k < 1000; k++ {
				sh.deleteChan <- DeleteItem{IDs: []string{strconv.Itoa(k)}, UserID: uuid.NewV4()}
			}
			go waitEmpty(sh)
			b.StartTimer()

			sh.flushDeleteItems()
		}
	})

	b.Run("flush delete items v2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for k := 1; k < 1000; k++ {
				sh.deleteChan <- DeleteItem{IDs: []string{strconv.Itoa(k)}, UserID: uuid.NewV4()}
			}
			go waitEmpty(sh)
			b.StartTimer()

			sh.flushDeleteItemsV2()
		}
	})

}

func BenchmarkGenerateRandomString(b *testing.B) {
	b.Run("generate random string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			generateRandomString(15)
		}
	})
	b.Run("generate random string faster", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			generateRandomStringFaster(15)
		}
	})
}
