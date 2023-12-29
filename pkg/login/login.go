package login

import (
	"bahamut/internal/browser"
	"bahamut/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type LoginData struct {
	Success  bool   `json:"success"`
	Userid   string `json:"userid"`
	Nickname string `json:"nickname"`
	Gold     int    `json:"gold"`
	Gp       int    `json:"gp"`
	Avatar   string `json:"avatar"`
	Avatar_s string `json:"avatar_s"`
	Lv       int    `json:"lv"`
}

type BahaCookies struct {
	BahaID      *http.Cookie
	BahaRune    *http.Cookie
	BahaEnur    *http.Cookie
	BahaHashID  *http.Cookie
	MB_BahaID   *http.Cookie
	MB_BahaRune *http.Cookie
}

// 登入巴哈，瀏覽器載入 Cookie
func Login(params *config.ModulesLogin, page playwright.Page) (bool, error) {
	res, err := requestLogin(
		params.Username,
		params.Password,
	)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var body map[string]interface{}
	if err = json.Unmarshal(rawData, &body); err != nil {
		return false, err
	}
	if body["message"] != nil {
		log.Fatalln("[登入] 登入失敗： ", body["message"])
		var errorMessage ErrorMessage
		json.Unmarshal(rawData, &errorMessage)
		return false, &errorMessage
	}

	var loginData LoginData
	if err = json.Unmarshal(rawData, &loginData); err != nil {
		return false, err
	}

	log.Println("[登入] 登入成功： ", loginData.Userid)

	bahaCookies := handleCookies(res.Cookies())
	if params.Debug {
		fmt.Println("data: ", loginData)
		fmt.Println("response cookies: ", res.Cookies())
		fmt.Println("Baha cookies: ", bahaCookies)
	}

	return goToHomePage(bahaCookies, page), nil
}

func goToHomePage(cookies *BahaCookies, page playwright.Page) bool {
	if cookies.BahaRune == nil || cookies.BahaEnur == nil {
		return false
	}
	browser.Goto(page, browser.Home)
	context := page.Context()
	context.AddCookies([]playwright.OptionalCookie{
		{
			Name:   cookies.BahaID.Name,
			Value:  cookies.BahaID.Value,
			Path:   &cookies.BahaID.Path,
			Domain: &cookies.BahaID.Domain,
		},
		{
			Name:   cookies.BahaRune.Name,
			Value:  cookies.BahaRune.Value,
			Path:   &cookies.BahaRune.Path,
			Domain: &cookies.BahaRune.Domain,
		},
		{
			Name:   cookies.BahaEnur.Name,
			Value:  cookies.BahaEnur.Value,
			Path:   &cookies.BahaEnur.Path,
			Domain: &cookies.BahaEnur.Domain,
		},
	})
	browser.Goto(page, browser.Home)
	log.Println("[登入] 成功載入 Cookie")
	return true
}

// 取得巴哈 Response 的 Cookie
func handleCookies(cookies []*http.Cookie) *BahaCookies {
	var baha BahaCookies
	for _, c := range cookies {
		switch c.Name {
		case "BAHARUNE":
			baha.BahaRune = c
		case "BAHAENUR":
			baha.BahaEnur = c
		case "BAHAID":
			baha.BahaID = c
		case "BAHAHASHID":
			baha.BahaHashID = c
		case "MB_BAHAID":
			baha.MB_BahaID = c
		case "MB_BAHARUNE":
			baha.MB_BahaRune = c
		}
	}

	return &baha
}

// 請求登入
func requestLogin(username string, password string) (*http.Response, error) {
	reqBody := url.Values{}
	reqBody.Set("uid", username)
	reqBody.Set("passwd", password)
	reqBody.Set("vcode", "6666")
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.gamer.com.tw/mobile_app/user/v3/do_login.php",
		strings.NewReader(reqBody.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ckAPP_VCODE=6666")
	return http.DefaultClient.Do(req)
}
