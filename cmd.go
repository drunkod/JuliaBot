package main

import (
	"main/modules"
	"os"
	"os/exec"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func initFunc(c *telegram.Client) {
	c.UpdatesGetState()

	if LOAD_MODULES {
		c.AddMessageHandler("/mz", modules.YtSongDL)
		c.AddMessageHandler("/sh", modules.ShellHandle)
		c.AddMessageHandler("/ul", modules.UploadHandle)

		c.AddMessageHandler("/start", modules.StartHandle)
		c.AddMessageHandler("/system", modules.GatherSystemInfo)
		c.AddMessageHandler("/info", modules.UserHandle)
		c.AddMessageHandler("/json", modules.JsonHandle)
		c.AddMessageHandler("/ping", modules.PingHandle)
		c.AddMessageHandler("?eval", modules.EvalHandle)
	}
}

func dowloadYTVid(url string) {
	// if not url: perform a search
	if !strings.HasPrefix(url, "http") {
		ytsearch := exec.Command("yt-dlp", "-o", "maxresdefault.webm", "ytsearch:"+url)
		ytsearch.Run()
	} else {
		// yt-dlp -f best -o maxresdefault.%(ext)s https://www.youtube.com/watch?v=S7tYeUBgGHU
		ytdl := exec.Command("yt-dlp", "-o", "maxresdefault.webm", url)
		ytdl.Run()
	}

	// convert to mkv
	ffmpeg := exec.Command("ffmpeg", "-i", "maxresdefault.webm", "-c:v", "copy", "-c:a", "copy", "maxresdefault.mp4", "-y")
	ffmpeg.Run()

	os.Remove("maxresdefault.webm")
}

func dowloadYTAud(url string) {
	// if not url: perform a search
	// embed the thumbnail, title, artist, and album art

	if !strings.HasPrefix(url, "http") {
		ytsearch := exec.Command("yt-dlp", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "--embed-metadata", "-o", "maxresdefault.mp3", "ytsearch:"+url)
		ytsearch.Run()
	} else {
		// yt-dlp -f best -o maxresdefault.%(ext)s https://www.youtube.com/watch?v=S7tYeUBgGHU
		ytdl := exec.Command("yt-dlp", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "--embed-metadata", "-o", "maxresdefault.mp3", url)
		ytdl.Run()
	}
}
