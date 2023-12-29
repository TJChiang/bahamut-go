package sign

import (
	"bahamut/internal/browser"
	"bahamut/internal/container"
	"encoding/json"
	"log"

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

	_, _, err = getSignStatus(page)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getSignStatus(page playwright.Page) (*SignInfo, *SignError, error) {
	// 在主頁跑簽到 API
	data, err := page.Evaluate(`async () => {
		const controller = new AbortController();
	    setTimeout(() => controller.abort(), 30000);
		const res = await fetch("https://www.gamer.com.tw/ajax/signin.php", {
	        method: "POST",
	        headers: {
	            "Content-Type": "application/x-www-form-urlencoded",
	        },
	        body: "action=2",
	        signal: controller.signal,
	    });
	    return res.json();
	}`)
	if err != nil {
		return nil, nil, err
	}

	log.Println("[Debug][簽到] 簽到資料:", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	// 判斷錯誤訊息
	var freeData map[string]interface{}
	json.Unmarshal(jsonData, &freeData)
	if freeData["error"] != nil {
		var errorMessage SignError
		json.Unmarshal(jsonData, &errorMessage)
		log.Fatalln("[簽到] 簽到失敗：", jsonData)
		return nil, nil, errorMessage
	}

	var signInfo SignInfo
	json.Unmarshal(jsonData, &signInfo)
	return &signInfo, nil, nil
}
