package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Initialize Discord bot

	discord, err := discordgo.New("Bot <bot token>")
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Add message handler
	discord.AddHandler(messageHandler)

	// Open Discord connection
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	// Wait for interruption signal to gracefully close connection
	fmt.Println("Bot is now running. Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord connection
	discord.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message starts with the command prefix
	if !strings.HasPrefix(m.Content, "!roll") {
		return
	}

	// Parse the dice roll command
	cmd := strings.TrimSpace(strings.TrimPrefix(m.Content, "!roll"))
	if cmd == "" {
		return
	}

	// Parse the number of dice and the number of sides per dice
	var numDice, numSides int
	_, err := fmt.Sscanf(cmd, "%dd%d", &numDice, &numSides)
	if err != nil || numDice <= 0 || numSides <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid dice roll command. Usage: !roll <num>d<sides>")
		return
	}

	// Simulate the dice roll
	total := 0
	for i := 0; i < numDice; i++ {
		total += rand.Intn(numSides) + 1
	}

	// Send the roll result to the channel
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s rolled %dd%d and got %d!", m.Author.Mention(), numDice, numSides, total))
}
