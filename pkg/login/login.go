package login

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
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

func Login(viper *viper.Viper) (*LoginData, *BahaCookies, error) {
	res, err := requestLogin(
		viper.GetString("login.username"),
		viper.GetString("login.password"),
	)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var body map[string]interface{}
	if err = json.Unmarshal(rawData, &body); err != nil {
		return nil, nil, err
	}
	if body["message"] != nil {
		log.Fatalln("登入失敗： ", body["message"])
		var errorMessage ErrorMessage
		json.Unmarshal(rawData, &errorMessage)
		return nil, nil, &errorMessage
	}

	var loginData LoginData
	if err = json.Unmarshal(rawData, &loginData); err != nil {
		return nil, nil, err
	}

	log.Println("登入成功： ", loginData.Userid)

	bahaCookies := handleCookies(res.Cookies())
	if viper.GetBool("login.debug") {
		fmt.Println("data: ", loginData)
		fmt.Println("response cookies: ", res.Cookies())
		fmt.Println("Baha cookies: ", bahaCookies)
	}
	return &loginData, bahaCookies, nil
}

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
