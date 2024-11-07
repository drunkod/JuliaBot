package main

import (
	"main/modules"

	"github.com/amarnathcjd/gogram/telegram"
)

func FilterOwner(m *telegram.NewMessage) bool {
	if m.SenderID() == ownerId {
		return true
	}
	m.Reply("You are not allowed to use this command")
	return false
}

func initFunc(c *telegram.Client) {
	c.UpdatesGetState()

	if LOAD_MODULES {
		c.On("message:/mz", modules.YtSongDL)
		c.On("message:/sh", modules.ShellHandle, telegram.FilterFunc(FilterOwner))
		c.On("message:/ul", modules.UploadHandle, telegram.FilterFunc(FilterOwner))

		c.On("message:/start", modules.StartHandle)
		c.On("message:/help", modules.HelpHandle)
		c.On("message:/system", modules.GatherSystemInfo)
		c.On("message:/info", modules.UserHandle)
		c.On("message:/json", modules.JsonHandle)
		c.On("message:/ping", modules.PingHandle)
		c.On("message:?eval", modules.EvalHandle, telegram.FilterFunc(FilterOwner))

		c.On("message:/file", modules.SendFileByIDHandle)
		c.On("message:/fid", modules.GetFileIDHandle)
		c.On("message:/dl", modules.DownloadHandle, telegram.FilterFunc(FilterOwner))
	}
}
