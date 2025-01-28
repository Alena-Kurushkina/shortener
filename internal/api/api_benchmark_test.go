package api

import (
	"strconv"
	"testing"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
)

func waitEmpty(sh *Shortener){
	for {
		if len(sh.deleteChan)==0{
			sh.done<-struct{}{}
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

	for i:=0; i<b.N; i++{
		b.StopTimer()
		for k:=1; k<1000; k++ {
			sh.deleteChan <- DeleteItem{IDs: []string{strconv.Itoa(k)}, UserID: uuid.NewV4()}
		}
		go waitEmpty(sh)
		b.StartTimer()

		sh.flushDeleteItems()
	}	
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for i:=0; i<b.N; i++{
		generateRandomString(15)
	}
}