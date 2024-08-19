package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/chatgpt-telegram-bot/cmd/run"
)

var rootCmd = &cobra.Command{
	Use:               "chatgpt-telegram-bot",
	Short:             "ChatGPT Telegram bot",
	Long:              `Telegram chat bot powered by OpenAI's GPT-3.`,
	PersistentPreRunE: logLevel,
	SilenceUsage:      true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(run.Cmd)
}
