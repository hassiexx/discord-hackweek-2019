# discord-hackweek-2019
A moderation bot in Go for Discord Hack Week 2019 using menu command systems for improved usability and interaction.
Still heavily in progress though, no commands running as of yet.

# Instructions

## Setup your environment variables
 - ``DISCORD_HACK_WEEK_2019_TOKEN`` - This is your bot token
 - ``DISCORD_HACK_WEEK_2019_MONGODB`` - This is your MongoDB URI in the form ``mongodb://[username:password@]host1[:port1][,...hostN[:portN]]][/[database][?options]]``, see <https://docs.mongodb.com/manual/reference/connection-string/>. If you are Discord staff testing the bot, feel free to DM me to get access to my hosted MongoDB cloud instance.

## Install Go
 - Ensure you have Go installed, this was built with 1.12.4.

## Run the bot
 - Clone the repo; as this uses Go modules, it does not need to be placed in the Go path.
 - Open your command shell and navigate to the repo folder.
 - Fetch all required dependencies for the bot: ``go get ./...``
 - Run the bot: ``go run main/main.go``
