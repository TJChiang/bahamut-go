package main

import (
	libBrowser "bahamut/internal/browser"
	"bahamut/internal/config"
	"bahamut/pkg/login"
	"bahamut/pkg/sign"
)

func main() {
	c, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	browser, context, err := libBrowser.Launch(c.GetBrowserConfig())
	if err != nil {
		panic(err)
	}
	defer browser.Close()
	defer context.Close()
	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}
	defer page.Close()

	// 登入
	_, err = login.Login(c.GetModulesConfig().Login, page)
	if err != nil {
		panic(err)
	}

	// 簽到
	_, err = sign.Sign(c.GetBrowserConfig(), page)
	if err != nil {
		panic(err)
	}
}
