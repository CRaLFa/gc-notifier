package notifier

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/samber/lo"
)

type parameters struct {
	GarbageType   string `json:"garbageType"`
	IsPreviousDay bool   `json:"isPreviousDay"`
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func init() {
	functions.HTTP("PostNotification", postNotification)
}

func postNotification(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var params parameters
	if err := json.Unmarshal(body, &params); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bot, err := messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create LINE bot client: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := lo.Ternary(params.IsPreviousDay, "明日", "今日")
	t := time.Now()
	if params.IsPreviousDay {
		t = t.AddDate(0, 0, 1)
	}

	req := &messaging_api.PushMessageRequest{
		To: os.Getenv("LINE_GROUP_ID"),
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: fmt.Sprintf("%s (%s) は %s の収集日です", d, weekdayToJa(t.In(jst).Format("1/2・Mon")), params.GarbageType),
			},
		},
	}
	_, err = bot.PushMessage(req, uuid.NewString())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send message: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func weekdayToJa(s string) string {
	return strings.NewReplacer(
		"Sun", "日",
		"Mon", "月",
		"Tue", "火",
		"Wed", "水",
		"Thu", "木",
		"Fri", "金",
		"Sat", "土",
	).Replace(s)
}
