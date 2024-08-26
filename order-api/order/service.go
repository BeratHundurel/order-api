package order

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func CalculateTotal(o *Order) float64 {
	var total float64
	for _, item := range o.LineItems {
		total += float64(item.Quantity) * item.Price
	}
	o.Total = total
	return total
}

func ParseParamToUint(r *http.Request, w http.ResponseWriter) (uint64, bool) {
	idParam := chi.URLParam(r, "id")
	const base = 10
	const bitSize = 64
	id, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 0, true
	}
	return id, false
}
