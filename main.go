package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Config holds the bot configuration loaded from environment variables
type Config struct {
	OpenAIAPIKey string
	BotToken     string
	BotName      string
}

// loadConfig reads configuration from environment variables
func loadConfig() (*Config, error) {
	openAIKey := os.Getenv("WARP_DEMO_OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("WARP_DEMO_OPENAI_API_KEY environment variable is not set")
	}

	botToken := os.Getenv("WARP_DEMO_BOT_KEY")
	if botToken == "" {
		return nil, fmt.Errorf("WARP_DEMO_BOT_KEY environment variable is not set")
	}

	botName := os.Getenv("WARP_DEMO_BOT_NAME")
	if botName == "" {
		botName = "English Improver Bot" // default name
	}

	return &Config{
		OpenAIAPIKey: openAIKey,
		BotToken:     botToken,
		BotName:      botName,
	}, nil
}

// improveText uses OpenAI to improve the English text and identify issues
func improveText(ctx context.Context, client *openai.Client, text string) (string, error) {
	// Create a prompt that asks OpenAI to improve the text and list issues
	prompt := fmt.Sprintf(`You are an English language assistant. Improve the following text to B2 level, semiformal American English.

Original text:
%s

Please provide:
1. The improved version of the text
2. A categorized list of issues found:
   - Grammar issues
   - Spelling issues
   - Punctuation issues
   - Style issues
   - Semantic issues

Format your response as:
IMPROVED TEXT:
[improved text here]

ISSUES FOUND:
Grammar:
- [issue 1]
- [issue 2]

Spelling:
- [issue 1]

Punctuation:
- [issue 1]

Style:
- [issue 1]

Semantics:
- [issue 1]

If a category has no issues, write "None"`, text)

	// Call OpenAI API
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	response := chatCompletion.Choices[0].Message.Content
	return response, nil
}

// formatResponseHTML formats the OpenAI response as HTML for Telegram
func formatResponseHTML(response string) string {
	var result strings.Builder

	// Split response into improved text and issues sections
	parts := strings.Split(response, "ISSUES FOUND:")
	
	if len(parts) >= 1 {
		improvedSection := strings.TrimSpace(strings.TrimPrefix(parts[0], "IMPROVED TEXT:"))
		result.WriteString("<b>✅ Improved Text:</b>\n")
		result.WriteString("<i>" + improvedSection + "</i>\n\n")
	}

	if len(parts) >= 2 {
		result.WriteString("<b>📋 Issues Found:</b>\n")
		issuesSection := strings.TrimSpace(parts[1])
		
		// Format issues section with proper HTML
		lines := strings.Split(issuesSection, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			
			// Check if it's a category header (e.g., "Grammar:", "Spelling:")
			if strings.HasSuffix(trimmed, ":") && !strings.HasPrefix(trimmed, "-") {
				result.WriteString("\n<b>" + trimmed + "</b>\n")
			} else {
				result.WriteString(trimmed + "\n")
			}
		}
	}

	return result.String()
}

// handleMessage processes incoming text messages
func handleMessage(cfg *Config, openAIClient *openai.Client) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		// Only process text messages
		if ctx.EffectiveMessage.Text == "" {
			return nil
		}

		// Get the message text
		text := ctx.EffectiveMessage.Text

		// Ignore commands
		if strings.HasPrefix(text, "/") {
			return nil
		}

		log.Printf("Received message: %s", text)

		// Send typing indicator
		_, err := b.SendChatAction(ctx.EffectiveChat.Id, "typing", nil)
		if err != nil {
			log.Printf("Failed to send typing action: %v", err)
		}

		// Call OpenAI to improve the text
		improved, err := improveText(context.Background(), openAIClient, text)
		if err != nil {
			log.Printf("Error improving text: %v", err)
			_, sendErr := ctx.EffectiveMessage.Reply(b, "❌ Sorry, I encountered an error processing your message.", nil)
			if sendErr != nil {
				log.Printf("Failed to send error message: %v", sendErr)
			}
			return err
		}

		// Format the response as HTML
		formattedResponse := formatResponseHTML(improved)

		// Send the response
		_, err = ctx.EffectiveMessage.Reply(b, formattedResponse, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		if err != nil {
			log.Printf("Failed to send response: %v", err)
			return err
		}

		log.Printf("Successfully processed and replied to message")
		return nil
	}
}

// handleStart handles the /start command
func handleStart(cfg *Config) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		welcomeMsg := fmt.Sprintf(`👋 Welcome to <b>%s</b>!

Send me any English text, and I'll:
✅ Improve it to B2 level semiformal American English
📋 List all grammar, spelling, punctuation, style, and semantic issues

Just send your text and I'll help you improve it!`, cfg.BotName)

		_, err := ctx.EffectiveMessage.Reply(b, welcomeMsg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		if err != nil {
			log.Printf("Failed to send welcome message: %v", err)
			return err
		}
		return nil
	}
}

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s...", cfg.BotName)

	// Create OpenAI client
	openAIClient := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	// Create Telegram bot
	bot, err := gotgbot.NewBot(cfg.BotToken, nil)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Create updater
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Printf("Error occurred: %v", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	// Register handlers
	dispatcher.AddHandler(handlers.NewCommand("start", handleStart(cfg)))
	dispatcher.AddHandler(handlers.NewMessage(nil, handleMessage(cfg, &openAIClient)))

	// Start polling
	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		log.Fatalf("Failed to start polling: %v", err)
	}

	log.Printf("%s is now running. Press Ctrl+C to stop.", cfg.BotName)

	// Block until stopped
	updater.Idle()
}
