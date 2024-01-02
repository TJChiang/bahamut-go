package main

import (
	libBrowser "bahamut/internal/browser"
	"bahamut/internal/config"
	"bahamut/internal/container"
	"bahamut/pkg/login"
	"bahamut/pkg/sign"
	"net/http"
)

func main() {
	conf, err := config.InitConfig("./configs/config.yaml")
	if err != nil {
		panic(err)
	}

	con := container.Register(&http.Client{}, conf)

	browser, context, err := libBrowser.Launch(con.Config().GetBrowserConfig())
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
	_, err = login.Login(con, page)
	if err != nil {
		panic(err)
	}

	// 簽到
	_, err = sign.Sign(con, page)
	if err != nil {
		panic(err)
	}
}
