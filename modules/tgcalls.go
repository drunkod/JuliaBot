package modules

import (
    "fmt"
    "os/exec"
    "strconv"
    "time"

    "main/libs"
    "github.com/amarnathcjd/gogram/telegram"
)

var ntg *libs.Client

func init() {
    ntg = libs.NTgCalls() // Initialize NTgCalls library
}

// Join a group call and start streaming audio
func JoinCall(m *telegram.NewMessage) error {
    group := m.Args() // Get group name from command arguments
    if group == "" {
        return m.Reply("Please provide the group to join the call.")
    }

    url := "default_audio.mp3" // Default audio file to stream
    if len(m.Args()) > 1 {     // If an audio file is specified in the command
        url = m.Args()[1]
    }

    media := libs.MediaDescription{
        Audio: &libs.AudioDescription{
            InputMode:     libs.InputModeShell,
            SampleRate:    128000,
            BitsPerSample: 16,
            ChannelCount:  2,
            Input:         fmt.Sprintf("ffmpeg -i %s -loglevel panic -f s16le -ac 2 -ar 128k pipe:1", url), // ffmpeg command to convert audio to s16le format and pipe it to stdout
        },
    }

    client := m.Client
    me, _ := client.GetMe() // Get the current bot user for JoinAs

    call, err := client.GetGroupCall(group) // Get the group call object
    if err != nil {
        return m.Reply("Error fetching group call: " + err.Error())
    }

    rawChannel, err := client.GetSendablePeer(group)
    if err != nil {
        return m.Reply("Error fetching group peer: " + err.Error())
    }
    channel := rawChannel.(*telegram.InputPeerChannel)

    jsonParams, err := ntg.CreateCall(channel.ChannelID, media) // Create call object with media description
    if err != nil {
        return m.Reply("Error creating call: " + err.Error())
    }

    callResRaw, err := client.PhoneJoinGroupCall(
        &telegram.PhoneJoinGroupCallParams{
            Muted:        false,
            VideoStopped: true, // false for video call
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
        return m.Reply("Error joining group call: " + err.Error())
    }

    callRes := callResRaw.(*telegram.UpdatesObj)
    for _, update := range callRes.Updates {
        switch u := update.(type) {
        case *telegram.UpdateGroupCallConnection: // Wait for connection params
            _ = ntg.Connect(channel.ChannelID, u.Params.Data)
        }
    }

    return m.Reply(fmt.Sprintf("Joined group call in %s and streaming audio.", group))
}

// Leave the group call
func LeaveCall(m *telegram.NewMessage) error {
    group := m.Args() // Get group name from command arguments
    if group == "" {
        return m.Reply("Please provide the group name to leave the call.")
    }

    rawChannel, err := m.Client.GetSendablePeer(group)
    if err != nil {
        return m.Reply("Error fetching group peer: " + err.Error())
    }
    channel := rawChannel.(*telegram.InputPeerChannel)

    // End the call using the NTgCalls library
    err = ntg.DestroyCall(channel.ChannelID)
    if err != nil {
        return m.Reply("Error leaving group call: " + err.Error())
    }

    return m.Reply(fmt.Sprintf("Left group call in %s.", group))
}

// Play a specific audio file in a group call
func PlayAudio(m *telegram.NewMessage) error {
    group := m.Args() // Get group name from command arguments
    if group == "" {
        return m.Reply("Please provide the group name to play audio.")
    }

    url := "default_audio.mp3" // Default audio file
    if len(m.Args()) > 1 {
        url = m.Args()[1]
    }

    media := libs.MediaDescription{
        Audio: &libs.AudioDescription{
            InputMode:     libs.InputModeShell,
            SampleRate:    128000,
            BitsPerSample: 16,
            ChannelCount:  2,
            Input:         fmt.Sprintf("ffmpeg -i %s -loglevel panic -f s16le -ac 2 -ar 128k pipe:1", url), // ffmpeg command to convert audio to s16le format and pipe it to stdout
        },
    }

    rawChannel, err := m.Client.GetSendablePeer(group)
    if err != nil {
        return m.Reply("Error fetching group peer: " + err.Error())
    }
    channel := rawChannel.(*telegram.InputPeerChannel)

    err = ntg.UpdateCall(channel.ChannelID, media) // Update the media being streamed in the current call
    if err != nil {
        return m.Reply("Error playing audio in group call: " + err.Error())
    }

    return m.Reply(fmt.Sprintf("Playing audio in group call in %s.", group))
}

func init() {
    // Register the commands to make them available in the Telegram bot
    Mods.AddModule("TgCalls", `<b>Group Call Commands:</b>

<code>/joincall [group_name]</code> - Join a group call and stream audio
<code>/leavecall [group_name]</code> - Leave the group call
<code>/play [group_name] [audio_file]</code> - Play a specific audio file in the group call`)

    // Add command handlers
    Mods.AddHandler("/joincall", JoinCall)
    Mods.AddHandler("/leavecall", LeaveCall)
    Mods.AddHandler("/play", PlayAudio)
}
