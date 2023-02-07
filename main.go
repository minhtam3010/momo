package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//run for download library `go get github.com/sony/sonyflake`
	"github.com/minhtam3010/momo/middleware"
	"github.com/sony/sonyflake"
)

//define a payload, reference in https://developers.momo.vn/#cong-thanh-toan-momo-phuong-thuc-thanh-toan
type Payload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}

//define a POS payload, reference in https://developers.momo.vn/#thanh-toan-pos-xu-ly-thanh-toan
type PosHash struct {
	PartnerCode  string `json:"partnerCode"`
	PartnerRefID string `json:"partnerRefId"`
	Amount       int    `json:"amount"`
	PaymentCode  string `json:"paymentCode"`
}

type PosPayload struct {
	PartnerCode  string `json:"partnerCode"`
	PartnerRefID string `json:"partnerRefId"`
	Hash         string `json:"hash"`
	Version      int    `json:"version"`
}

var pubKeyData = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkpa+qMXS6O11x7jBGo9W
3yxeHEsAdyDE40UoXhoQf9K6attSIclTZMEGfq6gmJm2BogVJtPkjvri5/j9mBnt
A8qKMzzanSQaBEbr8FyByHnf226dsLt1RbJSMLjCd3UC1n0Yq8KKvfHhvmvVbGcW
fpgfo7iQTVmL0r1eQxzgnSq31EL1yYNMuaZjpHmQuT24Hmxl9W9enRtJyVTUhwKh
tjOSOsR03sMnsckpFT9pn1/V9BE2Kf3rFGqc6JukXkqK6ZW9mtmGLSq3K+JRRq2w
8PVmcbcvTr/adW4EL2yc1qk9Ec4HtiDhtSYd6/ov8xLVkKAQjLVt7Ex3/agRPfPr
NwIDAQAB
-----END PUBLIC KEY-----
`)

func ProcessMoMo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello Momo!\n")

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	//randome orderID and requestID
	a, err := flake.NextID()
	b, err := flake.NextID()

	// QR Pay
	var endpoint = "https://test-payment.momo.vn/v2/gateway/api/create"
	var partnerCode = "MOMO"
	var accessKey = "F8BBA842ECF85"
	var secretKey = "K951B6PE1waDMi640xX08PD3vg6EkVlz"
	var orderInfo = "pay with MoMo"
	var redirectUrl = "http://localhost:5500/thanks"
	var ipnUrl = "http://localhost:5500/thanks"
	var amount = "50000"
	var requestType = "captureWallet"
	var extraData = "" // pass empty value or Encode base64 JsonString
	var partnerName = "MoMo Payment"

	/* Pay with all Methods*/
	// var endpoint = "https://test-payment.momo.vn/v2/gateway/api/create"
	// var accessKey = "F8BBA842ECF85"
	// var secretKey = "K951B6PE1waDMi640xX08PD3vg6EkVlz"
	// var orderInfo = "pay with MoMo"
	// var partnerCode = "MOMO"
	// var redirectUrl = "https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b"
	// var ipnUrl = "https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b"
	// var amount = "50000"
	var orderId = strconv.FormatUint(a, 16)
	var requestId = strconv.FormatUint(b, 16)
	// var extraData = ""
	// var partnerName = "MoMo Payment"
	var storeId = "Test Store"
	var orderGroupId = ""
	var autoCapture = true
	var lang = "vi"
	// var requestType = "payWithMethod"

	//build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(secretKey))

	// Write Data to it
	hmac.Write(rawSignature.Bytes())
	// fmt.Println("Raw signature: " + rawSignature.String())

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))

	var payload = Payload{
		PartnerCode:  partnerCode,
		AccessKey:    accessKey,
		RequestID:    requestId,
		Amount:       amount,
		RequestType:  requestType,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		OrderID:      orderId,
		StoreId:      storeId,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupId,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		ExtraData:    extraData,
		Signature:    signature,
	}

	var jsonPayload []byte
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Payload: " + string(jsonPayload))
	// fmt.Println("Signature: " + signature)

	//send HTTP to momo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalln("This problem", err)
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Response from Momo: ", result)

	w.Header().Set("Content-Type", "application/json")

	http.Redirect(w, r, fmt.Sprintf("%s", result["payUrl"]), http.StatusSeeOther)
}

func main() {
	fmt.Println("Server is running on port 5500")
	mux := http.NewServeMux()
	handler := middleware.NewHandler()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/buy", ProcessMoMo)
	mux.HandleFunc("/thanks", handler.CreateTransactionMomo)
	mux.HandleFunc("/commit", handler.Commit)

	http.ListenAndServe(":5500", mux)
}

// accessKey=F8BBA842ECF85&amount=50000&extraData=&ipnUrl=https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b&orderId=632c4d48100020b&orderInfo=pay with MoMo&partnerCode=MOMO&redirectUrl=https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b&requestId=632c4d48101020b&requestType=payWithMethod
// accessKey=SvDmj2cOTYZmQQ3H&amount=10000&extraData=&ipnUrl=https://webhook.site/3c5b6488-a159-4f8d-b038-29eed82fab1e&orderId=632c4e4ea00020b&orderInfo=momo all-in-one&partnerCode=MOMOIQA420180417&redirectUrl=https://webhook.site/3c5b6488-a159-4f8d-b038-29eed82fab1e&requestId=632c4e4ea01020b&requestType=captureWallet
