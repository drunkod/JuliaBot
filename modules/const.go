package modules

import (
	"math/rand"
	"os"
	"strconv"
)

var random_emoji_list = []string{
	"❤️",
	"👍",
	"😍",
	"😂",
	"😊",
	"🔥",
	"😘",
	"💕",
}

func getRandomEmoticon() string {
	return random_emoji_list[rand.Intn(len(random_emoji_list))]
}

var OWNER_ID, _ = strconv.Atoi(os.Getenv("OWNER_ID"))
