package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	_ "github.com/traefik/yaegi/interp"
	"io/ioutil"
	"log"
)

type JSONMaster struct {
	Commands []*JSONCommand `json:"Commands"`
}

type JSONCommand struct {
	name               string
	applicationcommand *discordgo.ApplicationCommand `json:"ApplicationCommand"`
	handler            string
}

func LoadCommands(s *discordgo.Session, file string) []*JSONCommand {
	var master JSONMaster
	fileData, err := ioutil.ReadFile(file)
	err = json.Unmarshal(fileData, &master)
	if err != nil {
		log.Panicf("Error loading commands: %v", err)
	}
	registeredCommands := make([]*JSONCommand, len(master.Commands))
	for i, command := range master.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", command.applicationcommand)
		if err != nil {
			return nil
		}
		command.applicationcommand = cmd
		registeredCommands[i] = command
	}
	return registeredCommands
}

func GenerateExampleCommandsJson() {
	var master JSONMaster
	master.Commands = []*JSONCommand{
		{
			name: "ping",
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping",
			},
		},
		{
			name: "echo",
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "echo",
				Description: "make the bot relay a message to a channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "The channel to send the message to",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "content",
						Description: "The content of the message",
						Required:    true,
					},
				},
			},
		},
	}
	fileData, err := json.Marshal(master)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("data/commands.json", fileData, 0644)
	if err != nil {
		return
	}
}
