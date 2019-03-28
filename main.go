package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"encoding/json"
	xmlBasic "encoding/xml"

	"github.com/gorilla/mux"
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

//Response a request
type Response struct {
	XMLName      xmlBasic.Name `xml:"lariz"`
	Date         string        `xml:"date"`
	Result       string        `xml:"result"`
	Msg          string        `xml:"msg"`
	Trxid        string        `xml:"trxid"`
	PartnerTrxid string        `xml:"partner_trxid"`
}

func callBackHandler(w http.ResponseWriter, r *http.Request) {
	// func callBackHandler(h http.Handler) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read XML: %v\n", err)
		w.WriteHeader(500)
	}

	v := new(Lariz)
	errUnmarshal := xmlBasic.Unmarshal(b, v)
	if errUnmarshal != nil {
		log.Printf("cannot unmarshall request: %v\n", errUnmarshal)
		w.WriteHeader(422)
	}

	w.Header().Set("Content-Type", "text/xml")
	buf, _ := xmlBasic.MarshalIndent(&Response{Result: "00", Msg: "Successfully callback", Trxid: v.Trxid, PartnerTrxid: v.PartnerTrxid}, "", "  ")

	io.WriteString(w, string(buf))

	req, _ := json.Marshal(v)
	log.Printf("Request from client: %s\n", req)

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	log.Printf("Ping: %s\n", time.Now())
}

// Define our struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers = make(map[string]string, 1)
	amw.tokenUsers["05f717e5"] = "mba"
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("X-Auth-Token")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/ping", ping).Methods("GET")

	amw := authenticationMiddleware{}
	amw.Populate()

	sub := r.PathPrefix("/mba").Subrouter()
	sub.Path("/callback").HandlerFunc(callBackHandler).Methods("POST")
	sub.Use(amw.authMiddleware)

	// r.HandleFunc("/MBACallback", callBackHandler).Methods("POST")
	// r.Use(amw.authMiddleware) // use for all path
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))

}
