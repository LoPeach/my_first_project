package handlers

import (
	"encoding/json"
	"my_crypto_project/internal/models"
	"my_crypto_project/internal/storage"
	"net/http"
)

type PriceHandler struct {
	st *storage.Storage
}

// New создает новый экземпляр хендлера с доступом к хранилищу
func New(st *storage.Storage) *PriceHandler {
	return &PriceHandler{st: st}
}

func (h *PriceHandler) GetPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	excha := r.URL.Query().Get("exchange")

	// Вызываем метод, который мы уже написали в storage.go
	prices, err := h.st.GetPrices(ctx, "30 seconds")
	if err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}

	var priceResult []models.PriceResult
	var avg float64

	for _, p := range prices {
		if excha != "" && p.Exchange != excha {
			continue
		}
		priceResult = append(priceResult, p)
		avg += p.Price
	}

	if len(priceResult) > 0 {
		avg = avg / float64(len(priceResult))
	}

	finalData := models.PriceResultAVG{
		PriceResult: priceResult,
		AvgPrice:    avg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalData)
}

func (h *PriceHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := h.st.GetStats(ctx)
	if err != nil {
		http.Error(w, "Ошибка при получении стратистики", http.StatusInternalServerError)
		return
	}
	response := []models.DataForStats{stats}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
