package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/emillindau/discord-music-bot-go/config"
	client "github.com/emillindau/discord-music-bot-go/discord"
	"github.com/emillindau/discord-music-bot-go/game"
	"github.com/emillindau/discord-music-bot-go/spotify"
)

func main() {
	cfg, err := config.GetConfig()
	discordClient, err := client.NewDiscordClient(cfg["token"])
	spotifyClient, err := spotify.NewSpotifyClient(cfg["clientId"], cfg["clientSecret"])

	_, err = game.NewGame(spotifyClient, discordClient, cfg["playlist"])

	if err != nil {
		panic("could not initialize")
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	game.Exit()
	discordClient.Exit()
}
