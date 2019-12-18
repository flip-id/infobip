package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/flip-id/infobip/pkg/infobip"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type response struct {
	Name string
}

type Description struct {
	To        string
	MessageId *string
}

type Destination struct {
	To string
}

type RequestMessage struct {
	From         string                `json:"from"`
	Text         string                `json:"text"`
	Destinations []infobip.Destination `json:"destinations"`
}

type CallbackData struct {
	Results []struct {
		BulkID       string `json:"bulkId"`
		CallbackData string `json:"callbackData"`
		DoneAt       string `json:"doneAt"`
		Error        struct {
			Description string `json:"description"`
			GroupID     int    `json:"groupId"`
			GroupName   string `json:"groupName"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Permanent   bool   `json:"permanent"`
		} `json:"error"`
		MccMnc    string `json:"mccMnc"`
		MessageID string `json:"messageId"`
		Price     struct {
			Currency        string  `json:"currency"`
			PricePerMessage float64 `json:"pricePerMessage"`
		} `json:"price"`
		SentAt   string `json:"sentAt"`
		SmsCount int    `json:"smsCount"`
		Status   struct {
			Description string `json:"description"`
			GroupID     int    `json:"groupId"`
			GroupName   string `json:"groupName"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
		} `json:"status"`
		To string `json:"to"`
	} `json:"results"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/v1/send-sms", handleSendSms).Methods("POST")
	r.HandleFunc("/v1/notify", handleNotify).Methods("POST")
	fmt.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleSendSms(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errMsg := fmt.Sprintf("error: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	var payload = RequestMessage{}

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	fmt.Println(payload)

	resBody := infobip.SendSMS(payload.From, payload.Destinations, payload.Text)
	resBodyBytes := new(bytes.Buffer)
	json.NewEncoder(resBodyBytes).Encode(resBody)

	w.WriteHeader(http.StatusOK)
	w.Write(resBodyBytes.Bytes())
}

func handleNotify(w http.ResponseWriter, r *http.Request) {
	callbackData := CallbackData{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &callbackData)

	if err != nil {
		panic(err)
	}

	infobip.UpdateStatus(callbackData.Results[0].MessageID, callbackData.Results[0].Status.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Success"}`))
}
