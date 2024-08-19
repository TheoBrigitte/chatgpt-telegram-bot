package run

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v3"

	"github.com/TheoBrigitte/chatgpt-telegram-bot/pkg/openai"
)

var (
	Cmd = &cobra.Command{
		Use:   "run",
		Short: "Run Telegram bot",
		RunE:  runner,
	}
)

func runner(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize OpenAI client
	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient := openai.New(openaiApiKey)

	// Initialize Telegram bot
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	pref := tele.Settings{
		Token:  telegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	// Handler for chat conversations
	b.Handle(tele.OnText, func(c tele.Context) error {
		sender := c.Sender()
		text := c.Text()

		response, err := openaiClient.ChatCompletion(ctx, sender.ID, text)
		if err != nil {
			return fmt.Errorf("Response error: %w", err)
		}

		return c.Reply(response)
	})

	// Handler to display chat history
	b.Handle("/history", func(c tele.Context) error {
		sender := c.Sender()

		history := openaiClient.GetHistory(sender.ID)

		var output string
		for _, h := range history {
			output += fmt.Sprintf("%s: %s\n", h.Role, h.Content)
		}

		return c.Reply(fmt.Sprintf("History\n%s", output))
	})

	// Handler to display chat role
	b.Handle("/role", func(c tele.Context) error {
		sender := c.Sender()

		role := openaiClient.GetRole(sender.ID)

		return c.Reply(fmt.Sprintf("Role is: %s", role))
	})

	// Handler to set chat role
	b.Handle("/setrole", func(c tele.Context) error {
		sender := c.Sender()
		role := c.Message().Payload

		if role == "" {
			return c.Reply("Usage: /setrole <role>")
		}

		openaiClient.SetRole(sender.ID, role)

		return c.Reply(fmt.Sprintf("Role set to: %s", role))
	})

	// Handler to clear chat session
	b.Handle("/clear", func(c tele.Context) error {
		sender := c.Sender()

		openaiClient.ClearHistory(sender.ID)
		openaiClient.ClearRole(sender.ID)

		return c.Reply("Session cleared")
	})

	// Handler to clear chat history
	b.Handle("/clearhistory", func(c tele.Context) error {
		sender := c.Sender()

		openaiClient.ClearHistory(sender.ID)

		return c.Reply("History cleared")
	})

	// Handler to clear chat role
	b.Handle("/clearrole", func(c tele.Context) error {
		sender := c.Sender()

		openaiClient.ClearRole(sender.ID)

		return c.Reply("Role cleared")
	})

	// Handler to display help
	b.Handle("/help", func(c tele.Context) error {
		return c.Send("Available commands:\n" +
			"/history\n" +
			"/role\n" +
			"/setrole\n" +
			"/clear\n" +
			"/clearhistory\n" +
			"/clearrole\n" +
			"/help",
		)
	})

	// Start the bot
	fmt.Println("Telegram bot is running")
	b.Start()

	return nil
}
