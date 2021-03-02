package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyClient struct {
	config *clientcredentials.Config;
	client *spotify.Client
}

func (sc *SpotifyClient) init() error {
	token, err := sc.config.Token(context.Background())
	if err != nil {
		return err
	}

	client := spotify.Authenticator{}.NewClient(token)
	sc.client = &client;
	return nil;
}

func (sc *SpotifyClient) GetPlaylist(id string) (sp []spotify.PlaylistTrack, err error) {
	var spotifyID spotify.ID = spotify.ID(id)
	// playlist, err := sc.client.GetPlaylist(spotifyID)
	tracks, err := sc.client.GetPlaylistTracks(spotifyID)

	if err != nil {
		return nil, err
	}

	sp = tracks.Tracks

	// paginate
	for page := 1; ; page++ {
		err = sc.client.NextPage(tracks)

		if err == spotify.ErrNoMorePages {
			break;
		}

		if err != nil {
			break;
		}

		sp = append(sp, tracks.Tracks...)
	}
	fmt.Printf("Playlist has %d tracks\n", len(sp))

	return sp, nil
}

func NewSpotifyClient(clientID string, clientSecret string) (*SpotifyClient, error) {
	spotifyClient := &SpotifyClient{
		config: &clientcredentials.Config{
			ClientID: clientID,
			ClientSecret: clientSecret,
			TokenURL: spotify.TokenURL,
		},
	}

	err := spotifyClient.init()

	if err != nil {
		return nil, err
	}

	return spotifyClient, nil;
}