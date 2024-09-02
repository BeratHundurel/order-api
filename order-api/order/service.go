package order

import (
	"context"
	"net/http"
	"strconv"

	pb "github.com/BeratHundurel/order-api/currency-api/proto"
	"github.com/go-chi/chi"
)

func (o *OrderHandler) CalculateTotal(order *Order) (float32, error) {
    var total float32
    currencyClient := pb.NewCurrencyConverterClient(o.CurrencyConn)

    for _, item := range order.LineItems {
        req := &pb.ConvertCurrencyRequest{
            Price:        item.Price,
            FromCurrency: order.Currency,
            ToCurrency:   "USD",
        }

        res, err := currencyClient.ConvertCurrency(context.Background(), req)
        if err != nil {
            return 0, err
        }

        convertedPrice := res.ConvertedPrice
        total += float32(item.Quantity) * convertedPrice
    }

    order.Total = total
    return total, nil
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
