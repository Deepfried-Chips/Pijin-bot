package sets

import (
	"github.com/bwmarrin/discordgo"
	_ "github.com/traefik/yaegi/interp"
)

type JSONMaster struct {
	Commands []*JSONCommand `json:"Commands"`
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type JSONCommand struct {
	Name               string
	ApplicationCommand *discordgo.ApplicationCommand `json:"ApplicationCommand"`
	Handler            CommandHandler
}
