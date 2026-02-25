package entities

import "time"

type ApiCall struct {
	ID             int64
	CreatedAt      time.Time
	Track          string
	Service        string
	Webtype        string
	Merchant       string
	Msisdn         string
	URL            string
	Method         string
	RequestQuery   string
	RequestBody    string
	ResponseBody   string
	StatusCode     int
	RequestHeader  string
	ResponseHeader string
	Latency        string
	Error          string
	TransactionID  string
}

type Incoming struct {
	ID            string
	TransactionID string
	CreatedAt     time.Time
	Track         string
	Service       string
	Webtype       string
	Path          string
	Merchant      string
	IP            string
	Method        string
	RequestQuery  string
	RequestHeader string
	RequestBody   string
	ResponseBody  string
	StatusCode    int
	Latency       string
	UserAgent     string
	Device        string
	Browser       string
	Callback      string
	Country       string
	ChannelCode   string
	CallbackUrl   string
	Description   string
	PaymentMethod string
	Event         string
	Email         string
	Curency       string
	Save          bool
}
