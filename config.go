package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
)

type handlerfunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type JSONMaster struct {
	Commands []*JSONCommand
}

type JSONCommand struct {
	name        string
	command     *discordgo.ApplicationCommand
	permissions *discordgo.ApplicationCommandPermissionsList
}

func LoadCommands(s *discordgo.Session, file string) {
	var master JSONMaster
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	err := json.Unmarshal([]byte(file), &master)
	if err != nil {
		log.Panicf("Error loading commands: %v", err)
	}
	for i, command := range master.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", command.command)
		if err != nil {
			return
		}
		err = s.ApplicationCommandPermissionsEdit(s.State.User.ID, "", command.command.ID, command.permissions)
		if err != nil {
			return
		}
		registeredCommands[i] = cmd
	}
}

func ExampleStructWithFuncSavingTest() {
	master := JSONMaster{
		Commands: []*JSONCommand{
			{
				name: "test",
			},
		},
	}
	marshal, err := json.Marshal(master)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("test.json", marshal, 0644)
	if err != nil {
		return
	}
}
