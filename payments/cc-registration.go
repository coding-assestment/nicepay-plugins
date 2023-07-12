package payments

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"nicepay-service/config"
	"nicepay-service/nicepay"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	RegistrationRequest struct {
		PayMethod   string `json:"payMethod"`
		ReferenceNo string `json:"referenceNo"`
		Amt         string `json:"amt" validate:"required"`
		InstmntMon  string `json:"instmntMon"  validate:"required"`
		InstmntType string `json:"instmntType"  validate:"required"`

		Currency        string      `json:"currency" validate:"required"`
		Description     string      `json:"description"`
		BillingNm       string      `json:"billingNm" validate:"required"`
		BillingPhone    string      `json:"billingPhone" validate:"required"`
		BillingEmail    string      `json:"billingEmail" validate:"required"`
		BillingAddr     string      `json:"billingAddr" validate:"required"`
		BillingCity     string      `json:"billingCity" validate:"required"`
		BillingState    string      `json:"billingState"  validate:"required"`
		BillingPostCd   string      `json:"billingPostCd"  validate:"required"`
		BillingCountry  string      `json:"billingCountry" validate:"required"`
		DeliveryNm      string      `json:"deliveryNm"`
		DeliveryPhone   string      `json:"deliveryPhone"`
		DeliveryAddr    string      `json:"deliveryAddr"`
		DeliveryCity    string      `json:"deliveryCity"`
		DeliveryState   string      `json:"deliveryState"`
		DeliveryPostCd  string      `json:"deliveryPostCd"`
		DeliveryCountry string      `json:"deliveryCountry"`
		RecurrOpt       string      `json:"recurrOpt"  validate:"required"`
		TimeStamp       string      `json:"timeStamp"`
		ReqDt           string      `json:"reqDt"`
		ReqTm           string      `json:"reqTm"`
		IMid            string      `json:"iMid"`
		MerchantToken   string      `json:"merchantToken"`
		DbProcessUrl    string      `json:"dbProcessUrl"`
		UserIP          string      `json:"userIP"`
		GoodsNm         string      `json:"goodsNm" validate:"required"`
		NotaxAmt        string      `json:"notaxAmt"`
		ReqDomain       string      `json:"reqDomain"`
		Fee             string      `json:"fee"`
		Vat             string      `json:"vat"`
		CartData        interface{} `json:"cartData"  validate:"required"`
	}

	RegistrationResponse struct {
		Amt          string `json:"amt"`
		BankCd       string `json:"bankCd"`
		BillingNm    string `json:"billingNm"`
		Currency     string `json:"currency"`
		Description  string `json:"description"`
		GoodsNm      string `json:"goodsNm"`
		MitraCd      string `json:"mitraCd"`
		PayMethod    string `json:"payMethod"`
		PayNo        string `json:"payNo"`
		PayValidDt   string `json:"payValidDt"`
		PayValidTm   string `json:"payValidTm"`
		PaymentExpDt string `json:"paymentExpDt"`
		PaymentExpTm string `json:"paymentExpTm"`
		QrContent    string `json:"qrContent"`
		QrUrl        string `json:"qrUrl"`
		ReferenceNo  string `json:"referenceNo"`
		RequestURL   string `json:"requestURL"`
		ResultCd     string `json:"resultCd"`
		ResultMsg    string `json:"resultMsg"`
		TXid         string `json:"tXid"`
		TransDt      string `json:"transDt"`
		TransTm      string `json:"transTm"`
		VacctNo      string `json:"vacctNo"`
		VacctValidDt string `json:"vacctValidDt"`
		VacctValidTm string `json:"vacctValidTm"`
	}
)

func (response *RegistrationResponse) toString() string {
	data, _ := json.Marshal(response)
	return string(data)
}

func BindResponse(data []byte) *RegistrationResponse {
	var response *RegistrationResponse
	/*var defaultResponse map[string]interface{} */
	json.Unmarshal([]byte(data), &response)
	return response
}

func (rr *RegistrationRequest) toString() string {
	data, _ := json.Marshal(rr)
	return string(data)
}

func (rr *RegistrationRequest) makeUniqRefNumber() {
	lock.Lock()
	defer lock.Unlock()
	seq++
	rr.ReferenceNo = rr.ReferenceNo + "-" + fmt.Sprintf("%x", seq)
}

