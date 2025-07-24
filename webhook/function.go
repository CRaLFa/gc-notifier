package webhook

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func init() {
	functions.HTTP("GetGroupID", getGroupID)
}

func getGroupID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cb, err := webhook.ParseRequest(os.Getenv("LINE_CHANNEL_SECRET"), r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook request: %v\n", err)
		http.Error(w, "Failed to parse webhook request", http.StatusBadRequest)
		return
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			fmt.Printf("MessageEvent: %#v\n", e)
		}
	}
}
