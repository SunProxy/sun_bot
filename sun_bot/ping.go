package sun_bot

import (
	"fmt"
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
)

type Ping struct {
	name string
}

func (p Ping) GetName() string {
	return p.name
}

func (p Ping) Setname(newname string) {
	p.name = newname
}

func (p Ping) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	em := &discordgo.MessageEmbed{Title: "API Ping!"}
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Websocket Latency: ",
		Value:  fmt.Sprintf("%vms", session.HeartbeatLatency().Milliseconds()),
		Inline: false,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  "Fetching.....",
		Inline: false,
	})
	msg, _ := session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
	em.Fields[1] = &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  "Fetching......",
		Inline: false,
	}
	_, _ = session.ChannelMessageEditEmbed(ctx.GetChannel().ID, msg.ID, em)
	org, _ := msg.Timestamp.Parse()
	edit, _ := msg.Timestamp.Parse()
	em.Fields[1] = &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  fmt.Sprintf("%vms", org.Sub(edit).Milliseconds()),
		Inline: false,
	}
	_, _ = session.ChannelMessageEditEmbed(ctx.GetChannel().ID, msg.ID, em)
	return nil
}

func (p Ping) GetHelp() HelpMsg {
	return HelpMsg{"sun@root ping", "Returns the ping of the discord api websocket and the http requests."}
}
