package browser

import (
	"strings"

	"github.com/playwright-community/playwright-go"
)

func Goto(page playwright.Page, location string) (playwright.Response, error) {
	var u string
	switch strings.ToLower(location) {
	case "login":
		u = "https://user.gamer.com.tw/login.php"
	case "user":
		u = "https://home.gamer.com.tw/homeindex.php?owner=<owner>"
	case "home":
	default:
		u = "https://www.gamer.com.tw/"
	}

	return page.Goto(u)
}
