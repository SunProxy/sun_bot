package sun_bot

import (
	"fmt"
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Ping struct {
}

func (p Ping) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	loc, _ := time.LoadLocation("UTC")
	em := &discordgo.MessageEmbed{Title: "API Ping!"}
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	em.Color = rgb(211, 20, 124).ToInteger()
	em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: sunimg}
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Websocket Latency: ",
		Value:  fmt.Sprintf("%vms", session.HeartbeatLatency().Milliseconds()),
		Inline: false,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  "Fetching.....",
		Inline: false,
	})
	msg, _ := session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
	org, _ := msg.Timestamp.Parse()
	em.Fields[1] = &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  "Fetching......",
		Inline: false,
	}
	msg, _ = session.ChannelMessageEditEmbed(ctx.GetChannel().ID, msg.ID, em)
	edit, _ := msg.EditedTimestamp.Parse()
	em.Fields[1] = &discordgo.MessageEmbedField{Name: "HTTP Latency: ",
		Value:  fmt.Sprintf("%+vms", edit.In(loc).Sub(org.In(loc)).Milliseconds()),
		Inline: false,
	}
	_, _ = session.ChannelMessageEditEmbed(ctx.GetChannel().ID, msg.ID, em)
	return nil
}

func (p Ping) GetHelp() HelpMsg {
	return HelpMsg{"sun@root ping", "Returns the ping of the discord api websocket and the http requests."}
}
