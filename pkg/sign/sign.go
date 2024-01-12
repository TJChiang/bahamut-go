package sign

import (
	"bahamut/internal/browser"
	"bahamut/internal/container"
	"encoding/json"
	"log"

	"github.com/playwright-community/playwright-go"
)

type SignData struct {
	Data SignInfo `json:"data"`
}

type SignInfo struct {
	Days          int `json:"days"`
	Signin        int `json:"signin"`
	FinishedAd    int `json:"finishedAd"`
	PrjSigninDays int `json:"prjSigninDays"`
}

func Sign(con *container.Container, page playwright.Page) (*SignInfo, error) {
	log.Println("[簽到] 開始執行簽到")
	homeRes, err := browser.Goto(page, browser.Home)
	if err != nil {
		return nil, err
	}

	log.Println("[簽到] 主頁狀態: ", homeRes.Status())

	info, err := getSignStatus(con, page)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func getSignStatus(con *container.Container, page playwright.Page) (*SignInfo, error) {
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
		return nil, err
	}

	if con.Config().GetModulesConfig().Sign.Debug {
		log.Println("[Debug][簽到] 簽到資料:", data)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 判斷錯誤訊息
	var freeData map[string]interface{}
	if err = json.Unmarshal(jsonData, &freeData); err != nil {
		return nil, err
	}

	if freeData["error"] != nil {
		var errorMessage SignError
		json.Unmarshal(jsonData, &errorMessage)
		log.Fatalln("[簽到] 簽到失敗：", jsonData)
		return nil, errorMessage
	}

	var signData SignData
	json.Unmarshal(jsonData, &signData)

	if con.Config().GetModulesConfig().Sign.Debug {
		log.Println("[Debug][簽到] json資料:", data)
	}

	log.Printf("[簽到] 已連續簽到 %d 天 \n", signData.Data.Days)
	if signData.Data.Signin != 1 {
		if err = page.Locator("a#signin-btn").Click(); err != nil {
			log.Println("[簽到] 手動觸發失敗：", err)
			return nil, err
		}
		log.Println("[簽到] 手動簽到成功！")
	} else {
		log.Println("[簽到] 今日已簽到！")
	}

	return &signData.Data, nil
}
