package spotify

import (
	"context"

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

func (sc *SpotifyClient) GetPlaylist(id string) (*spotify.FullPlaylist, error) {
	var spotifyID spotify.ID = spotify.ID(id)
	playlist, err := sc.client.GetPlaylist(spotifyID)
	// tracks, err := sc.client.GetPlaylistTracks(spotifyID)

	if err != nil {
		return nil, err
	}

	return playlist, nil
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