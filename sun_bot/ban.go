package sun_bot

import (
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

type Ban struct {
}

func (b Ban) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	if ctx.GetAuthor().ID != "789633382933069906" {
		return nil
	}
	var reason string
	var days int
	for k, arg := range ctx.GetArgs() {
		if arg == "reason" || arg == "-r" {
			reason = strings.Join(ctx.GetArgs()[k:], " ")
		} else if arg == "-d" {
			tmp, err := strconv.Atoi(ctx.GetArgs()[k+1])
			if err != nil {
				return err
			}
			days = tmp
		}
	}
	for _, mention := range ctx.GetMessage().Mentions {
		_ = session.GuildBanCreateWithReason(ctx.GetGuild().ID, mention.ID, reason, days)
		ch, _ := session.UserChannelCreate(mention.ID)
		if reason == "" {
			_, _ = session.ChannelMessageSend(ch.ID, "You have been banned from Sun Discord for: No Reason Provided")
			continue
		}
		_, _ = session.ChannelMessageSend(ch.ID, "You have been banned from Sun Discord for: "+reason)
	}
	em := &discordgo.MessageEmbed{}
	em.Title = "Success!"
	em.Description = "The following members have been banned for: " + reason + "\n"
	for _, mention := range ctx.GetMessage().Mentions {
		em.Description += mention.Mention()
	}
	em.Color = rgb(200, 100, 240).ToInteger()
	_, _ = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
	return nil
}

func (b Ban) GetHelp() HelpMsg {
	return HelpMsg{Usage: "sun@root ban @member", Description: "Bans a member from the sun discord!"}
}
