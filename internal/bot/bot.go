package bot

import (
	"encoding/json"
	"fmt"
	logg "go-spooty/internal/log"
	spotify "go-spooty/internal/spotify"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Direct struct {
	OriginChannelId string
	Name            string
	Login           string
}

type Quote struct {
	Text  string
	Title string
	Name  string
}

var (
	BotName  string
	BotApp   string
	BotGuild string
	BotToken string
	song     spotify.Song
	DM       map[string]Direct = map[string]Direct{}
)

func CommandList() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "song",
			Description: "Get current Spotify song",
		},
		{
			Name:        "quote",
			Description: "Get random quote",
		},
		{
			Name:        "first",
			Description: "Showcase of a first slash command",
		},
		{
			Name:        "second",
			Description: "Showcase of a second slash command",
		},
	}
}

func (d *Direct) ToDirectMessageUser() discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title: "Your Name & Login:",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Name",
				Value: d.Name,
			},
			{
				Name:  "Login",
				Value: d.Login,
			},
		},
	}
}

func (d *Direct) ToDirectMessageGuild() discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title: "New User:",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Name",
				Value: d.Name,
			},
		},
	}
}

func Messages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	} else {
		logg.MessageLogger.Printf("%s: %s", m.Author.Username, m.Content)
	}

	if m.Content == "token" {
		DirectMessages(s, m)
	}

	if m.GuildID == "" {
		answers, ok := DM[m.ChannelID]
		if !ok {
			return
		}

		if answers.Name == "" {
			answers.Name = m.Content
			s.ChannelMessageSend(m.ChannelID, "Login?")
			DM[m.ChannelID] = answers
			return
		} else {
			answers.Login = m.Content
			embedUser := answers.ToDirectMessageUser()
			embedGuild := answers.ToDirectMessageGuild()
			s.ChannelMessageSendEmbed(m.ChannelID, &embedUser)
			s.ChannelMessageSendEmbed(answers.OriginChannelId, &embedGuild)
			delete(DM, m.ChannelID)
		}
	}
}

func SongArtists() string {
	json.Unmarshal([]byte(spotify.GetPlaybackState()), &song)
	artist := ""
	for i, art := range song.Item.Artists {
		switch {
		case i == 0:
			artist += art.Name
		case i > 0:
			artist = fmt.Sprintf("%s, %s", artist, art.Name)
		}
	}

	return artist
}

func EmbedMessageSpotify() []*discordgo.MessageEmbed {
	json.Unmarshal([]byte(spotify.GetPlaybackState()), &song)

	if song.Playing {
		embed := []*discordgo.MessageEmbed{{
			URL:   song.Item.External.URL,
			Type:  discordgo.EmbedTypeRich,
			Title: song.Item.Name,
			// Description:
			// Timestamp:
			Color: 5763719,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: "https://storage.googleapis.com/pr-newsroom-wp/1/2023/05/Spotify_Primary_Logo_RGB_Green-300x300.png",
				Text:    "Listening on Spotify",
			},
			Image: &discordgo.MessageEmbedImage{
				URL:    song.Item.Album.Images[0].URL,
				Width:  song.Item.Album.Images[0].Width,
				Height: song.Item.Album.Images[0].Height,
			},
			// Thumbnail:
			// Video:
			// Provider:
			Author: &discordgo.MessageEmbedAuthor{
				Name: SongArtists(),
			},
		},
		}
		return embed
	} else {
		embed := []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Something went wrong!",
			Description: "Please try again later!",
			Color:       15548997,
		},
		}
		return embed
	}
}

func Commands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "song":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: 4, // discordgo.InteractionResponseChannelMessageWithSource = 4
				Data: &discordgo.InteractionResponseData{
					Embeds: EmbedMessageSpotify(),
				},
			},
		)
		logg.CommandLogger.Println("song command")
		if err != nil {
			logg.SystemLogger.Println(err)
		}
	case "quote":
		file, err := os.ReadFile("internal/bot/quotes.json")
		if err != nil {
			logg.SystemLogger.Fatalln("Error when opening file: ", err)
		}

		var data []Quote
		err = json.Unmarshal(file, &data)
		if err != nil {
			logg.SystemLogger.Fatalln("Error during Unmarshal(): ", err)
		}

		selection := rand.Intn(len(data))

		err = s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: 4, // discordgo.InteractionResponseChannelMessageWithSource = 4
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{{ // Must be 256 or fewer in length.
						Title: data[selection].Text,
						Author: &discordgo.MessageEmbedAuthor{
							Name: data[selection].Name + " (" + data[selection].Title + ")",
						},
						Color: 2123412,
					}},
				},
			},
		)
		logg.CommandLogger.Println("quote command")
		if err != nil {
			logg.SystemLogger.Panic(err)
		}
	case "first":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					TTS:     true,
					Content: "First command!",
				},
			},
		)
		logg.CommandLogger.Println("first command")
		if err != nil {
			logg.SystemLogger.Panic(err)
		}
	case "second":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Second command!",
				},
			},
		)
		logg.CommandLogger.Println("second command")
		if err != nil {
			logg.SystemLogger.Panic(err)
		}
	default:
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
			},
		)
		if err != nil {
			s.FollowupMessageCreate(
				i.Interaction,
				true,
				&discordgo.WebhookParams{
					Content: "Something went wrong",
				},
			)
		}
		return
	}
}

// func Roles(s *discordgo.Session, ra *discordgo.MessageReactionAdd, rr *discordgo.MessageReactionRemove) {
// 	if ra.Emoji.Name == "" {
// 		s.GuildMemberRoleAdd(ra.GuildID, ra.UserID, "roleID")
// 		s.ChannelMessageSend(ra.ChannelID, fmt.Sprintf("%v has been added to %v", ra.UserID, ra.Emoji.Name))

// 		s.GuildMemberRoleRemove(rr.GuildID, rr.UserID, "roleID")
// 		s.ChannelMessageSend(rr.ChannelID, fmt.Sprintf("%v has been removed to %v", rr.UserID, rr.Emoji.Name))
// 	}
// }

func DirectMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		logg.SystemLogger.Panic(err)
	}

	if _, ok := DM[channel.ID]; !ok {
		DM[channel.ID] = Direct{
			OriginChannelId: m.ChannelID,
			Name:            "",
			Login:           "",
		}
		s.ChannelMessageSend(channel.ID, "Hello from Spooty Bot!")
		s.ChannelMessageSend(channel.ID, "Please enter your name")
	} else {
		s.ChannelMessageSend(channel.ID, "Uh... Hello... 🧐")
	}
}

func Run() {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		logg.SystemLogger.Fatal(err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(BotApp, BotGuild, CommandList())
	if err != nil {
		logg.SystemLogger.Println(err)
	}

	session.AddHandler(Messages)
	session.AddHandler(Commands)
	// session.AddHandler(Roles)
	session.Open()
	defer session.Close()

	fmt.Println("Bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logg.SystemLogger.Println("Bot shutdown")
}
