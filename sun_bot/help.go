/**
      ___           ___           ___
     /  /\         /__/\         /__/\
    /  /:/_        \  \:\        \  \:\
   /  /:/ /\        \  \:\        \  \:\
  /  /:/ /::\   ___  \  \:\   _____\__\:\
 /__/:/ /:/\:\ /__/\  \__\:\ /__/::::::::\
 \  \:\/:/~/:/ \  \:\ /  /:/ \  \:\~~\~~\/
  \  \::/ /:/   \  \:\  /:/   \  \:\  ~~~
   \__\/ /:/     \  \:\/:/     \  \:\
     /__/:/       \  \::/       \  \:\
     \__\/         \__\/         \__\/

MIT License

Copyright (c) 2020 Jviguy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package sun_bot

import (
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/Jviguy/SpeedyCmds/utils"
	"github.com/bwmarrin/discordgo"
)

var sunimg = "https://cdn.discordapp.com/attachments/787786916840865803/789586385090772992/sun.png"

type Help struct {
	name string
}

func (h Help) GetHelp() HelpMsg {
	return HelpMsg{"~help <COMMAND>","Returns the help message for the given command!"}
}

func (h Help) GetName() string {
	return h.name
}

func (h Help) Setname(newname string) {
	h.name = newname
}

func (h Help) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	if len(ctx.GetArgs()) > 0 {
		if cmd, ok := cmds.GetCommands()[ctx.GetArgs()[0]]; ok {
			em := &discordgo.MessageEmbed{}
			em.Title = ctx.GetArgs()[0]
			em.Description = "Usage: " + cmd.(SunCommand).GetHelp().Usage + "\n" +
				"Description: " + cmd.(SunCommand).GetHelp().Description
			em.Color = rgb(150,200,255).ToInteger()
			em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: sunimg}
			_,_ = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
			return nil
		}
		em := & discordgo.MessageEmbed{Title: "Unknown Command"}
		em.Description = "You might have meant " + utils.FindClosest(ctx.GetArgs()[0], utils.GetAllKeysCommands(cmds.GetCommands()))
		em.Color = rgb(255, 0, 0).ToInteger()
		em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: sunimg}
		_,_ = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
		return nil
	}
	em := &discordgo.MessageEmbed{}
	em.Title = "Help"
	em.Description = "Usage: " + h.GetHelp().Usage + "\n" +
		"Description: " + h.GetHelp().Description
	em.Color = rgb(150,200,255).ToInteger()
	em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: sunimg}
	_,_ = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
	return nil
}

