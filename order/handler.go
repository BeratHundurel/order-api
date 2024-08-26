package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Repo interface {
	Insert(ctx context.Context, order Order) error
	FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
	GetByID(ctx context.Context, id uint64) (Order, error)
	Update(ctx context.Context, order Order) error
	DeleteByID(ctx context.Context, id uint64) error
}

type OrderHandler struct {
	Repo Repo
}

func (o *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID  `json:"customer_id"`
		LineItems  []LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()

	order := Order{
		ID:         rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}
	order.Total = CalculateTotal(&order)

	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(r.Context(), FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []Order `json:"items"`
		Next  uint64  `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (o *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, shouldReturn := ParseParamToUint(r, w)
	if shouldReturn {
		return
	}

	data, err := o.Repo.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (o *OrderHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, shouldReturn := ParseParamToUint(r, w)
	if shouldReturn {
		return
	}

	data, err := o.Repo.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch body.Status {
	case "shipped":
		if data.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shippedAt := time.Now().UTC()
		data.ShippedAt = &shippedAt
	case "completed":
		if data.CompletedAt != nil || data.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		completedAt := time.Now().UTC()
		data.CompletedAt = &completedAt
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.Update(r.Context(), data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func (o *OrderHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, shouldReturn := ParseParamToUint(r, w)
	if shouldReturn {
		return
	}

	err := o.Repo.DeleteByID(r.Context(), id)
	if errors.Is(err, ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
