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

func songArtists() string {
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
	logg.MessageLogger.Printf("%s: %s", m.Author.Username, m.Content)

	if m.Author.ID == sess.State.User.ID {
		return
	}

	if m.Content == "1" {
		sess.ChannelMessageSend(m.ChannelID, "2")
	}
}

func embedMessage() []*discordgo.MessageEmbed {
	json.Unmarshal([]byte(spotify.GetPlaybackState()), &song)

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
				Value:  songArtists(),
				Inline: false,
			},
		},
	},
	}

	return embed
}

func Commands(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "song":
		err := sess.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: embedMessage(),
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
		if err != nil {
			logg.SystemLogger.Println(err)
		}
	}
}
