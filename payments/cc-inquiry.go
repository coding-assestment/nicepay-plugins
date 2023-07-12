package payments

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nicepay-service/config"
	"nicepay-service/nicepay"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	InquiryRequest struct {
		TXid        string `json:"tXid" validate:"required,lte=30"`
		Amt         string `json:"amt" validate:"required,gte=0"`
		ReferenceNo string `json:"referenceNo" validate:"required,lte=40"`

		TimeStamp     string `json:"timeStamp"`
		IMid          string `json:"iMid"`
		MerchantToken string `json:"merchantToken"`
	}

	InquiryNicePayResponse struct {
		TXid           string `json:"tXid"`
		IMid           string `json:"iMid"`
		Currency       string `json:"currency"`
		Amt            string `json:"amt"`
		InstmntMon     string `json:"instmntMon"`
		InstmntType    string `json:"instmntType"`
		ReferenceNo    string `json:"referenceNo"`
		GoodsNm        string `json:"goodsNm"`
		PayMethod      string `json:"payMethod"`
		BillingNm      string `json:"billingNm"`
		ReqDt          string `json:"reqDt"`
		ReqTm          string `json:"reqTm"`
		Status         string `json:"status"`
		ResultCd       string `json:"resultCd"`
		ResultMsg      string `json:"resultMsg"`
		CardNo         string `json:"cardNo"`
		PreauthToken   string `json:"preauthToken"`
		AcquBankCd     string `json:"acquBankCd"`
		IssuBankCd     string `json:"issuBankCd"`
		VacctValidDt   string `json:"vacctValidDt"`
		VacctValidTm   string `json:"vacctValidTm"`
		VacctNo        string `json:"vacctNo"`
		BankCd         string `json:"bankCd"`
		PayNo          string `json:"payNo"`
		MitraCd        string `json:"mitraCd"`
		ReceiptCode    string `json:"receiptCode"`
		CancelAmt      string `json:"cancelAmt"`
		TransDt        string `json:"transDt"`
		TransTm        string `json:"transTm"`
		RecurringToken string `json:"recurringToken"`
		CcTransType    string `json:"ccTransType"`
		PayValidDt     string `json:"payValidDt"`
		PayValidTm     string `json:"payValidTm"`
		MRefNo         string `json:"mRefNo"`
		AcquStatus     string `json:"acquStatus"`
		CardExpYymm    string `json:"cardExpYymm"`
		AcquBankNm     string `json:"acquBankNm"`
		IssuBankNm     string `json:"issuBankNm"`
		DepositDt      string `json:"depositDt"`
		DepositTm      string `json:"depositTm"`
		PaymentExpDt   string `json:"paymentExpDt"`
		PaymentExpTm   string `json:"paymentExpTm"`
		PaymentTrxSn   string `json:"paymentTrxSn"`
		CancelTrxSn    string `json:"cancelTrxSn"`
		UserId         string `json:"userId"`
		ShopId         string `json:"shopId"`
		AuthNo         string `json:"authNo"`
	}
	InqueryResponse struct {
		TXid        string `json:"tXid,omitempty"`
		ReqDt       string `json:"reqDt,omitempty"`
		ReqTm       string `json:"reqTm,omitempty"`
		TransDt     string `json:"transDt,omitempty"`
		TransTm     string `json:"transTm,omitempty"`
		ResultCd    string `json:"resultCd,omitempty"`
		ResultMsg   string `json:"resultMsg,omitempty"`
		Amt         string `json:"amt,omitempty"`
		ReferenceNo string `json:"referenceNo,omitempty"`
	}
)

func NewInquiryResponse(data []byte, txid string, refNo string) *InqueryResponse {
	var response *InqueryResponse
	json.Unmarshal([]byte(data), &response)
	response.TXid = txid
	response.ReferenceNo = refNo
	return response
}

func (response *InquiryNicePayResponse) toString() string {
	data, _ := json.Marshal(response)
	return string(data)
}
func (response *InqueryResponse) toString() string {
	data, _ := json.Marshal(response)
	return string(data)
}

func BindResponseInquiry(data []byte) *InquiryNicePayResponse {
	var response *InquiryNicePayResponse
	json.Unmarshal([]byte(data), &response)
	return response
}

func (rr *InquiryRequest) toString() string {
	data, _ := json.Marshal(rr)
	return string(data)
}

func (rr *InquiryRequest) setTimeStamp() {
	t := time.Now()

	formatted := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	rr.TimeStamp = formatted
}

func (rr *InquiryRequest) MakeMerchantToken() {
	rr.setTimeStamp()
	combinedTokenMentah := rr.TimeStamp + rr.IMid + rr.ReferenceNo + string(rr.Amt) + config.NICEPAY_MERCHANT_KEY

	h := sha256.New()
	h.Write([]byte(combinedTokenMentah))
	bs := h.Sum(nil)

	rr.MerchantToken = fmt.Sprintf("%x", bs)
}

func (rr *InquiryRequest) prepareInqueryReq() {
	rr.MakeMerchantToken()
}

func NewInquiryRequest() *InquiryRequest {
	return &InquiryRequest{
		TXid:        "",
		Amt:         "",
		ReferenceNo: "",

		TimeStamp:     "",
		IMid:          config.NICEPAY_IMID,
		MerchantToken: "",
	}
}

func GetInquiry(c echo.Context) error {
	inquiryRequest := NewInquiryRequest()
	if err := c.Bind(inquiryRequest); err != nil {
		return err
	}
	ClientRequest := inquiryRequest
	log.Println(ClientRequest)
	inquiryRequest.prepareInqueryReq()

	if err := c.Validate(inquiryRequest); err != nil {
		return err
	}

	NicePay := nicepay.NewInstance()
	NicePay.Operation("checkPaymentStatus")
	nicePayStringRequest := inquiryRequest.toString()
	resultNicePayRequest, _ := NicePay.ApiRequest(nicePayStringRequest)

	nicePayResponse := BindResponseInquiry(resultNicePayRequest)
	response := NewInquiryResponse(resultNicePayRequest, inquiryRequest.TXid, inquiryRequest.ReferenceNo)

	result := struct {
		ClientRequest        interface{} `json:"ClientRequest,omitempty"`
		NicePayStringRequest interface{} `json:"NicePayStringRequest,omitempty"`
		NicePayResponse      interface{} `json:"NicePayResponse,omitempty"`
		Response             interface{} `json:"Response,omitempty"`
	}{ClientRequest: ClientRequest, NicePayStringRequest: nicePayStringRequest, NicePayResponse: nicePayResponse, Response: response}

	if response.ResultCd == "0000" {
		log.Println("Success Iquiry: " + response.toString())
	} else {
		log.Println("Failed Iquiry: " + response.toString())
	}
	log.Print("NicePay Request Iquiry: " + ClientRequest.toString())
	log.Print("NicePay Response Iquiry: " + nicePayResponse.toString())

	return c.JSON(http.StatusCreated, result)

}
