package main

import (
	libBrowser "bahamut/internal/browser"
	"bahamut/internal/config"
	"bahamut/pkg/login"
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
	login.Login(c.GetModulesConfig().Login, &page)
}
