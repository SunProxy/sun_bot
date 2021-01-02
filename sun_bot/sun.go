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
	"fmt"
	"github.com/Jviguy/SpeedyCmds"
	"github.com/Jviguy/SpeedyCmds/command/commandMap"
	"github.com/bwmarrin/discordgo"
	"github.com/jszwedko/go-circleci"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var client circleci.Client

func Start() error {
	cfg := LoadConfig()
	client = circleci.Client{Token: cfg.CircleCI.Token}
	dg, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		return err
	}
	cmdMap := commandMap.New()
	dg.StateEnabled = true
	h := SpeedyCmds.New(dg, cmdMap, true, "sun@root")
	RegisterCommands(h.GetCommandMap())
	dg.AddHandler(onMessage)
	dg.AddHandler(onJoin)
	dg.AddHandler(onLeave)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMembers)
	err = dg.Open()
	if err != nil {
		return err
	}
	fmt.Println("sun_bot is now online!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	return dg.Close()
}

func onMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if len(msg.Mentions) > 4 {
		_ = session.GuildMemberDeleteWithReason(msg.GuildID, msg.Author.ID, "You have been auto kicked for pinging too many members at once!")
		return
	}
	// Groups -> 1: repo; 2: file; 3: line-from; 4: line-to
	var re = regexp.MustCompile(`http(?:s|)://github\.com/(.*?/.*?/)blob/(.*?/.*?)#L([0-9]+)-?L?([0-9]+)?`)
	match := re.FindAllStringSubmatch(msg.Content, -1)
	if len(match) == 0 {
		return
	}
	resp, err := http.Get("https://raw.githubusercontent.com/" + match[0][1] + match[0][2])
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")
	lfrom, _ := strconv.Atoi(match[0][3])
	lto, _ := strconv.Atoi(match[0][4])
	send := ""
	if lto != -1 {
		for lfrom-1 < lto {
			lfrom++
			send += fmt.Sprintf("%v > %v", lfrom, lines[lfrom]+"\n")
		}
	}
	fmt.Println(send)
	_, _ = session.ChannelMessageSend(msg.ChannelID, send)
}

func onJoin(session *discordgo.Session, member *discordgo.GuildMemberAdd) {
	_ = session.GuildMemberRoleAdd(member.GuildID, member.User.ID, "787812497796104192")
	_, _ = session.ChannelMessageSend("790271339322671135", "Welcome to SunProxy's discord "+member.Mention())
	em := &discordgo.MessageEmbed{Title: "Welcome!"}
	em.Description = "This discord provides support and also development ideas and recommendations on sun proxy!"
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	build, _ := client.ListRecentBuildsForProject("SunProxy", "sun", "master", "", 1, 0)
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Latest Build: ",
		Value:  fmt.Sprintf("[%v](%s), [Download Here](%s)", build[0].BuildNum, build[0].BuildURL, GenerateArtifactUrl(build[0].BuildNum)),
		Inline: false,
	})
	em.Fields = append(em.Fields, &discordgo.MessageEmbedField{Name: "Channels: ",
		Value: fmt.Sprintf("For Support head to <#787786998995615756>\n" +
			"For General Chatting head to <#787786916840865803>\n" +
			"For Development help or api usage head to <#787786964853850144>\n" +
			"For keeping up with builds and commits head to <#787794545152491570>\n" +
			"For using me or other bots head to <#789710042299236392>"),
		Inline: false,
	})
	em.Color = rgb(150, 200, 255).ToInteger()
	em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: build[0].VCSURL + "/blob/master/SunProxy.png?raw=true"}
	_, _ = session.ChannelMessageSendEmbed("790271339322671135", em)
}

func onLeave(session *discordgo.Session, member *discordgo.GuildMemberRemove) {
	em := &discordgo.MessageEmbed{Title: "I caught a moon lover on his way out!"}
	em.Description = "His name was " + member.Mention() + " || " + member.User.Username + "!"
	em.Color = rgb(255, 0, 0).ToInteger()
	em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: "https://i.redd.it/iysl5f5vrxe31.jpg"}
	_, _ = session.ChannelMessageSendEmbed("790300380001075231", em)
}
