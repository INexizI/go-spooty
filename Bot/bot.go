package bot

import (
	"encoding/json"
	"fmt"
	spotify "go-spooty/Spotify"
	logg "go-spooty/log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotName  string
	BotToken string
)

func Run() {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		logg.SystemLogger.Fatal(err)
	}

	session.AddHandler(Messages)
	session.Open()
	defer session.Close()

	fmt.Println("Bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logg.SystemLogger.Println("Bot shutdown")
}

func Messages(s *discordgo.Session, m *discordgo.MessageCreate) {
	logg.MessageLogger.Printf("%s: %s", m.Author.Username, m.Content)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "1" {
		s.ChannelMessageSend(m.ChannelID, "2")
	}

	if m.Content == "!song" {
		var song spotify.Song
		json.Unmarshal([]byte(spotify.GetPlaybackState()), &song)

		currentSong := embedMessage(&song)
		s.ChannelMessageSendComplex(m.ChannelID, currentSong)
	}
}

func embedMessage(s *spotify.Song) *discordgo.MessageSend {
	artist := ""
	for i, art := range s.Item.Artists {
		switch {
		case i == 0:
			artist += art.Name
		case i > 0:
			artist = fmt.Sprintf("%s, %s", artist, art.Name)
		}
	}
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			URL:   s.Item.External.URL,
			Type:  discordgo.EmbedTypeRich,
			Title: "Spotify",
			// Description:
			// Timestamp:
			Color: 5763719,
			// Footer:
			Image: &discordgo.MessageEmbedImage{
				URL:    s.Item.Album.Images[0].URL,
				Width:  s.Item.Album.Images[0].Width,
				Height: s.Item.Album.Images[0].Height,
			},
			// Thumbnail:
			// Video:
			// Provider:
			// Author:
			Fields: []*discordgo.MessageEmbedField{
				{
					// Name:   "Title",
					Value:  s.Item.Name,
					Inline: true,
				},
				{
					// Name:   "Artist",
					Value:  artist,
					Inline: false,
				},
			},
		},
		},
	}

	return embed
}
