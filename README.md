# Tip Aggregator
This Application combines Tips/Donations, Follows & Subscriptions of most popular streaming-platforms and displays them in a nice overview.
It also collects statistics, and offers a websocket connection to let other local applications use the data (eg. UE, Unity, Warudo, Lovense)

## Supported platforms
I'm trying to add as many features as possible for all platforms, but some (eg. Fansly) do not offer api's and rely on webscraping.
Platforms supported or planned:
- Chaturbate
- Stripchat (planned)
- Fansly (planned)
- Onlyfans (planned)
- Youtube (planned)
- Twitch (planned)

## Features
Current implemented or planned features:
- Overview with latest events (half-done)
- Settings page (half-done)
- Statistics page (planned)
- Websocket sending events as they come in (planned)

# Setting up
(TODO) Go to the [Releases Page](https://github.com/Nyoob/tip-aggregator/releases) and download the latest release.

## Building it yourself
Clone the repo, run `go mod download`, head into /frontend, run `yarn install`.
Then run `wails build`. Check out the [Wails Documentation](https://wails.io/docs/reference/cli#build) for more info.

# Develop

## About
This project is written in golang, using Wails with React for the frontend.
Data is stored in an sqlite DB.

## Live Development
To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
