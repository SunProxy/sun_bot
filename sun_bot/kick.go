package sun_bot

import (
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Kick struct {
	name string
}

func (k Kick) GetName() string {
	return k.name
}

func (k Kick) Setname(newname string) {
	k.name = newname
}

func (k Kick) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	if ctx.GetAuthor().ID != "789633382933069906" {
		return nil
	}
	var reason = ""
	for k, arg := range ctx.GetArgs() {
		if arg == "reason" || arg == "r" {
			reason = strings.Join(ctx.GetArgs()[k:], " ")
		}
	}
	for _, mention := range ctx.GetMessage().Mentions {
		_ = session.GuildMemberDeleteWithReason(ctx.GetGuild().ID, mention.ID, reason)
		ch, _ := session.UserChannelCreate(mention.ID)
		if reason == "" {
			_, _ = session.ChannelMessageSend(ch.ID, "You have been kicked from Sun Discord for: No Reason Provided")
			continue
		}
		_, _ = session.ChannelMessageSend(ch.ID, "You have been kicked from Sun Discord for: "+reason)
	}
	em := &discordgo.MessageEmbed{}
	em.Title = "Success!"
	em.Description = "The following members have been kicked for: " + reason + "\n"
	for _, mention := range ctx.GetMessage().Mentions {
		em.Description += mention.Mention()
	}
	em.Color = rgb(200, 100, 240).ToInteger()
	_, _ = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
	return nil
}

func (k Kick) GetHelp() HelpMsg {
	return HelpMsg{Usage: "sun@root kick @member", Description: "Kicks a member from the sun discord!"}
}
