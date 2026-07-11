package main

import (
	"fmt"

	"github.com/thalesraymond/gator/internal/config"
)

func main() {
	fmt.Println("Hello World")

	cfg, err := config.ReadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.SetUser("thales")

	fmt.Println(cfg)
}
