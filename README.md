# TipsySync
This Application combines Tips/Donations, Follows & Subscriptions of most popular streaming-platforms and displays them in a nice overview.
It also collects statistics (soon), and offers a websocket connection to let other local applications use the data (eg. UE, Unity, Warudo, Lovense)

## Features
Currently implemented features:
- Event list (latest follows, subs, tips)
- Websocket sending events as they come in

## Planned Features
- More Providers (see Supported Platforms below)
- Infobuttons (on every page, display an Iconbutton at the top right, which expands a menu from the top explaining details of the current page)
- Chatlog page (combines all chatlogs from all sources)
- Overview (customizable widgets, resizable, drag+drop - basically components of pages in smaller version)
- Lovense & PiShock integration (maybe other toys aswell?)
- Statistics page (gets income by platform/date/day etc)
- Stream Overlays (since we already got relevant tip/sub data, why not create some overlays for OBS aswell?)

### Technical todo:
- improve logger:
    - log to file
    - add toast to frontend, when logger.Toast in golang, display toast in UI
    - add log-page in UI, displaying logger.Info and above
- Add Websocket settings (port, eventtypes to send)

## Supported platforms
I'm trying to add as many features as possible for all platforms, but some (eg. Fansly) do not offer api's and rely on webscraping or hijacking chat websockets.
Platforms ✅ supported or 🛠️ planned:

| Implemented | Provider   | Tips | Un-/Follow | Subscriptions | Chat |
|-------------|------------|------|------------|---------------|------|
| ✅          | Chaturbate | ✅   | ✅         | ✅            | ✅   |
| 🛠️💢        | Stripchat  | ✅   | ❌         | ❌            | ✅   |
| ✅💢        | Fansly     | ✅   | ❌         | ✅            | ✅   |
| ❌          | Onlyfans   |      |            |               |      |
| ❌          | YouTube    |      |            |               |      |
| ❌          | Twitch     |      |            |               |      |
| ❌          | Streamlabs |      |            |               |      |
| ❌          | Kick       |      |            |               |      |
| ❌          | Joystick   |      |            |               |      |

**Additional Info:**
* Platforms marked with 💢 do not offer official API's. This basically means, they could break at any time and require manual fixing.
* Fansly un-/follow should be possible, likely requires polling of some API endpoint or the followers page. i might implement it in the future, rn i got other priorities tho.
* Will consider checking if un-/follow & sub is possible on stripchat, got other priorities rn tho
* About joystick, this seems to be a new upcoming platform slowly getting traction. I might look into it after other platforms are finished.
* Platforms to consider in the future are CamSoda, BongaCams, LiveJasmin, ManyVids. Might start working on them if there's interest.

# Setting up
(TODO) Go to the [Releases Page](https://github.com/Nyoob/TipsySync/releases) and download the latest release.

#### Chaturbate
Go to the [Chaturbate Token Authorization Page](https://chaturbate.com/statsapi/authtoken/) and create an Events API Token
Paste the full Token URL into the Chaturbate Settings.

#### Fansly
Put in your Fansly Username into the Settings, without the leading @.

#### Stripchat
Put in your Stripchat Username into the Settings, without the leading @.

# Building it yourself
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

