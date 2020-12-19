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
	RegisterCommands(h.GetCommandHandler())
	dg.AddHandler(onMessage)
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
