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
	session, err := discordgo.New("Bot " + "MTE2NDI5NzA5MTgzNzg3NDM1Ng.Gz0bAr.WNGZsbj5IYqR92dJrloekpxyaf-L-E7fx5n93w")
	if err != nil {
		fmt.Print(err)
	}
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")

		if args[0] != prefix || len(args) < 2 {
			return
		}

		if args[1] == "user" {
            player := printNames(args[2])
            message := fmt.Sprintf("Nome: %s\nLiga: %s\nPaís: %s", player.Name, player.League, player.Country)
			s.ChannelMessageSend(m.ChannelID, message)
		}
	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	fmt.Print("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	//printNames()
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
