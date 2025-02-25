package api

import (
	"testing"

	"github.com/Alena-Kurushkina/shortener/internal/generator"
)

// func waitEmpty(sh *Shortener) {
// 	for {
// 		if len(sh.deleteChan) == 0 {
// 			close(sh.done)
// 			return
// 		}
// 	}
// }

// func BenchmarkFlushDeleteItems(b *testing.B) {
// 	ctrl := gomock.NewController(b)
// 	m := NewMockStorager(ctrl)
// 	m.EXPECT().DeleteRecords(gomock.Any(), gomock.Any()).Return(nil)

// 	cfg = config.InitConfig()

// 	b.ResetTimer()

// 	b.Run("flush delete items", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			b.StopTimer()
// 			sh := newShortenerObject(m, cfg)
// 			for k := 1; k < 10; k++ {
// 				sh.deleteChan <- DeleteItem{IDs: []string{strconv.Itoa(k)}, UserID: uuid.NewV4()}
// 			}
// 			go waitEmpty(sh)
// 			b.StartTimer()

// 			sh.flushDeleteItems()
// 		}
// 	})

// 	b.Run("flush delete items v2", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			b.StopTimer()
// 			sh := newShortenerObject(m, cfg)
// 			for k := 1; k < 10; k++ {
// 				sh.deleteChan <- DeleteItem{IDs: []string{strconv.Itoa(k)}, UserID: uuid.NewV4()}
// 			}
// 			go waitEmpty(sh)
// 			b.StartTimer()

// 			sh.flushDeleteItems()
// 		}
// 	})

// }

func BenchmarkGenerateRandomString(b *testing.B) {
	b.Run("generate random string faster", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			generator.GenerateRandomString(15)
		}
	})
}
