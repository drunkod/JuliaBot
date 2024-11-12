package modules

/*
#cgo CFLAGS: -I../libs
#cgo LDFLAGS: -L../libs -lntgcalls -Wl,-rpath=../libs
*/
import "C"
import (
	"fmt"

	"main/libs"

	// "flag"

	"github.com/amarnathcjd/gogram/telegram"
	dotenv "github.com/joho/godotenv"
)

func StartStream(filePath string, groupID string) error {
    dotenv.Load(".env")
    
    client, err := telegram.NewClient(telegram.ClientConfig{
        AppID:   AskInputOrEnv[int32]("API_KEY"),
        AppHash: AskInputOrEnv[string]("API_HASH"),
    })

    if err != nil {
        return err
    }

    client.AuthPrompt()

    ntg := libs.NTgCalls()
    defer ntg.Free()

    media := libs.MediaDescription{
        Audio: &libs.AudioDescription{
            InputMode:     libs.InputModeShell,
            SampleRate:    128000,
            BitsPerSample: 16,
            ChannelCount:  2,
            Input:         fmt.Sprintf("ffmpeg -i %s -loglevel panic -f s16le -ac 2 -ar 128k pipe:1", filePath),
        },
    }

    err = joinGroupCall(client, ntg, groupID, media)
    if err != nil {
        return err
    }

    client.Idle()
    return nil
}

func joinGroupCall(client *telegram.Client, ntg *libs.Client, chatId interface{}, media libs.MediaDescription) error {
    me, err := client.GetMe()
    if err != nil {
        return err
    }

    call, err := client.GetGroupCall(chatId)
    if err != nil {
        return err
    }

    rawChannel, err := client.GetSendablePeer(chatId)
    if err != nil {
        return err
    }
    
    channel := rawChannel.(*telegram.InputPeerChannel)
    jsonParams, err := ntg.CreateCall(channel.ChannelID, media)
    if err != nil {
        return err
    }

    callResRaw, err := client.PhoneJoinGroupCall(
        &telegram.PhoneJoinGroupCallParams{
            Muted:        false,
            VideoStopped: true,
            Call:         call,
            Params: &telegram.DataJson{
                Data: jsonParams,
            },
            JoinAs: &telegram.InputPeerUser{
                UserID:     me.ID,
                AccessHash: me.AccessHash,
            },
        },
    )

    if err != nil {
        return err
    }

    callRes := callResRaw.(*telegram.UpdatesObj)
    for _, update := range callRes.Updates {
        switch u := update.(type) {
        case *telegram.UpdateGroupCallConnection:
            if err := ntg.Connect(channel.ChannelID, u.Params.Data); err != nil {
                return err
            }
        }
    }

    return nil
}