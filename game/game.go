package game

import (
	"fmt"
	"os"
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
	tracks, err := g.spotifyClient.GetPlaylist(playlistID)

	if err != nil {
		return err
	}

	// TODO: Fetch all songs + youtube

	for _, t := range tracks {
		pu := t.Track.PreviewURL

		if (pu != "") {
			ar := []string{}
			for _, a := range t.Track.Artists {
				ar = append(ar, a.Name)
			}

			artists := strings.Join(ar, ", ")
			song := newSong(artists, t.Track.Name, t.Track.PreviewURL, t.Track.ID.String())
			g.songs = append(g.songs, song)
		}
	}

	fmt.Printf("Filtered out to %d tracks\n", len(g.songs))

	sg := g.songs[:100]

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, s := range sg {
			downloadFile(s.id, s.downloadUrl)
		}
	}()
	wg.Wait()

	// fmt.Println("songs", len(g.songs))

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