func (rr *RegistrationRequest) setTimeStamp() {
	t := time.Now()

	formatted := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	rr.TimeStamp = formatted
}

func (rr *RegistrationRequest) MakeMerchantToken() {
	rr.setTimeStamp()
	rr.makeUniqRefNumber()
	combinedTokenMentah := rr.TimeStamp + rr.IMid + rr.ReferenceNo + rr.Amt + config.NICEPAY_MERCHANT_KEY

	h := sha256.New()
	h.Write([]byte(combinedTokenMentah))
	bs := h.Sum(nil)

	rr.MerchantToken = fmt.Sprintf("%x", bs)
}

func (rr *RegistrationRequest) prepareRegistrationReq() {
	rr.MakeMerchantToken()
	rr.Description = rr.Description + rr.ReferenceNo
	rr.GoodsNm = rr.GoodsNm + rr.ReferenceNo

	t := time.Now()

	formattedDate := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	rr.ReqDt = formattedDate

	formattedTime := fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
	rr.ReqTm = formattedTime

	rr.PayMethod = "01" // for cc
}

func NewRegistrationRequest() *RegistrationRequest {
	/* default / template value */
	return &RegistrationRequest{
		/*  data required*/
		/*
			Amt:         "600",
			InstmntMon:  "1",
			InstmntType: "2",

			Currency:        "IDR",
			BillingNm:       "John Doe",
			BillingPhone:    "082111111111",
			BillingEmail:    "john@example.com",
			BillingAddr:     "Jl. Jend. Sudirman No. 28",
			BillingCity:     "Jakarta Pusat",
			BillingState:    "DKI Jakarta",
			BillingPostCd:   "10210",
			BillingCountry:  "Indonesia",
			RecurrOpt:       "1",
			CartData:      "{}",
		*/
		/* data optional */
		PayMethod:       "01",
		ReferenceNo:     "NCPAY-",
		Description:     "Payment of Invoice No ",
		GoodsNm:         "Payment of Invoice No ",
		DeliveryNm:      "John Doe",
		DeliveryPhone:   "02112345678",
		DeliveryAddr:    "Jl. Jend. Sudirman No. 28",
		DeliveryCity:    "Jakarta Pusat",
		DeliveryState:   "DKI Jakarta",
		DeliveryPostCd:  "10210",
		DeliveryCountry: "Indonesia",
		TimeStamp:       "",
		ReqDt:           "",
		ReqTm:           "",
		IMid:            config.NICEPAY_IMID,
		MerchantToken:   "",
		DbProcessUrl:    "http://httpresponder.com/nicepay",
		UserIP:          "::1",
		NotaxAmt:        "0",
		ReqDomain:       "http://localhost/",
		Fee:             "0",
		Vat:             "0",
	}
}

func CreateRegistration(c echo.Context) error {
	registrationReq := NewRegistrationRequest()
	if err := c.Bind(registrationReq); err != nil {
		return err
	}
	if err := c.Validate(registrationReq); err != nil {
		c.Logger().Error(err)
		log.Println("Registration validation Failed : " + err.Error())
		log.Print("NicePay Request Iquiry: " + registrationReq.toString())
		return err
	}

	registrationReq.prepareRegistrationReq()

	NicePay := nicepay.NewInstance()
	NicePay.Operation("requestCC")
	requestString := registrationReq.toString()
	resultNicePayRequest, _ := NicePay.ApiRequest(requestString)

	responseNicePay := BindResponse(resultNicePayRequest)

	result := struct {
		/* Request *RegistrationRequest `json:"request"` */
		Request         interface{} `json:"Request"`
		ResponseNicePay interface{} `json:"ResponseNicePay"`
	}{Request: registrationReq, ResponseNicePay: responseNicePay}

	if responseNicePay.ResultCd == "0000" {
		c.Logger().Info("Registration Success")
		log.Println("Registration Success: " + responseNicePay.toString())
	} else {
		c.Logger().Error(errors.New("Registration Failed"))
		log.Println("Registration Failed: " + responseNicePay.toString())
	}

	log.Print("NicePay Request Iquiry: " + registrationReq.toString())
	log.Print("NicePay Response Iquiry: " + responseNicePay.toString())

	return c.JSON(http.StatusCreated, result)

}
