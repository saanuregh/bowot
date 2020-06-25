# BOWOT
A fun discord bot made with Golang using [disgord](https://github.com/andersfylling/disgord), [gommand](https://github.com/auttaja/gommand) and [FaunaDB](https://fauna.com/).

## Features
* Multi-guild support with guild specific settings.
  * Guild specific prefix.
  * Unlimited custom commands per guild.
  * Wakephrase support which lets bot to listen to messages for certain phrases and react to them (ratelimited to avoid spam). 
* Low memory footprint.
* Selfroles
  * Autodetects guild selfroles using regex set per guild by the administrator.
  * Supports default selfrole for new members.
  * User settings to pick your own selfrole.
* Easy and free database using [FaunaDB](https://fauna.com/).
* Music
  * Using [youtube-dl](https://github.com/ytdl-org/youtube-dl) and [ffmpeg](https://ffmpeg.org/).
  * Supports every site supported by [youtube-dl](https://github.com/ytdl-org/youtube-dl) such as Twitch livestreams, soundcloud, etc...
  * Supports spotify (converted to youtube links) and youtube playlist. 
  * Queues
* Reddit commands to get random reddit posts.
  * Subreddits configurable using config file.
* Economy - get some cOwOins.
  * Gamble, daily, leaderboards
* Hydrate reminder
  * Reminds added user by DM to drink some water every 45 mins.
* Intelligent bot reply using [Botlibre](https://www.botlibre.com) (This chatbot AI will be replaced by a more powerful chatbot AI in future). Users can talk to bowot by mentioning him (@).
* Error reporting using sentry.
* Configuration using YAML.
* Docker deployment.
  * Can be deployed easily to heroku for free using their [Container Registry CLI](https://devcenter.heroku.com/articles/container-registry-and-runtime)

## Deployment
1. Clone the repository
2. Rename the configuration file `config.yaml.example` to `config.yaml` and edit it.
3. Build the docker and run it (the config file is placed in /root/ if you want to provide configuration file dynamically).

## Commands available

| Command     |                  Description                  |
| ----------- | :-------------------------------------------: |
| about       |                 About bowot.                  |
| balance     |          Check your balance cOwOins.          |
| chuck       |        Gets random Chuck Norris joke.         |
| customs     |      Get all the guild custom commands.       |
| daily       |           Grab your daily cowoins.            |
| define      |                Defines a word.                |
| dice        |                 Throw a dice.                 |
| gamble      |                Gamble cowoins.                |
| leaderboard |              cOwOin leaderboard.              |
| ping        |        Get average latency of the bot.        |
| poll        |          Make a reaction based poll.          |
| meme        |         Gets random meme from reddit.         |
| whoosh      |        Gets random whoosh from reddit.        |
| copypasta   |      Gets random copypasta from reddit.       |
| selfroles   |         Get all the guild selfroles.          |
| wakephrases |        Get all the guild wakephrases.         |
| translate   |    Translate from one language to another.    |
| user        |            User specific settings.            |
| uwufy       |               Uwufy a sentence.               |
| valorant    |      Get current valorant server status.      |
| settings    | Guild specific settings (administrator only). |
| baka        |              Call somebody baka.              |
| cuddle      |               Cuddle somebody.                |
| hug         |                 Hug somebody.                 |
| kiss        |                Kiss somebody.                 |
| pat         |                 Pat somebody.                 |
| poke        |                Poke somebody.                 |
| slap        |                Slap somebody.                 |
| smug        |                Smug somebody.                 |
| tickle      |               Tickle somebody.                |
| join        |        Join the current voice channel.        |
| play        |            Play a song or unpause.            |
| pause       |                Pause playback.                |
| stop        |                Stop playback.                 |
| skip        |              Skip current song.               |
| shuffle     |              Shuffle the queue.               |
| clear       |               Clear the queue.                |
| remove      |        Remove a track from the queue.         |
| leave       |       Leave the current voice channel.        |
| queue       |               Shows the queue."               |