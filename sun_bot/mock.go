package sun_bot

import (
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
)

type Mock struct {
}

func (m Mock) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	msg := strings.Join(ctx.GetArgs(), " ")
	var send string
	for _, c := range msg {
		if rand.Intn(100) < 50 && c > 96 && c < 123 {
			send += string(c ^ 0x20)
			continue
		}
		send += string(c)
	}
	_, _ = session.ChannelMessageSend(ctx.GetChannel().ID, send)
	return nil
}

func (m Mock) GetHelp() HelpMsg {
	return HelpMsg{
		Usage:       "sun@root mock <MSG>",
		Description: "Mocks a message or saying as inputted!",
	}
}
