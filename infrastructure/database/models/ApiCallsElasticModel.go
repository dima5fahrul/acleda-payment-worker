package models

import "time"

type ApiCallsElasticModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	Track          string    `gorm:"column:track"`
	Service        string    `gorm:"column:service"`
	Webtype        string    `gorm:"column:webtype"`
	Merchant       string    `gorm:"column:merchant"`
	Msisdn         string    `gorm:"column:msisdn"`
	URL            string    `gorm:"column:url"`
	Method         string    `gorm:"column:method"`
	RequestQuery   string    `gorm:"column:request_query"`
	RequestBody    string    `gorm:"column:request_body"`
	ResponseBody   string    `gorm:"column:response_body"`
	StatusCode     int       `gorm:"column:status_code"`
	RequestHeader  string    `gorm:"column:request_header"`
	ResponseHeader string    `gorm:"column:response_header"`
	Latency        string    `gorm:"column:latency"`
	Error          string    `gorm:"column:error"`
	TransactionID  string    `gorm:"column:transaction_id;index"`
}
