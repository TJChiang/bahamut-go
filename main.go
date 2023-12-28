package main

import (
	"bahamut/internal/config"
	"bahamut/pkg/login"
)

func main() {
	v, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	login.Login(v)
}
