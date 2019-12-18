package infobip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/flip-id/infobip/pkg/infobip/database"
	"github.com/flip-id/infobip/pkg/infobip/models"
	"github.com/t-tiger/gorm-bulk-insert"
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

func SendSMS(from string, to []Destination, text string) ResponseBody {
	url := fmt.Sprintf("https://%s%s", baseUrl, SMS_ENDPOINT)
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

	logToDB(resBody)

	return resBody
}

func logToDB(r ResponseBody) {
	db, _ := database.Connect("root", "root", "infobip", "localhost")
	var statusLogs []interface{}

	for _, m := range r.Messages {
		statusLogs = append(statusLogs, models.StatusLog{
			BulkId:      r.BulkId,
			MessageId:   m.MessageId,
			PhoneNumber: m.To,
			StatusCode:  m.Status.Name,
		})
	}

	err := gormbulk.BulkInsert(db, statusLogs, 100)

	if err != nil {
		panic(err)
	}
}

// Update status log status
func UpdateStatus(messageId, newStatus string) {
	db, _ := database.Connect("root", "root", "infobip", "localhost")
	var statusLog = models.StatusLog{}

	db.First(&statusLog, "message_id = ?", messageId)

	db.Model(&statusLog).Update("status_code", newStatus)
}
