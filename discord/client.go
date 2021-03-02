package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	token string
	Session *discordgo.Session
}

func NewDiscordClient(token string) (*DiscordClient, error) {
	client := &DiscordClient{
		token: token,
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

func (dc *DiscordClient) ListenForMessage() {
	fmt.Println("Starting to listen for messages")
	dc.Session.AddHandler(handleMessage)
}

func (dc *DiscordClient) Exit() {
	dc.Session.Close()
}

func (dc *DiscordClient) Play() {

}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
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

		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			return
		}

		// Look up voice channel
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, guild.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing ", err)
				}
				return
			}
		}
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func loadSound(filePath string) ([][]byte, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	var opuslen int16
	var buffer = make([][]byte, 0)
	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// eof
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from file");
			return nil, err
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)

		if err != nil {
			fmt.Println("Error reading pcm")
			return nil, err
		}

		buffer = append(buffer, inBuf)
	}
}

func playSound(s *discordgo.Session, guildID string, channelID string) (err error) {
	// join voice channel
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)

	// Send buffer bla bla
	buffer, err := loadSound("temp/1.dca")
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	return nil;
}