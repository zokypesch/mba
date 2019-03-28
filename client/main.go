package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	xmlBasic "encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/vigneshuvi/GoDateFormat"
	utl "github.com/zokypesch/mba/utils"
)

// Lariz set a response
type Lariz struct {
	XMLName xmlBasic.Name `xml:"lariz"`
	// Text         string        `xml:",chardata"`
	Date         string `xml:"date"`
	Result       string `xml:"result"`
	Msg          string `xml:"msg"`
	Trxid        string `xml:"trxid"`
	PartnerTrxid string `xml:"partner_trxid"`
	Saldo        string `xml:"saldo"`
}

// XMLRPCCall func test a service
// func XMLRPCCall(req *Lariz) (resps *Lariz, err error) {
func XMLRPCCall(req map[string]string) (resps *Lariz, err error) {

	buf, _ := xmlBasic.MarshalIndent(utl.Evoucher(req), "", "  ")
	// buf, _ := xmlBasic.MarshalIndent(req, "", "  ")

	client := &http.Client{}
	// requestClient, _ := http.NewRequest("POST", "http://localhost/mba/callback", bytes.NewBuffer(buf))
	requestClient, _ := http.NewRequest("POST", "http://202.158.48.172:2929/topup", bytes.NewBuffer(buf))

	// requestClient.Header.Add("X-Auth-Token", "05f717e5")
	requestClient.Header.Add("Content-Type", "text/xml")

	resp, err := client.Do(requestClient)

	if err != nil {
		log.Printf("errClient : %v", err)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Http Status Code : %d", resp.StatusCode)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Forbidden access")
	}

	v := &Lariz{}
	errUnmarshal := xmlBasic.Unmarshal(bodyBytes, v)
	if errUnmarshal != nil {
		log.Printf("Cannot unmarshall %v", errUnmarshal)
	}

	return v, errUnmarshal
}

// XMLRPCCallLocal for reply for local
func XMLRPCCallLocal(req *Lariz) (resps *Lariz, err error) {
	buf, _ := xmlBasic.MarshalIndent(req, "", "  ")

	client := &http.Client{}
	requestClient, _ := http.NewRequest("POST", "http://goexample.staging.svc.cluster.local/mba/callback", bytes.NewBuffer(buf))
	// requestClient, _ := http.NewRequest("POST", "http://localhost/mba/callback", bytes.NewBuffer(buf))

	requestClient.Header.Add("X-Auth-Token", "05f717e5")
	requestClient.Header.Add("Content-Type", "text/xml")

	resp, err := client.Do(requestClient)

	if err != nil {
		log.Printf("errClient : %v", err)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Http Status Code : %d", resp.StatusCode)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Forbidden access")
	}

	v := &Lariz{}
	errUnmarshal := xmlBasic.Unmarshal(bodyBytes, v)
	if errUnmarshal != nil {
		log.Printf("Cannot unmarshall %v", errUnmarshal)
	}

	return v, errUnmarshal
}

// GetToday for get today format
func GetToday(format string) (todayString string) {
	today := time.Now()
	todayString = today.Format(format)
	return
}

// GetSignature for get today format
func GetSignature(phone string) (times string, signature string) {
	timeSignature := GetToday(GoDateFormat.ConvertFormat("HHMMSS"))
	rotateNum := phone[len(phone)-4 : len(phone)]
	val1, _ := strconv.Atoi(timeSignature)
	val2, _ := strconv.Atoi(rotateNum)

	return timeSignature, b64.StdEncoding.EncodeToString([]byte(strconv.Itoa(val1 + val2)))
}

// GenerateNewRequest for new request
func GenerateNewRequest(
	command string,
	msisdn string,
	partnerTrxid string,
	product string,
	userid string,
) map[string]string {
	times, sign := GetSignature(msisdn)

	return map[string]string{
		"command":       command,
		"msisdn":        msisdn,
		"partner_trxid": partnerTrxid,
		"product":       product,
		"signature":     sign,
		"time":          times,
		"userid":        userid,
	}
}

func main() {
	min := 10
	max := 30

	newRequest := GenerateNewRequest("TOPUP", "081216074307", strconv.Itoa(rand.Intn(max-min)+min), "S10", "62812345678905")

	v, _ := XMLRPCCall(newRequest)
	// call a service
	reply, _ := json.Marshal(v)
	log.Printf("Result : %s\n", reply)

	l, _ := XMLRPCCallLocal(&Lariz{Date: fmt.Sprintf("%s", time.Now()), Result: "00", Msg: "success", Trxid: "xxx", PartnerTrxid: "xxx-xxx", Saldo: "10000"})

	replyLocal, _ := json.Marshal(l)
	log.Printf("Result From local : %s\n", replyLocal)
}
