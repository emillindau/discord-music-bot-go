package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/emillindau/discord-music-bot-go/utils"
)

type DiscordClient struct {
	token string
	Session *discordgo.Session
	Voice *discordgo.VoiceConnection
	Channel *discordgo.Channel
	onStop chan bool
	isPlaying bool
}

func NewDiscordClient(token string) (*DiscordClient, error) {
	client := &DiscordClient{
		token: token,
		onStop: make(chan bool),
	}

	err := client.init()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (dc *DiscordClient) init() error {
	fmt.Printf("token %s\n", dc.token)
	discord, err := discordgo.New("Bot " + dc.token)

	if err != nil {
		return errors.New("could not initialize")
	}

	dc.Session = discord

	// discord.AddHandler(handleMessage)
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates


	wsErr := discord.Open()

	if wsErr != nil {
		fmt.Println("websocket error", wsErr)
		return errors.New("could not open connection")
	}

	return nil
}

// func (dc *DiscordClient) ListenForMessage(c chan<- string) {
// 	fmt.Println("Starting to listen for messages")
// 	dc.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
// 		handleMessage(s, m, dc, c)
// 	})
// }

type MessageCallback func(*discordgo.User, string)
func (dc *DiscordClient) ListenForMessage(cb MessageCallback) {
	dc.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handleMessage(s, m, dc, cb)
	})
}

func (dc *DiscordClient) Exit() {
	dc.Session.Close()
}

func (dc *DiscordClient) Stop() {
	if (dc.isPlaying) {
		dc.onStop <- true
	}
}

func (dc *DiscordClient) Play(path string, end chan<- bool) {
	dc.Stop()
	dc.Voice.Speaking(true)

	// Send buffer bla bla
	buffer, err := utils.LoadSound(path)
	if err != nil {
		dc.Voice.Speaking(false)
		fmt.Println("could not play song")
		return
	}

	dc.isPlaying = true
	for _, buff := range buffer {
		select {
			case <-dc.onStop:
				fmt.Println("Received stop")
				return
			case dc.Voice.OpusSend <- buff:
		}
	}

	end <- true

	dc.isPlaying = false
	dc.Voice.Speaking(false)
}

func (dc *DiscordClient) SendMessage(message string) {
	fmt.Println("trying to send message")
	_, err := dc.Session.ChannelMessageSend(dc.Channel.ID, message)

	if err != nil {
		fmt.Println(err)
	}
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, dc *DiscordClient, cb MessageCallback) {
	// Ignore messages by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!quiz") {
		// Find channel
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}
		dc.Channel = channel

		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			return
		}

		// Look up voice channel
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := joinChannel(s, guild.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error joining channel ", err)
					return
				}
				dc.Voice = vc
			}
		}

	} else {
		cb(m.Author, m.Content)
	}

}

func joinChannel(s *discordgo.Session, guildID string, channelID string) (*discordgo.VoiceConnection, error) {
	// join voice channel
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	return vc, nil
}