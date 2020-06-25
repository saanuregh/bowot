package music

import (
	"bowot/internal/embeds"
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/lrita/cmap"
	dca "github.com/yyewolf/dca-disgord"
)

const (
	PLAYER_STOP  = 0
	PLAYER_PLAY  = 1
	PLAYER_PAUSE = 2
)

type GuildPlayer struct {
	ID              string
	Queue           []Track
	VoiceConnection disgord.VoiceConnection
	VoiceStates     []*disgord.VoiceState
	PlayerStatus    PlayerStatus
	CurrentTrack    Track
	StatusChan      chan int
	StreamSession   *dca.StreamingSession
	streamerMutex   *sync.Mutex
}

type Track struct {
	Title       string
	URL         string
	DownloadURL string
	Source      string
	IsLive      bool
	Duration    time.Duration
	Requester   disgord.Snowflake
	Channel     disgord.Snowflake
}

type PlayerStatus struct {
	Current  int
	Previous int
}

var GuildPlayers cmap.Cmap

func addGuild(guildID string) *GuildPlayer {
	g := &GuildPlayer{
		ID:              guildID,
		CurrentTrack:    Track{},
		Queue:           []Track{},
		VoiceStates:     make([]*disgord.VoiceState, 0),
		VoiceConnection: nil,
		PlayerStatus:    PlayerStatus{Current: PLAYER_STOP, Previous: PLAYER_STOP},
		StatusChan:      make(chan int, 10),
		StreamSession:   nil,
		streamerMutex:   &sync.Mutex{},
	}
	GuildPlayers.Store(guildID, g)
	go g.PlayerWorker()
	return g
}

func removeGuild(guildID string) *GuildPlayer {
	if tmp, ok := GuildPlayers.Load(guildID); !ok {
		guild := tmp.(*GuildPlayer)
		GuildPlayers.Delete(guildID)
		return guild
	}
	return nil
}

func (G *GuildPlayer) PlayerWorker() {
	options := dca.StdEncodeOptions
	options.CompressionLevel = 5
	for {
		select {
		case status := <-G.StatusChan:
			G.PlayerStatus.Previous = G.PlayerStatus.Current
			G.PlayerStatus.Current = status
			if G.PlayerStatus.Current == PLAYER_PLAY {
				if G.PlayerStatus.Previous == PLAYER_PAUSE {
					G.StreamSession.SetPaused(false)
				} else {
					go func() {
						G.streamerMutex.Lock()
						defer G.streamerMutex.Unlock()
						G.CurrentTrack = G.Queue[0]
						G.RemoveTrack(0)
						fields := []*disgord.EmbedField{
							embeds.Field("Title", G.CurrentTrack.Title, false),
							embeds.Field("Requester", fmt.Sprintf("<@%s>", G.CurrentTrack.Requester), true),
							embeds.Field("Source", G.CurrentTrack.Source, true),
							nil,
							embeds.Field("URL", G.CurrentTrack.URL, false),
						}
						if G.CurrentTrack.IsLive {
							fields[3] = embeds.Field("Duration", "Live", true)
						} else {
							fields[3] = embeds.Field("Duration", G.CurrentTrack.Duration.String(), true)
						}
						client.SendMsg(
							context.Background(),
							G.CurrentTrack.Channel,
							embeds.Info(
								"Music - Now playing",
								"",
								"",
								fields...,
							),
						)
						encodingSession, err := dca.EncodeFile(G.CurrentTrack.DownloadURL, options)
						if err != nil {
							client.Logger().Error(err)
							return
						}
						defer encodingSession.Cleanup()
						err = G.VoiceConnection.StartSpeaking()
						if err != nil {
							client.Logger().Error(err)
						}
						done := make(chan error)
						G.StreamSession = dca.NewStream(encodingSession, G.VoiceConnection, done)
						G.StreamSession.SetPaused(false)
						err = <-done
						if err != nil && err != io.EOF {
							client.Logger().Error(err)
						}
						if G.PlayerStatus.Current == PLAYER_PLAY {
							if len(G.Queue) > 0 {
								G.Play()
							} else {
								client.SendMsg(
									context.Background(),
									G.CurrentTrack.Channel,
									"Stopping playback.",
								)
								G.Stop()
							}
						}
					}()
				}
			}
			if G.PlayerStatus.Current == PLAYER_PAUSE {
				if G.StreamSession != nil {
					if !G.StreamSession.Paused() {
						G.StreamSession.SetPaused(true)
					}
				}
			}
			if G.PlayerStatus.Current == PLAYER_STOP {
				if G.StreamSession != nil {
					if G.PlayerStatus.Previous == PLAYER_PLAY {
						G.StreamSession.SetPaused(true)
						finished, err := G.StreamSession.Finished()
						if !finished && err == nil {
							G.StreamSession.Stop()
						}
					}
					G.StreamSession = nil
					G.CurrentTrack = Track{}
					G.VoiceConnection.StopSpeaking()
				}
			}
		}
	}
}

func (G *GuildPlayer) AddVoiceState(vs *disgord.VoiceState) {
	G.VoiceStates = append(G.VoiceStates, vs)
}

func (G *GuildPlayer) GetVoiceState(userID disgord.Snowflake) (*disgord.VoiceState, int) {
	for i, v := range G.VoiceStates {
		if v.UserID == userID {
			return v, i
		}
	}
	return nil, -1
}

func (G *GuildPlayer) RemoveVoiceState(userID disgord.Snowflake) {
	_, i := G.GetVoiceState(userID)
	if i > -1 {
		G.VoiceStates = append(G.VoiceStates[:i], G.VoiceStates[i+1:]...)
	}
}

func (G *GuildPlayer) UpdateVoiceState(userID disgord.Snowflake, vs *disgord.VoiceState) {
	_, i := G.GetVoiceState(userID)
	if i > -1 {
		G.VoiceStates[i] = vs
	}
}

func (G *GuildPlayer) UpdateVoiceConnection(vc disgord.VoiceConnection) {
	G.VoiceConnection = vc
}

func (G *GuildPlayer) AddTrack(title, source, URL, downloadURL string, isLive bool, duration int64, requesterID, channelID disgord.Snowflake) {
	G.Queue = append(G.Queue, Track{
		Title:       title,
		URL:         URL,
		DownloadURL: downloadURL,
		Source:      source,
		Requester:   requesterID,
		Channel:     channelID,
		IsLive:      isLive,
		Duration:    time.Duration(duration) * time.Second,
	})
}

func (G *GuildPlayer) RemoveTrack(i int) {
	G.Queue = append(G.Queue[:i], G.Queue[i+1:]...)
}

func (G *GuildPlayer) RemoveAllTracks() {
	G.Queue = make([]Track, 0)
}

func (G *GuildPlayer) ShuffleTracks() {
	rand.Shuffle(len(G.Queue), func(i, j int) {
		G.Queue[i], G.Queue[j] = G.Queue[j], G.Queue[i]
	})
}

func (G *GuildPlayer) Play() {
	G.StatusChan <- PLAYER_PLAY
}

func (G *GuildPlayer) Pause() {
	G.StatusChan <- PLAYER_PAUSE
}

func (G *GuildPlayer) Stop() {
	G.StatusChan <- PLAYER_STOP
}

func (G *GuildPlayer) Skip() {
	G.StatusChan <- PLAYER_STOP
	G.streamerMutex.Lock()
	G.streamerMutex.Unlock()
	G.StatusChan <- PLAYER_PLAY
}
