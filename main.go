package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tg "github.com/amarnathcjd/gogram/telegram"
	dotenv "github.com/joho/godotenv"
)

const LOAD_MODULES = true

var startTimeStamp = time.Now().Unix()

func main() {
	dotenv.Load()

	appId, _ := strconv.Atoi(os.Getenv("APP_ID"))
	client, err := tg.NewClient(tg.ClientConfig{
		AppID:   int32(appId),
		AppHash: os.Getenv("APP_HASH"),
	})

	if err != nil {
		panic(err)
	}

	client.Conn()
	client.LoginBot(os.Getenv("BOT_TOKEN"))

	initFunc(client)
	me, err := client.GetMe()

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Bot started as  @%s", me.Username))
	client.Idle()
}
