package infobip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	To string
	// MessageId string
}

type ReqMessage struct {
	From               string
	Destinations       []Destination
	Text               string
	IntermediateReport bool
	NotifyContentType  string
}

type ReqMessages struct {
	// BulkId string
	Messages []ReqMessage
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

func SendSMS(from string, to []Destination, text string) ResponseBody {
	url := fmt.Sprintf("https://%s%s", baseUrl, SmsEndpoint)
	payload, err := json.Marshal(ReqMessages{Messages: []ReqMessage{
		{from, to, text, true, "application/json"},
	}})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Authorization", fmt.Sprintf("App %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	resBody := ResponseBody{}

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resBody)

	if err != nil {
		panic(err)
	}

	return resBody
}
