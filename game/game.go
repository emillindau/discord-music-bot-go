package game

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	client "github.com/emillindau/discord-music-bot-go/discord"
	"github.com/emillindau/discord-music-bot-go/spotify"
)

type Game struct {
	spotifyClient *spotify.SpotifyClient
	discordClient *client.DiscordClient
	songs []*Song
}

func (g *Game) init(playlistID string) error {
	playlist, err := g.spotifyClient.GetPlaylist(playlistID)

	if err != nil {
		return err
	}

	// TODO: Fetch all songs + youtube

	for _, t := range playlist.Tracks.Tracks {
		pu := t.Track.PreviewURL

		if (pu != "") {
			ar := []string{}
			for _, a := range t.Track.Artists {
				ar = append(ar, a.Name)
			}

			artists := strings.Join(ar, ", ")
			song := newSong(artists, t.Track.Name, t.Track.PreviewURL)
			g.songs = append(g.songs, song)
		}
	}

	var wg sync.WaitGroup
	for i, s := range g.songs {
		wg.Add(1)
		go downloadFile(strconv.Itoa(i), s.url, &wg)
	}
	wg.Wait()

	fmt.Println("songs", len(g.songs))

	return nil
}

func Exit() {
	fmt.Println("Trying to remove files")
	os.RemoveAll("temp/")
	os.Mkdir("temp", 0755)
}

func NewGame(sc *spotify.SpotifyClient, dc *client.DiscordClient, p string) (*Game, error) {
	game := &Game{
		spotifyClient: sc,
		discordClient: dc,
	}

	err := game.init(p)

	if err != nil {
		return nil, err
	}

	
	game.discordClient.ListenForMessage()

	// fmt.Println("playlist", playlist)
	return game, nil
}