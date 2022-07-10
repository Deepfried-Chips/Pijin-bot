package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// BotToken Flags
var (
	BotToken      = flag.String("token", "", "Bot token")
	UnsplashToken = flag.String("unsplash-token", "", "Unsplash token")
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

type ExtendedCommand struct {
	applicationcommand *discordgo.ApplicationCommand
}

type UnsplashRandom struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Color       string `json:"color"`
	Description string `json:"description"`
	URLs        struct {
		Raw     string `json:"raw"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Full    string `json:"full"`
		Thumb   string `json:"thumb"`
		S3      string `json:"small-s3"`
	}
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	User struct {
		ID        string `json:"id"`
		UpdatedAt string `json:"updated_at"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		Links     struct {
			Self      string `json:"self"`
			HTML      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
		} `json:"links"`
	} `json:"user"`
}

func init() { flag.Parse() }

func UnsplashImageFromApi(query string) *UnsplashRandom {
	resp, err := http.Get(fmt.Sprintf("https://api.unsplash.com/photos/random?query=%v&client_id=%v", query, *UnsplashToken))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	var unsplash *UnsplashRandom

	err = json.NewDecoder(resp.Body).Decode(&unsplash)
	if err != nil {
		var invalid *UnsplashRandom
		return invalid
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return unsplash
}

func init() {
	*BotToken = ReadTokenFromFile("token.txt")
	if *BotToken != "" {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		*BotToken = os.Getenv("TOKEN")
	}
	*UnsplashToken = ReadTokenFromFile("unsplash-token.txt")
	if *UnsplashToken != "" {
		log.Println("Unsplash token read from file")
	} else {
		log.Println("Unsplash token not read from file, fetching from env")
		*UnsplashToken = os.Getenv("UNSPLASH_TOKEN")
	}
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func ReadTokenFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
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
	return string(buf)
}

func init() {
	flag.Parse()
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
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
		{
			applicationcommand: &discordgo.ApplicationCommand{
				Name:        "randompigeon",
				Description: "uses unsplash to get a random image of a pigeon",
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
		"randompigeon": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var (
				query = "pigeon"
				err   error
			)

			var unsplash = UnsplashImageFromApi(query)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Random Pigeon (Click to view image on unsplash)",
							URL:   fmt.Sprintf("%s?utm_source=Pijin-Bot&utm_medium=referral", unsplash.Links.HTML),
							Image: &discordgo.MessageEmbedImage{
								URL:      unsplash.URLs.Raw,
								ProxyURL: fmt.Sprintf("%s?utm_source=Pijin-Bot&utm_medium=referral", unsplash.Links.HTML),
								Width:    unsplash.Width,
								Height:   unsplash.Height,
							},
							Description: unsplash.Description,
							Author: &discordgo.MessageEmbedAuthor{
								Name: fmt.Sprintf("Photo by %s, taken from unsplash", unsplash.User.Name),
								URL:  unsplash.User.Links.Self,
							},
						},
					},
				},
			})
			if err != nil {
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
