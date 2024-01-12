package main

import (
	libBrowser "bahamut/internal/browser"
	"bahamut/internal/config"
	"bahamut/internal/container"
	linenotify "bahamut/pkg/line_notify"
	"bahamut/pkg/login"
	"bahamut/pkg/sign"
	"fmt"
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
	signInfo, err := sign.Sign(con, page)
	if err != nil {
		e := fmt.Sprintf("\n[簽到失敗] %s", err.Error())
		linenotify.Notify(con, e)
		panic(err)
	}

	n := fmt.Sprintf("\n[簽到成功] 已簽到 %d 天", signInfo.Days)
	// Line Notify
	if err = linenotify.Notify(con, n); err != nil {
		panic(err)
	}
}
