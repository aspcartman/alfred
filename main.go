package main

import (
	"fmt"
	"github.com/aspcartman/alfred/env"
)

func main() {
	fmt.Println(env.TelegramToken())
}