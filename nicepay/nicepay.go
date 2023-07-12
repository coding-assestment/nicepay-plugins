package nicepay

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"nicepay-service/config"
	"strings"
)

type (
	SocketConfig struct {
		ApiUrl    string
		Port      int64
		Status    string
		Headers   string
		body      string
		Request   string
		Errorcode string
		Errormsg  string
		Log       string
		Timeout   int64
	}
)

func NewInstance() *SocketConfig {
	return &SocketConfig{
		ApiUrl:    "",
		Port:      443,
		Status:    "",
		Headers:   "",
		body:      "",
		Request:   "",
		Errorcode: "",
		Errormsg:  "",
		Log:       "",
		Timeout:   config.NICEPAY_TIMEOUT_CONNECT,
	}
}

func (sc *SocketConfig) Operation(typeOp string) {
	if typeOp == "testing" {
		sc.ApiUrl = config.TESTING_REQ_URL
	}
	if typeOp == "requestCC" {
		sc.ApiUrl = config.NICEPAY_REQ_URL
	}
	if typeOp == "creditCard" {
		sc.ApiUrl = config.NICEPAY_REQ_CC_URL
	}
	if typeOp == "checkPaymentStatus" {
		sc.ApiUrl = config.NICEPAY_ORDER_STATUS_URL
	}
	if typeOp == "cancel" {
		sc.ApiUrl = config.NICEPAY_CANCEL_URL
	}
}

func (sc *SocketConfig) getHost() string {

	url, err := url.Parse(sc.ApiUrl)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname
}

func (sc *SocketConfig) SetBody(payloads string) {
	sc.body = payloads
}

func (sc *SocketConfig) GetBody() string {
	return sc.body
}

func (sc *SocketConfig) ApiRequest(payloads string) ([]byte, error) {
	sc.body = payloads

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	/*
		fmt.Println(payloads)
		test := `{
		    "tXid": "IONPAYTEST01202307121544390194",
		    "amt": "123",
		    "referenceNo": "NCPAY-64ae67ea289600.02399925",
		    "timeStamp": "20230712105506",
		    "iMid": "IONPAYTEST",
		    "merchantToken": "2748eb990d8aa494226509f5d1e1397e48fe68f8ff5313f58bbde99d9ae5931f"
		}`
	*/
	fmt.Println(payloads)
	response, err := client.Post(
		sc.ApiUrl,
		"application/json",
		bytes.NewBufferString(payloads),
		// bytes.NewBufferString(test),
	)

	if err != nil {
		fmt.Printf("%s", err)
		return []byte{}, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return contents, err
		}
		return contents, err
	}

	return []byte{}, err
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
