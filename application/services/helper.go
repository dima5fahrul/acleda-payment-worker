package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"payment-airpay/infrastructure/database"
	"payment-airpay/infrastructure/database/models"
	pkg "payment-airpay/infrastructure/gateway"
)

func formatErrorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func SaveAPICall(
	ctx context.Context,
	param pkg.APICall,
	merchant string,
	err error,
	service string,
	track string,
	msisdn string,
	webtype string,
	transactionId string,
) {
	resp := param.GetAPICall()
	data := models.ApiCallsElasticModel{
		ID:             0,
		CreatedAt:      time.Now(),
		Track:          track,
		Service:        service,
		Webtype:        webtype,
		Merchant:       merchant,
		Msisdn:         msisdn,
		URL:            resp.RequestURL,
		Method:         resp.Method,
		RequestQuery:   resp.RequestQuery,
		RequestBody:    resp.RequestBody,
		ResponseBody:   resp.ResponseBody,
		StatusCode:     resp.ResponseStatusCode,
		RequestHeader:  resp.RequestHeaders,
		ResponseHeader: resp.ResponseHeaders,
		Latency:        resp.RequestLatency,
		Error:          formatErrorToString(err),
		TransactionID:  transactionId,
	}

	if database.ElasticsearchClient != nil {
		_, err := database.ElasticsearchClient.Index("api_call_logs").
			Request(&data).
			Do(context.Background())
		if err != nil {
			fmt.Printf("Failed to save API call to Elasticsearch: %v\n", err)
		}
	} else {
		fmt.Println("ElasticsearchClient is not initialized")
	}
}

func ValidateRequest(req interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range validationErrs {
				errMsgs = append(errMsgs, e.Field()+" is "+e.Tag())
			}
			return errors.New(strings.Join(errMsgs, ", "))
		}
		return err
	}
	return nil
}
