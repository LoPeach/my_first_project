package collector

import (
	"context"
	"fmt"
	"math/rand/v2"
	"my_crypto_project/internal/models"
	"sync"
	"time"
)

func CheckPrice(ctx context.Context, exc string, fileChan chan<- models.PriceResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Работа для %s отменена!\n", exc)
			return
		default:
		}
		n := rand.Float64() * 60000
		newPrice := models.PriceResult{Exchange: exc, Price: n, Timestamp: time.Now()}
		fileChan <- newPrice
		time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
	}

}
