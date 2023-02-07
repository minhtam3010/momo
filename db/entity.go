package db

import "time"

type MomoPayment struct {
	OderID      string `json:"orderId"`
	RequestID   string `json:"requestId"`
	TransID     string `json:"transId"`
	Amount      int    `json:"amount"`
	PayType     string `json:"payType"`
	DateCreated int64  `json:"dateCreated"`
	DateUpdated int64  `json:"dateUpdated"`
}

func (c *ConnectDB) CreateMomoPayment(data MomoPayment) (err error) {
	date := ConvertUnixDateToString(time.Now().Unix())
	_, err = c.tx.Exec("INSERT INTO momo_payment (order_id, request_id, trans_id, amount, pay_type, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?, ?)", data.OderID, data.RequestID, data.TransID, data.Amount, data.PayType, date, nil)
	if err != nil {
		return c.tx.Rollback()
	}
	return nil
}

func (c *ConnectDB) Commit() (err error) {
	return c.tx.Commit()
}

func ConvertUnixDateToString(unixDate int64) string {
	return time.Unix(unixDate, 0).Format("2006-01-02 15:04:05")
}

// partnerCode=MOMO&
// orderId=632c9b49a00020b&
// requestId=632c9b49a01020b&
// amount=50000&
// orderInfo=pay+with+MoMo&
// orderType=momo_wallet&
// transId=2843828825&
// resultCode=0&
// message=Giao+d%E1%BB%8Bch+th%C3%A0nh+c%C3%B4ng.&
// payType=qr&
// responseTime=1675748448289&
// extraData=&
// signature=9ecfc80e2a0650375080ec80e34ca34a81d42d4b4561b582128398b2b5ebb2bf
