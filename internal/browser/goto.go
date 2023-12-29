package browser

import (
	"github.com/playwright-community/playwright-go"
)

type Page string

const (
	Home  Page = "home"
	Login Page = "login"
	User  Page = "user"
)

func Goto(page playwright.Page, location Page) (playwright.Response, error) {
	var u string
	switch {
	case location == "login":
		u = "https://user.gamer.com.tw/login.php"
	case location == "user":
		u = "https://home.gamer.com.tw/homeindex.php?owner=<owner>"
	case location == "home":
	default:
		u = "https://www.gamer.com.tw/"
	}

	return page.Goto(u)
}
