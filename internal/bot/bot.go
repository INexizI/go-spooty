package bot

import (
	"encoding/json"
	"fmt"
	logg "go-spooty/internal/log"
	spotify "go-spooty/internal/spotify"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotName  string
	BotApp   string
	BotGuild string
	BotToken string
	song     spotify.Song
)

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

func Run() {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		logg.SystemLogger.Fatal(err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(BotApp, BotGuild, []*discordgo.ApplicationCommand{
		{
			Name:        "song",
			Description: "Get current Spotify song",
		},
		{
			Name:        "first",
			Description: "Showcase of a first slash command",
		},
		{
			Name:        "second",
			Description: "Showcase of a second slash command",
		},
	})
	if err != nil {
		logg.SystemLogger.Println(err)
	}

	session.AddHandler(Messages)
	session.AddHandler(Commands)
	session.Open()
	defer session.Close()

	fmt.Println("Bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logg.SystemLogger.Println("Bot shutdown")
}

func Messages(sess *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == sess.State.User.ID {
		return
	} else {
		logg.MessageLogger.Printf("%s: %s", m.Author.Username, m.Content)
	}

	if m.Content == "1" {
		sess.ChannelMessageSend(m.ChannelID, "2")
	}
}

func EmbedMessageSpotify() []*discordgo.MessageEmbed {
	json.Unmarshal([]byte(spotify.GetPlaybackState()), &song)

	if song.Playing {
		embed := []*discordgo.MessageEmbed{{
			URL:   song.Item.External.URL,
			Type:  discordgo.EmbedTypeRich,
			Title: "Spotify",
			// Description:
			// Timestamp:
			Color: 5763719,
			// Footer:
			Image: &discordgo.MessageEmbedImage{
				URL:    song.Item.Album.Images[0].URL,
				Width:  song.Item.Album.Images[0].Width,
				Height: song.Item.Album.Images[0].Height,
			},
			// Thumbnail:
			// Video:
			// Provider:
			// Author:
			Fields: []*discordgo.MessageEmbedField{
				{
					// Name:   "Title",
					Value:  song.Item.Name,
					Inline: true,
				},
				{
					// Name:   "Artist",
					Value:  SongArtists(),
					Inline: false,
				},
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

func Commands(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "song":
		err := sess.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: 4, // discordgo.InteractionResponseChannelMessageWithSource = 4
				Data: &discordgo.InteractionResponseData{
					Embeds: EmbedMessageSpotify(),
				},
			},
		)
		if err != nil {
			logg.SystemLogger.Println(err)
		}
	case "first":
		err := sess.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					TTS:     true,
					Content: "First command!",
				},
			},
		)
		logg.CommandLogger.Printf("first command")
		if err != nil {
			logg.SystemLogger.Println(err)
		}
	case "second":
		err := sess.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Second command!",
				},
			},
		)
		logg.CommandLogger.Printf("second command")
		if err != nil {
			logg.SystemLogger.Println(err)
		}
	default:
		err := sess.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
			},
		)
		if err != nil {
			sess.FollowupMessageCreate(
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

// func Roles(sess *discordgo.Session, ra *discordgo.MessageReactionAdd, rr *discordgo.MessageReactionRemove) {
// 	if ra.Emoji.Name == "" {
// 		sess.GuildMemberRoleAdd(ra.GuildID, ra.UserID, "roleID")
// 		sess.ChannelMessageSend(ra.ChannelID, fmt.Sprintf("%v has been added to %v", ra.UserID, ra.Emoji.Name))

// 		sess.GuildMemberRoleRemove(rr.GuildID, rr.UserID, "roleID")
// 		sess.ChannelMessageSend(rr.ChannelID, fmt.Sprintf("%v has been removed to %v", rr.UserID, rr.Emoji.Name))
// 	}
// }
