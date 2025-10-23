## Goals
Build a Telegram bot that is improving English message (B2, semiformal American) and send as response improved message with lists of issues (grammar, spelling, punctuation, style, semantics).

## Constraints
- Use Go 1.21+.
- Use https://github.com/PaulSonOfLars/gotgbot for working with Telegram API
- Use https://github.com/openai/openai-go for working with OpenAI API 
- No external network calls beyond Telegram/OpenAI in core flow.

## Agent Behavior
- Propose diffs, do not push to remote without confirmation.
- Keep code clear and commented.
- Prefer HTML formatting in Telegram replies to avoid Markdown escaping issues.


## Enviroment Variables
WARP_DEMO_OPENAI_API_KEY
WARP_DEMO_BOT_KEY
WARP_DEMO_BOT_NAME
