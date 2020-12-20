package sun_bot

import (
	"fmt"
	"github.com/Jviguy/SpeedyCmds/command/ctx"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

type Build struct {
	name string
}

func (b Build) GetHelp() HelpMsg {
	return HelpMsg{"sun@root build <INT or STRING> <INT or STRING>", "Returns build info on Builds and can rebuild " +
		"/ find artifacts of certain ones!"}
}

func (b Build) GetName() string {
	return b.name
}

func (b Build) Setname(newname string) {
	b.name = newname
}

func (b Build) Execute(ctx ctx.Ctx, session *discordgo.Session) error {
	if len(ctx.GetArgs()) > 0 {
		num, err := strconv.Atoi(ctx.GetArgs()[0])
		if err != nil {
			switch ctx.GetArgs()[0] {
			case "latest":
				builds, _ := client.ListRecentBuildsForProject("SunProxy", "sun", "master", "", 1, 0)
				build, _ := client.GetBuild("SunProxy", "sun", builds[0].BuildNum)
				em := &discordgo.MessageEmbed{Title: "Information on Build " + ctx.GetArgs()[0] + "!"}
				fields := make([]*discordgo.MessageEmbedField, 0)
				for _, step := range build.Steps {
					fields = append(fields, &discordgo.MessageEmbedField{Name: step.Name, Value: step.Actions[0].Status, Inline: true})
				}
				em.Fields = fields
				em.Color = rgb(211, 20, 124).ToInteger()
				_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
				return nil
			case "recent":
				builds, err := client.ListRecentBuildsForProject("SunProxy", "sun", "master", "", -1, 0)
				if err != nil {
					log.Fatal(err)
				}
				em := &discordgo.MessageEmbed{Title: "All Builds On SunProxy/sun!"}
				fields := make([]*discordgo.MessageEmbedField, 0)
				for _, build := range builds {
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:   strconv.Itoa(build.BuildNum) + ":",
						Value:  "Status: " + build.Status + "\nLink: " + build.BuildURL,
						Inline: false})
				}
				em.Fields = fields
				em.Color = rgb(211, 20, 124).ToInteger()
				_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
				return err
			case "status":
				if len(ctx.GetArgs()) > 1 {
					num, err := strconv.Atoi(ctx.GetArgs()[1])
					if err != nil {
						return err
					}
					build, err := client.GetBuild("SunProxy", "sun", num)
					if err != nil {
						return fmt.Errorf("unknown build: %v", num)
					}
					em := &discordgo.MessageEmbed{Title: "Status of Build " + ctx.GetArgs()[1] + "!"}
					fields := make([]*discordgo.MessageEmbedField, 0)
					fields = append(fields, &discordgo.MessageEmbedField{Name: "Status: ", Value: build.Status, Inline: false})
					fields = append(fields, &discordgo.MessageEmbedField{Name: "Build Link: ", Value: build.BuildURL, Inline: false})
					fields = append(fields, &discordgo.MessageEmbedField{Name: "Artifacts: ", Value: GenerateArtifactUrl(num), Inline: false})
					fields = append(fields, &discordgo.MessageEmbedField{Name: "Timing: ",
						Value:  "from " + build.StartTime.String() + " to " + build.StopTime.String(),
						Inline: false})
					fields = append(fields, &discordgo.MessageEmbedField{Name: "Build Starter: ", Value: build.Why, Inline: false})
					for _, step := range build.Steps {
						fields = append(fields, &discordgo.MessageEmbedField{Name: step.Name, Value: step.Actions[0].Status, Inline: true})
					}
					em.Fields = fields
					em.Color = rgb(211, 20, 124).ToInteger()
					em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: build.VCSURL + "/blob/master/SunProxy.png?raw=true"}
					_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
					return nil
				}
			case "artifacts":
				if len(ctx.GetArgs()) > 1 {
					num, err := strconv.Atoi(ctx.GetArgs()[1])
					if err != nil {
						switch ctx.GetArgs()[1] {
						case "latest":
							build, err := client.ListRecentBuildsForProject("SunProxy", "sun", "master", "", 1, 0)
							if err != nil {
								return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
							}
							em := &discordgo.MessageEmbed{Title: "Artifacts of Build " + ctx.GetArgs()[1] + "!"}
							em.Description = "**Artifacts:** [Click Here](" + GenerateArtifactUrl(build[0].BuildNum) + ")"
							em.Color = rgb(211, 20, 124).ToInteger()
							_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
							return nil
						}
						return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
					}
					_, err = client.GetBuild("SunProxy", "sun", num)
					if err != nil {
						return fmt.Errorf("unknown build: %v", num)
					}
					em := &discordgo.MessageEmbed{Title: "Artifacts of Build " + ctx.GetArgs()[1] + "!"}
					em.Description = "**Artifacts:** [Click Here](" + GenerateArtifactUrl(num) + ")"
					em.Color = rgb(211, 20, 124).ToInteger()
					_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
					return nil
				}
			case "rebuild":
				if len(ctx.GetArgs()) > 1 && ctx.GetAuthor().ID == "431853152518668294" {
					num, err := strconv.Atoi(ctx.GetArgs()[1])
					if err != nil {
						switch ctx.GetArgs()[1] {
						case "latest":
							builds, err := client.ListRecentBuildsForProject("SunProxy", "sun", "master", "", 1, 0)
							if err != nil {
								return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
							}
							build, err := client.RetryBuild("SunProxy", "sun", builds[0].BuildNum)
							if err != nil {
								return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
							}
							em := &discordgo.MessageEmbed{Title: fmt.Sprintf("Rebuilding %v!", build.BuildNum)}
							em.Description = "Starting now!"
							em.Color = rgb(211, 20, 124).ToInteger()
							em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: build.VCSURL + "/blob/master/SunProxy.png?raw=true"}
							_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
							return nil
						}
						return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
					}
					build, err := client.RetryBuild("SunProxy", "sun", num)
					if err != nil {
						return fmt.Errorf("unknown build: %v", ctx.GetArgs()[1])
					}
					em := &discordgo.MessageEmbed{Title: fmt.Sprintf("Rebuilding %v!", build.BuildNum)}
					em.Description = "Starting now!"
					em.Color = rgb(211, 20, 124).ToInteger()
					em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: build.VCSURL + "/blob/master/SunProxy.png?raw=true"}
					_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
				}
			}
		}
		build, err := client.GetBuild("SunProxy", "sun", num)
		if err != nil {
			return fmt.Errorf("unknown build: %v", num)
		}
		em := &discordgo.MessageEmbed{Title: "Information on Build " + ctx.GetArgs()[0] + "!"}
		fields := make([]*discordgo.MessageEmbedField, 0)
		for _, step := range build.Steps {
			fields = append(fields, &discordgo.MessageEmbedField{Name: step.Name, Value: step.Actions[0].Status, Inline: true})
		}
		em.Fields = fields
		em.Color = rgb(211, 20, 124).ToInteger()
		_, err = session.ChannelMessageSendEmbed(ctx.GetChannel().ID, em)
		return nil
	}
	return nil
}

func GenerateArtifactUrl(buildNum int) string {
	return fmt.Sprintf("https://app.circleci.com/pipelines/github/SunProxy/sun/jobs/%v/artifacts",
		buildNum)
}
