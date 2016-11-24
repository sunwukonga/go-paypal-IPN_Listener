package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/notify/", ipnHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page, nothing to see here.")
	})
	fmt.Println("Router initiated.")
	log.Fatal(http.ListenAndServe(":80", mux))
}

func ipnHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ipnHandler called.")
	// Switch for production and live
	isProduction := false

	urlSimulator := "https://www.sandbox.paypal.com/cgi-bin/webscr"
	urlLive := "https://www.paypal.com/cgi-bin/webscr"
	paypalURL := urlSimulator

	if isProduction {
		paypalURL = urlLive
	}
	fmt.Println("URL to return to is: %v", paypalURL)

	// Verify that the POST HTTP Request method was used.
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("No route for %v", r.Method), http.StatusNotFound)
		return
	}
	// *********************************************************
	// HANDSHAKE STEP 1 -- Write back an empty HTTP 200 response
	// *********************************************************
	fmt.Printf("Write Status 200")
	w.WriteHeader(http.StatusOK)

	// *********************************************************
	// HANDSHAKE STEP 2 -- Send POST data (IPN message) back as verification
	// *********************************************************
	// Get Content-Type of request to be parroted back to paypal
	contentType := r.Header.Get("Content-Type")
	// Read the raw POST body
	body, err := ioutil.ReadAll(r.Body)
	// Prepend POST body with required field
	body = append([]byte("cmd=_notify-validate&"), body...)
	// Make POST request to paypal
	req, err := http.NewRequest(http.MethodPost, paypalURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	log.Printf("Response: %v", string(body))
}
