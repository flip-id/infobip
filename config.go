package infobip

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	SmsEndpoint = "/sms/2/text/advanced"
)

var (
	baseUrl   string = os.Getenv("INFOBIP_BASE_URL")
	apiKey    string = os.Getenv("INFOBIP_API_KEY")
	notifyUrl string = os.Getenv("INFOBIP_NOTIFY_URL")
)
