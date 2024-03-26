package main

import (
	bot "go-spooty/Bot"
	spotify "go-spooty/Spotify"
	logg "go-spooty/log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Discord struct {
		Name  string `yaml:"name"`
		App   string `yaml:"app_id"`
		Guild string `yaml:"guild_id"`
		Token string `yaml:"token"`
	} `yaml:"discord"`
	Spotify struct {
		Client  string `yaml:"client_id"`
		Key     string `yaml:"client_secret"`
		Refresh string `yaml:"refresh_token"`
	} `yaml:"spotify"`
	Endpoint struct {
		Token string `yaml:"token"`
		State string `yaml:"playback_state"`
	} `yaml:"endpoint"`
}

func main() {
	f, err := os.Open("config.yml")
	if err != nil {
		logg.SystemLogger.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logg.SystemLogger.Fatal(err)
	}

	// Start the Bot
	logg.SystemLogger.Println("Bot running")

	spotify.Client = cfg.Spotify.Client
	spotify.Key = cfg.Spotify.Key
	spotify.Refresh = cfg.Spotify.Refresh
	spotify.Token = cfg.Endpoint.Token
	spotify.State = cfg.Endpoint.State

	bot.BotName = cfg.Discord.Name
	bot.BotApp = cfg.Discord.App
	bot.BotGuild = cfg.Discord.Guild
	bot.BotToken = cfg.Discord.Token
	bot.Run()
}
