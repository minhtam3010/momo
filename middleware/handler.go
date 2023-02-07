package middleware

import (
	"net/http"
	"strconv"

	"github.com/minhtam3010/momo/db"
)

type Handler struct {
	C *db.ConnectDB
}

func NewHandler() *Handler {
	return &Handler{C: db.NewDB()}
}

func (h *Handler) CreateTransactionMomo(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("orderId")
	requestId := r.URL.Query().Get("requestId")
	transID := r.URL.Query().Get("transId")
	amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Create transaction momo failed cuz internal code has some problem about logic conversion"}`))
	}

	payType := r.URL.Query().Get("payType")

	data := db.MomoPayment{
		OderID:    orderId,
		RequestID: requestId,
		TransID:   transID,
		Amount:    amount,
		PayType:   payType,
	}
	err = h.C.CreateMomoPayment(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Create transaction momo failed"}`))
		panic(err)
	}
	http.Redirect(w, r, "http://localhost:5500/thanks.html", http.StatusSeeOther)
}

func (h *Handler) Commit(w http.ResponseWriter, r *http.Request) {
	err := h.C.Commit()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Commit transaction momo failed"}`))
		panic(err)
	}
}
