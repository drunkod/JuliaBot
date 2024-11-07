package modules

import (
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func UploadHandle(m *telegram.NewMessage) error {
	filename := m.Args()
	if filename == "" {
		m.Reply("No filename provided")
		return nil
	}

	spoiler := false
	if strings.Contains(filename, "-s") {
		spoiler = true
		filename = strings.ReplaceAll(filename, "-s", "")
	}

	msg, _ := m.Reply("Uploading...")
	uploadStartTimestamp := time.Now()

	var pm *telegram.ProgressManager

	if _, err := m.RespondMedia(filename, telegram.MediaOptions{
		Spoiler: spoiler,
		ProgressCallback: func(total, curr int64) {
			if pm == nil {
				pm = telegram.NewProgressManager(total, 5)
			}
			if pm.ShouldEdit() {
				m.Client.EditMessage(m.ChatID(), msg.ID, pm.GetStats(curr))
			}
		},
	}); err != nil {
		msg.Edit("Error: " + err.Error())
		return nil
	} else {
		msg.Edit("Uploaded " + filename + " in <code>" + time.Since(uploadStartTimestamp).String() + "</code>")
	}

	return nil
}
