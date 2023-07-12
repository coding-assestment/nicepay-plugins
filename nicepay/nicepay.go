package nicepay

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"nicepay-service/config"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

func (sc *SocketConfig) ApiRequest(payloads string) ([]byte, string) {
	sc.body = payloads

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	// fmt.Println(payloads)
	// test := `{
	// 	"amt": "15000",
	// 	"billingAddr": "Jl. Jend. Sudirman No. 28",
	// 	"billingCity": "Jakarta Pusat",
	// 	"billingCountry": "Indonesia",
	// 	"billingEmail": "john@example.com",
	// 	"billingNm": "John Doe",
	// 	"billingPhone": "082111111111",
	// 	"billingPostCd": "10210",
	// 	"billingState": "DKI Jakarta",
	// 	"cartData": "{\"count\": \"1\",\"item\": [{\"img_url\": \"https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/iphone11-select-2019-family?wid=882&hei=1058&fmt=jpeg&qlt=80&op_usm=0.5,0.5&.v=1567022175704\",\"goods_name\": \" iPhone 11 \",\"goods_detail\": \"A new dual‑camera system captures more of what you see and love. The fastest chip ever in a smartphone and all‑day battery life let you do more and charge less. And the highest‑quality video in a smartphone, so your memories look better than ever.\",\"goods_amt\":\"15000\"}]}",
	// 	"currency": "IDR",
	// 	"dbProcessUrl": "http://httpresponder.com/nicepay",
	// 	"deliveryAddr": "Jl. Jend. Sudirman No. 28",
	// 	"deliveryCity": "Jakarta Pusat",
	// 	"deliveryCountry": "Indonesia",
	// 	"deliveryNm": "John Doe",
	// 	"deliveryPhone": "02112345678",
	// 	"deliveryPostCd": "10210",
	// 	"deliveryState": "DKI Jakarta",
	// 	"description": "Payment of Invoice No NCPAY-64ad872fc76ec4.25417383",
	// 	"fee": "0",
	// 	"goodsNm": "Payment of Invoice No NCPAY-64ad872fc76ec4.25417383",
	// 	"iMid": "IONPAYTEST",
	// 	"instmntMon": "1",
	// 	"instmntType": "2",
	// 	"merchantToken": "09b1f556e9537d2997cd545f8fd1e834edcf86e6a505c494243d76853b774c16",
	// 	"notaxAmt": "0",
	// 	"payMethod": "01",
	// 	"recurrOpt": "1",
	// 	"referenceNo": "NCPAY-64ad872fc76ec4.25417383",
	// 	"reqDomain": "http://localhost/",
	// 	"reqDt": "20230711",
	// 	"reqTm": "184732",
	// 	"timeStamp": "20230711184732",
	// 	"userIP": "::1",
	// 	"vat": "0"
	//   }`
	response, err := client.Post(
		sc.ApiUrl,
		"application/json; charset=UTF-8",
		bytes.NewBufferString(payloads),
		// bytes.NewBufferString(test),
	)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		return contents, payloads
	}

	return []byte{}, payloads
}

func (sc *SocketConfig) ApiRequestOld() {

	// con, err := net.Dial("tcp", sc.ApiUrl+":"+sc.Port)
	spew.Dump(fmt.Sprintf("%s:%d", sc.ApiUrl, sc.Port))

	// con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sc.getHost(), sc.Port))
	// con, err := net.Dial("tcp", sc.ApiUrl)
	// con, err := net.Dial("tcp", "webcode.me:80")
	con, err := net.Dial("tcp", "dev.nicepay.co.id:443")
	checkError(err)

	req := "POST / HTTP/1.0\r\n" +
		// "Host: " + sc.getHost() + "\r\n" +
		"Host: " + "dev.nicepay.co.id" + "\r\n" +
		"Path: " + "/nicepay/direct/v2/registration" + "\r\n" +
		"User-Agent: Go client\r\n\r\n" +
		"Connection: close\r\n" +
		"Content-type: application/json\r\n" +
		"Content-length: 0\r\n" +
		"Accept: */*\r\n" +
		"{} \r\n" +
		"\r\n"

	// req := "GET / HTTP/1.0\r\n" +
	// 	"Host: webcode.me\r\n" +
	// 	"User-Agent: Go client\r\n\r\n"

	// req := "POST " . $uri . " HTTP/1.0\r\n";
	// $request .= "Connection: close\r\n";
	// $request .= "Host: " . $host . "\r\n";
	// $request .= "Content-type: application/json\r\n";
	// $request .= "Content-length: " . strlen ( $postdata ) . "\r\n";
	// $request .= "Accept: */*\r\n";
	// $request .= "\r\n";
	// $request .= $postdata . "\r\n";
	// $request .= "\r\n";

	_, err = con.Write([]byte(req))
	checkError(err)

	res, err := ioutil.ReadAll(con)
	checkError(err)

	fmt.Println(string(res))
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
