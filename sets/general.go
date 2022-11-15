package sets

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func newGenModule() *[]JSONCommand {
	return GenModule
}

var GenModule = &[]JSONCommand{
	{
		Name: "ping",
		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Ping",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Ping: %s", s.HeartbeatLatency()),
				},
			})
			if err != nil {
				return
			}
		},
	},
}
