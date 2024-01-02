package login

import (
	"bahamut/internal/browser"
	"bahamut/internal/container"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"text/template"

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
func Login(con *container.Container, page playwright.Page) (bool, error) {
	params := con.Config().Modules.Login
	res, err := requestLogin(
		con.HttpClient(),
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

	bahaCookies := getBahaCookies(res.Cookies())
	if params.Debug {
		log.Println("[Debug][登入] 登入資料:", loginData)
		log.Println("[Debug][登入] Response Cookies:", res.Cookies())
		log.Println("[Debug][登入] Cookies: ", bahaCookies)
	}

	if err = setHttpClientCookie(con, bahaCookies); err != nil {
		return false, err
	}

	if err = handleScriptCookies(bahaCookies, page); err != nil {
		return false, err
	}

	return true, nil
}

func handleScriptCookies(cookies *BahaCookies, page playwright.Page) error {
	if cookies.BahaRune == nil || cookies.BahaEnur == nil {
		return errors.New("登入 cookie 不存在")
	}

	browser.Goto(page, browser.Home)
	context := page.Context()

	// main.go 所在的位置為根目錄
	t, err := template.ParseFiles("./pkg/login/add_cookies.js")
	if err != nil {
		return err
	}
	var js bytes.Buffer
	t.Execute(&js, cookies)
	if err != nil {
		return err
	}

	// 設定起始腳本，將 cookie 存入，讓後續模組 fetch request 可以取得 cookie
	content := js.String()
	context.AddInitScript(playwright.Script{
		Content: &content,
	})

	log.Println("[登入] Script 成功載入 Cookie")

	return nil
}

// 取得巴哈 Response 的 Cookie
func getBahaCookies(cookies []*http.Cookie) *BahaCookies {
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

func setHttpClientCookie(con *container.Container, cookies *BahaCookies) error {
	cookieUrl, err := url.Parse("https://api.gamer.com.tw/mobile_app/user/v3/do_login.php")
	if err != nil {
		return err
	}
	h := con.HttpClient()
	h.Jar, err = cookiejar.New(nil)
	if err != nil {
		return err
	}

	var jar []*http.Cookie
	if cookies.BahaID != nil {
		jar = append(jar, cookies.BahaID)
	}
	if cookies.BahaRune != nil {
		jar = append(jar, cookies.BahaRune)
	}
	if cookies.BahaEnur != nil {
		jar = append(jar, cookies.BahaEnur)
	}
	if cookies.BahaHashID != nil {
		jar = append(jar, cookies.BahaHashID)
	}
	if cookies.MB_BahaID != nil {
		jar = append(jar, cookies.MB_BahaID)
	}
	if cookies.MB_BahaRune != nil {
		jar = append(jar, cookies.MB_BahaRune)
	}

	h.Jar.SetCookies(cookieUrl, jar)

	log.Println("[登入] HTTP Client 成功載入 cookie")
	return nil
}

// 請求登入
func requestLogin(h *http.Client, username string, password string) (*http.Response, error) {
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
	return h.Do(req)
}
