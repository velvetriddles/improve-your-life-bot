# Improve Your Life Bot

Improve Your Life Bot is a Telegram bot that helps users track useful activities and rewards. Users can perform various activities to earn coins and then spend these coins on rewards.

## Features

- Track useful activities to earn coins.
- Spend coins on rewards.
- View current balance.
- Interactive menus for ease of use.

## Getting Started

### Prerequisites

- Go (version 1.15 or higher)
- A Telegram bot token. You can obtain one by talking to [BotFather](https://core.telegram.org/bots#6-botfather).

### Installation

1. Fork && Clone the repository:

    ```sh
    git clone https://github.com/yourusername/improve-your-life-bot.git
    cd improve-your-life-bot
    ```

2. Set up environment variables:

    Create a `.env` file in the root directory of the project and add your Telegram bot token:

    ```env
    TOKEN_NAME_IN_OS=your_telegram_bot_token
    ```

3. Build the project:

    ```sh
    go build
    ```

4. Run the project:

    ```sh
    ./improve-your-life-bot
    ```

### Usage

Once the bot is running, you can interact with it via Telegram. Here are some basic commands:

- `/start`: Start interacting with the bot and see the introductory message.
- `View Introduction`: View the introduction and rules of the bot.
- `Skip Introduction`: Skip the introduction and go to the main menu.
- `Current Balance`: View your current coin balance.
- `Useful Activities`: View a list of activities you can perform to earn coins.
- `Rewards`: View a list of rewards you can purchase with your coins.
- `Main Menu`: Return to the main menu.

### Project Structure

- `main.go`: The main entry point of the bot.
- `constants.go`: Contains constants and configurations used throughout the project.
- `go.mod`: The Go module file, which defines the module path and its dependencies.
- `go.sum`: The Go checksum file, which contains the expected cryptographic checksums of the content of specific module versions.
- `README.md`: This file.

### Contributing

If you would like to contribute to this project, please open an issue or submit a pull request. We welcome all contributions!

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Acknowledgements

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Go Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api)

---

Made with ❤️ by velvetriddles
