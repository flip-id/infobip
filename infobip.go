package infobip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	To        string
	MessageId string
	Status    struct {
		GroupId     int
		GroupName   string
		Id          int
		Name        string
		Description string
	}
}

type ResponseBody struct {
	BulkId   string
	Messages []Message
}

type Destination struct {
	To string `json:"to"`
	// MessageId string
}

type ReqMessage struct {
	From               string        `json:"from"`
	Destinations       []Destination `json:"destinations"`
	Text               string        `json:"text"`
	IntermediateReport bool          `json:"intermediateReport"`
	NotifyContentType  string        `json:"notifyContentType"`
	NotifyUrl          string        `json:"notifyUrl"`
}

type ReqMessages struct {
	// BulkId string
	Messages []ReqMessage `json:"messages"`
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

func SendSMS(from string, to []Destination, text string) (ResponseBody, error) {
	url := fmt.Sprintf("%s%s", baseUrl, SmsEndpoint)
	payload, err := json.Marshal(ReqMessages{Messages: []ReqMessage{
		{
			From:               from,
			Destinations:       to,
			Text:               text,
			IntermediateReport: true,
			NotifyContentType:  "application/json",
			NotifyUrl:          notifyUrl,
		},
	}})

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("App %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	if res.StatusCode >= 400 {
		log.Error(res)
		return ResponseBody{}, errors.New("failed to send SMS")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	resBody := ResponseBody{}

	if err != nil {
		log.Error(res)
		return ResponseBody{}, err
	}

	err = json.Unmarshal(body, &resBody)

	if err != nil {
		log.Error(res)
		return ResponseBody{}, err
	}

	return resBody, nil
}
