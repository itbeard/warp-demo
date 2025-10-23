# English Improver Telegram Bot

A Telegram bot that improves English messages to B2 level semiformal American English and provides detailed feedback on grammar, spelling, punctuation, style, and semantic issues.

## Features

- ✅ Improves English text to B2 level semiformal American English
- 📋 Provides categorized lists of issues:
  - Grammar
  - Spelling
  - Punctuation
  - Style
  - Semantics
- 🤖 Powered by OpenAI GPT-4
- 💬 HTML-formatted responses for better readability

## Prerequisites

- Go 1.21 or higher
- Telegram Bot Token (get it from [@BotFather](https://t.me/botfather))
- OpenAI API Key (get it from [OpenAI Platform](https://platform.openai.com/))

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd warp-demo
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export WARP_DEMO_OPENAI_API_KEY="your-openai-api-key"
export WARP_DEMO_BOT_KEY="your-telegram-bot-token"
export WARP_DEMO_BOT_NAME="English Improver Bot"  # optional, has default value
```

## Running the Bot

```bash
go run main.go
```

The bot will start polling for messages. Press Ctrl+C to stop.

## Usage

1. Start a conversation with your bot on Telegram
2. Send `/start` to see the welcome message
3. Send any English text message
4. The bot will respond with:
   - The improved version of your text
   - A categorized list of all issues found

## Example

**You send:**
```
i go to store yesterday and buy some apples
```

**Bot responds:**
```
✅ Improved Text:
I went to the store yesterday and bought some apples.

📋 Issues Found:

Grammar:
- "go" should be "went" (past tense)
- "buy" should be "bought" (past tense)

Punctuation:
- Missing period at the end of the sentence
- First word should be capitalized

Spelling:
None

Style:
- Added article "the" before "store" for clarity

Semantics:
None
```

## Project Structure

```
warp-demo/
├── main.go          # Main bot implementation
├── go.mod           # Go module dependencies
├── go.sum           # Dependency checksums
├── warp.md          # Project specifications
└── README.md        # This file
```

## Configuration

The bot reads configuration from environment variables:

- `WARP_DEMO_OPENAI_API_KEY` (required): Your OpenAI API key
- `WARP_DEMO_BOT_KEY` (required): Your Telegram bot token
- `WARP_DEMO_BOT_NAME` (optional): Custom name for the bot (default: "English Improver Bot")

## Technologies Used

- **Go 1.21+**: Programming language
- **gotgbot v2**: Telegram Bot API library ([github.com/PaulSonOfLars/gotgbot](https://github.com/PaulSonOfLars/gotgbot))
- **openai-go**: OpenAI API client ([github.com/openai/openai-go](https://github.com/openai/openai-go))
- **OpenAI GPT-4o**: Language model for text improvement

## Development

The bot follows these principles:
- Clear, commented code
- HTML formatting for Telegram responses to avoid Markdown escaping issues
- No external network calls beyond Telegram and OpenAI
- Comprehensive error handling and logging

---

Дата создания проекта: 2024-10-23
