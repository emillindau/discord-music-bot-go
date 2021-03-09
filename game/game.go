package game

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/bwmarrin/discordgo"
	client "github.com/emillindau/discord-music-bot-go/discord"
	"github.com/emillindau/discord-music-bot-go/spotify"
	"github.com/emillindau/discord-music-bot-go/utils"
)

const maxAllowedDistance = 1
const numberOfSongs = 4

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	spotifyClient *spotify.SpotifyClient
	discordClient *client.DiscordClient
	songs []*Song
	currentSong *Song
	state string
	onEnd chan bool
	users map[string]*User
	playedSongs int
}

func (g *Game) init(playlistID string) error {
	tracks, err := g.spotifyClient.GetPlaylist(playlistID)
	g.state = "initialized"
	g.onEnd = make(chan bool)

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

func (g *Game) handleGuess(guess string) bool {
	distance := levenshtein.ComputeDistance(strings.ToLower(guess), strings.ToLower(g.currentSong.name))
	fmt.Printf("The distance between %s and %s is %d.\n", guess, g.currentSong.name, distance)

	if (distance <= maxAllowedDistance) {
		return true
	}
	return false
}

func (g *Game) nextSong() error {
	g.state = "next"
	// handle game end
	if g.playedSongs >= numberOfSongs {
		g.state = "finished"

		var points int
		var winner string
		for _, user := range g.users {
			if user.points > points {
				winner = user.name
				points = user.points
			}
		}

		g.discordClient.SendMessage("Game is done! The winner was " + winner + " with " + strconv.Itoa(points) + " points!")
		g.discordClient.Stop()
		return nil
	}

	song := g.getRandomSong()
	g.currentSong = song
	fmt.Println("trying to play: ", song)

	filepath, err := utils.DownloadFile(song.id, song.downloadUrl)
	if err != nil {
		fmt.Println("could not download song")
		return err
	}

	// wait until playing next
	time.Sleep(3 * time.Second) 

	go g.discordClient.Play(filepath, g.onEnd)

	go func() {
		isEnd := <-g.onEnd
		if isEnd {
			g.discordClient.SendMessage("Nobody guessed right :( The correct song was " + g.currentSong.artist + " - " + g.currentSong.name)
			g.nextSong()
		}
	}()

	// increment played songs
	g.playedSongs++

	g.state = "started"

	return nil
}

func (g *Game) Start() {
	g.state = "started"
	err := g.nextSong()

	if err != nil {
		fmt.Println("could not start song")
	}
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
		users: make(map[string]*User),
	}

	err := game.init(p)

	if err != nil {
		return nil, err
	}

	go func() {
		game.discordClient.ListenForMessage(func(u *discordgo.User, m string) {
			if (game.state != "started" && m == "start") {
				game.discordClient.SendMessage("Sit tight! Starting soon")
				game.Start()
			} else if (game.state == "started") {
				currentUser, ok := game.users[u.ID]
				if !ok {
					 currentUser = &User{
						id: u.ID,
						name: u.Username,
						points: 0,
					}
					game.users[u.ID] = currentUser
				}

				res := game.handleGuess(m)

				if (res) {
					// Just give 10p for now
					currentUser.points += 10
					correctMsg := "That was indeed right! " + u.Username 
					game.discordClient.SendMessage(correctMsg)
					game.nextSong()
				}
			}
		})
	}()

	// fmt.Println("playlist", playlist)
	return game, nil
}