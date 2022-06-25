package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

// BotToken Flags
var (
	BotToken = flag.String("token", "", "Bot token")
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

type ExtendedCommand struct {
	applicationcommand *discordgo.ApplicationCommand
}

func init() { flag.Parse() }

func init() {
	if ReadTokenFromFile("token.txt") {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*BotToken = os.Getenv("TOKEN")
	}
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func ReadTokenFromFile(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {

	}
	buf = buf[:n]
	*BotToken = string(buf)
	return true
}

func init() {
	flag.Parse()
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func DMUser(s *discordgo.Session, m *discordgo.Message, msg string) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		if err != nil {
			return
		}
		return
	}
	// Then we send the message through the channel we created.
	_, err = s.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		fmt.Println("error sending DM message:", err)
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
		if err != nil {
			return
		}
	}
}

var (
	commands = []*ExtendedCommand{
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping",
			},
		},
		{
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
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "~~why do u hurt me~~ Pong!",
				},
			})
			if err != nil {
				return
			}
		},
		"echo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			if i.Interaction.Member.Permissions&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages {
				options := i.ApplicationCommandData().Options

				optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
				for _, opt := range options {
					optionMap[opt.Name] = opt
				}

				var (
					channelobj     *discordgo.Channel
					messagecontent string
				)

				// Get the value from the option map.
				// When the option exists, ok = true
				if option, ok := optionMap["channel"]; ok {
					// Option values must be type asserted from interface{}.
					// Discordgo provides utility functions to make this simple.
					channelobj = option.ChannelValue(s)
				}

				if opt, ok := optionMap["content"]; ok {
					messagecontent = opt.StringValue()
				}

				_, err := s.ChannelMessageSendComplex(channelobj.ID, &discordgo.MessageSend{
					Content: messagecontent,
				})
				if err != nil {
					return
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Echo sent to %s:\n```%s```", channelobj.Mention(), messagecontent),
					},
				})
				if err != nil {
					return
				}
			} else {
				return
			}

		},
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v.applicationcommand)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.applicationcommand.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {

		}
	}(s)
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Graceful shutdown")

}
