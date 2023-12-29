package sign

import (
	"bahamut/internal/browser"
	"bahamut/internal/container"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type SignInfo struct {
	Days          int `json:"data.days"`
	Signin        int `json:"data.signin"`
	FinishedAd    int `json:"data.finishedAd"`
	PrjSigninDays int `json:"data.prjSigninDays"`
}

func Sign(con *container.Container, page playwright.Page) (bool, error) {
	log.Println("[簽到] 開始執行簽到")
	homeRes, err := browser.Goto(page, browser.Home)
	if err != nil {
		return false, err
	}

	log.Println("[簽到] 主頁狀態: ", homeRes.Status())

	_, _, err = getSignStatus(con, page)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getSignStatus(con *container.Container, page playwright.Page) (*SignInfo, *SignError, error) {
	reqBody := url.Values{}
	reqBody.Set("action", "2")
	req, err := http.NewRequest(
		http.MethodPost,
		"https://www.gamer.com.tw/ajax/signin.php",
		strings.NewReader(reqBody.Encode()),
	)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := con.HttpClient().Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	dataByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var resBody map[string]interface{}
	if err = json.Unmarshal(dataByte, &resBody); err != nil {
		return nil, nil, err
	}

	if con.Config().GetModulesConfig().Sign.Debug {
		log.Println("[Debug][簽到] 簽到資料:", resBody)
	}

	// 判斷錯誤訊息
	if resBody["error"] != nil {
		var errorMessage SignError
		json.Unmarshal(dataByte, &errorMessage)
		log.Fatalln("[簽到] 簽到失敗：", errorMessage)
		return nil, nil, errorMessage
	}

	var signInfo SignInfo
	json.Unmarshal(dataByte, &signInfo)
	return &signInfo, nil, nil
}
