package linenotify

import (
	"bahamut/internal/container"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type NotifyResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Notify(con *container.Container, message string) error {
	reqBody := url.Values{}
	reqBody.Set("message", message)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://notify-api.line.me/api/notify",
		strings.NewReader(reqBody.Encode()),
	)
	if err != nil {
		return err
	}

	token := con.Config().Modules.LineNotify.Token
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
	res, err := con.HttpClient().Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var body NotifyResponse
	if err = json.Unmarshal(rawData, &body); err != nil {
		return err
	}

	if body.Status != 200 {
		e := fmt.Sprintf("[Line Notify] statue: %d, Message: %s", body.Status, body.Message)
		return errors.New(e)
	}

	return nil
}
