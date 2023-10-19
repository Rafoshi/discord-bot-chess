package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	Avatar    string `json:"avatar"`
	URL       string `json:"url"`
	UserName  string `json:"username"`
	Name      string `json:"name"`
	Followers int    `json:"followers"`
	Country   string `json:"country"`
	League    string `json:"league"`
}

const prefix string = "chess"

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Print(err)
	}
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")

		if args[0] != prefix || len(args) < 2 {
			printHelp(s, m)
			return
		}

		if args[1] == "user" {
			player := printNames(args[2])
			author := discordgo.MessageEmbedAuthor{
				Name: player.Name,
				URL:  player.URL,
			}

			image := discordgo.MessageEmbedImage{
				URL: player.Avatar,
			}

			country := strings.TrimPrefix(player.Country, "https://api.chess.com/pub/country/")

			footer := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Seguidores: %d", player.Followers),
			}

			embed := discordgo.MessageEmbed{
				Title:       player.League,
				Description: country,
				Footer:      &footer,
				Author:      &author,
				Image:       &image,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		} else if args[1] == "help" {
			s.ChannelMessageSend(m.ChannelID, "Comandos: \n ```chess user <username>``` ```chess help```")
		}
	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Print("Bot is running")

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}()
}

func printHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Use ```chess help``` para ver os comandos")
}

func printNames(playername string) Player {
	url := "https://api.chess.com/pub/player/" + playername

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Status not OK")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var player Player
	err = json.Unmarshal(body, &player)
	if err != nil {
		panic(err)
	}
	return player
}
