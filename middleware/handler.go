package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/minhtam3010/momo/db"
)

type UserInfo struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type Handler struct {
	C *db.ConnectDB
}

func NewHandler() *Handler {
	return &Handler{
		C: db.NewDB(),
	}
}

func decodeFromBase64(v interface{}, enc string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(enc))).Decode(v)
}

func (h *Handler) CreateTransactionMomo(w http.ResponseWriter, r *http.Request) {
	tx := db.NewTx(h.C)

	values := r.URL.Query()
	for k, v := range values {
		fmt.Println(k, " => ", v)
	}

	extraData := r.URL.Query().Get("extraData")
	var userInfo UserInfo
	err := decodeFromBase64(&userInfo, extraData)
	if err != nil {
		return
	}
	log.Println(userInfo)

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
	err = tx.CreateMomoPayment(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Create transaction momo failed"}`))
		panic(err)
	}
	http.Redirect(w, r, "http://localhost:5500/thanks.html", http.StatusSeeOther)
}
