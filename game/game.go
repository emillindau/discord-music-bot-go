package game

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	client "github.com/emillindau/discord-music-bot-go/discord"
	"github.com/emillindau/discord-music-bot-go/spotify"
	"github.com/emillindau/discord-music-bot-go/utils"
)

var messages chan string
func init() {
	rand.Seed(time.Now().UnixNano())
	messages = make(chan string)
}


type Game struct {
	spotifyClient *spotify.SpotifyClient
	discordClient *client.DiscordClient
	songs []*Song
	state string
}

func (g *Game) init(playlistID string) error {
	tracks, err := g.spotifyClient.GetPlaylist(playlistID)
	g.state = "initialized"

	if err != nil {
		return err
	}

	for _, t := range tracks {
		pu := t.Track.PreviewURL

		// filter out without previewurl
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

	return nil
}

func (g *Game) getRandomSong() *Song {
	i := rand.Intn(len(g.songs))
	return g.songs[i]
}

func (g *Game) nextSong() {
	song := g.getRandomSong()
	fmt.Println(song)
	filepath, err := utils.DownloadFile(song.id, song.downloadUrl)
	if err != nil {
		fmt.Println("could not download song")
		return
	}

	end := make(chan bool)
	go g.discordClient.Play(filepath, end)
	isEnd := <-end

	if isEnd {
		g.nextSong()
	}
}

func (g *Game) Start() {
	g.nextSong()
	g.state = "started"
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

	game.discordClient.ListenForMessage(messages)

	go func() {
		for {
			msg := <-messages
			if msg == "start" {
				game.discordClient.SendMessage("Sit tight! Starting soon")
				game.Start()
			} else if msg == "end" {
				break;
			} else {
				if game.state == "started" {

				}
			}
		}
	}()

	// fmt.Println("playlist", playlist)
	return game, nil
}