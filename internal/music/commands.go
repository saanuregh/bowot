package music

import (
	"bowot/internal/embeds"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

var cmds []*gommand.Command

func init() {
	musicCategory := &gommand.Category{
		Name:        "Music",
		Description: "Play some tunes.",
	}
	cmds = append(
		cmds,
		&gommand.Command{
			Name:        "join",
			Description: "Join the current voice channel.",
			Category:    musicCategory,
			Function:    join,
		},
		&gommand.Command{
			Name:        "play",
			Description: "Play a song or unpause.",
			Category:    musicCategory,
			Function:    play,
		},
		&gommand.Command{
			Name:        "pause",
			Description: "Pause playback.",
			Category:    musicCategory,
			Function:    pause,
		},
		&gommand.Command{
			Name:        "stop",
			Description: "Stop playback.",
			Category:    musicCategory,
			Function:    stop,
		},
		&gommand.Command{
			Name:        "skip",
			Description: "Skip current song.",
			Category:    musicCategory,
			Function:    skip,
		},
		&gommand.Command{
			Name:        "shuffle",
			Description: "Shuffle the queue.",
			Category:    musicCategory,
			Function:    shuffle,
		},
		&gommand.Command{
			Name:        "clear",
			Description: "Clear the queue.",
			Category:    musicCategory,
			Function:    clear,
		},
		&gommand.Command{
			Name:        "remove",
			Description: "Remove a track from the queue.",
			Category:    musicCategory,
			Function:    remove,
			ArgTransformers: []gommand.ArgTransformer{
				{Function: gommand.IntTransformer},
			},
		},
		&gommand.Command{
			Name:        "leave",
			Description: "Leave the current voice channel.",
			Category:    musicCategory,
			Function:    leave,
		},
		&gommand.Command{
			Name:        "queue",
			Description: "Shows the queue.",
			Category:    musicCategory,
			Function:    queue,
		},
		&gommand.Command{
			Name:        "now",
			Description: "Shows now playing.",
			Category:    musicCategory,
			Function:    now,
		},
	)
}

func _join(guild *GuildPlayer, ctx *gommand.Context) error {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			vc, err := ctx.Session.VoiceConnect(ctx.Message.GuildID, vs.ChannelID)
			if err != nil {
				return err
			}
			guild.UpdateVoiceConnection(vc)
			ch, err := ctx.Session.GetChannel(context.Background(), vs.ChannelID)
			if err != nil {
				return err
			}
			_, err = ctx.Reply(fmt.Sprintf("Joined channel - %s.", ch.Name))
			return err
		}
	}
	return fmt.Errorf("voice channel")
}

func join(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		err := _join(guild, ctx)
		if err != nil && err.Error() != "voice channel" {
			return err
		}
		return nil
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func play(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		s := guild.PlayerStatus.Current
		if ctx.RawArgs == "" {
			if s == PLAYER_PAUSE {
				guild.Play()
				_, err := ctx.Reply("Playing.")
				return err
			} else if s == PLAYER_STOP {
				if len(guild.Queue) > 0 {
					guild.Play()
					_, err := ctx.Reply("Playing.")
					return err
				} else {
					_, err := ctx.Reply("Nothing playing.")
					return err
				}
			} else {
				_, err := ctx.Reply("Already playing something.")
				return err
			}
		} else {
			if guild.VoiceConnection == nil {
				err := _join(guild, ctx)
				if err != nil {
					if err.Error() == "voice channel" {
						_, err := ctx.Reply("Join a voice channel first.")
						return err
					}
					return err
				}
			}
			msg, err := ctx.Reply("Loading...")
			videos, err := GetYoutubeDLResponse(ctx.RawArgs)
			if err != nil {
				return err
			}
			updatedMsg := ctx.Session.UpdateMessage(
				context.Background(),
				msg.ChannelID,
				msg.ID,
			)
			if len(videos) < 1 {
				_, err := updatedMsg.SetContent("Couldn't find the video.").Execute()
				return err
			}
			for _, v := range videos {
				var duration int64
				if v.Duration != nil {
					duration = int64(*v.Duration)
				} else {
					duration = 0
				}
				guild.AddTrack(v.Title, v.ExtractorKey, v.Query, v.URL, v.IsLive, duration, ctx.Message.Author.ID, ctx.Message.ChannelID)
			}
			_, err = updatedMsg.SetContent(fmt.Sprintf("Queing - %v songs.", len(videos))).Execute()
			if s == PLAYER_STOP {
				guild.Play()
			}
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func pause(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if guild.PlayerStatus.Current == PLAYER_PLAY {
			guild.Pause()
			_, err := ctx.Reply("Pausing.")
			return err
		} else if guild.PlayerStatus.Current == PLAYER_PAUSE {
			_, err := ctx.Reply("Already paused.")
			return err
		} else {
			_, err := ctx.Reply("Nothing playing.")
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func stop(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if guild.PlayerStatus.Current == PLAYER_STOP {
			_, err := ctx.Reply("Nothing playing.")
			return err
		} else {
			guild.Stop()
			_, err := ctx.Reply("Stopping.")
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func skip(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if guild.PlayerStatus.Current == PLAYER_STOP {
			_, err := ctx.Reply("Nothing playing.")
			return err
		} else if len(guild.Queue) < 1 {
			guild.Stop()
			_, err := ctx.Reply("Queue empty.")
			return err
		} else {
			guild.Skip()
			_, err := ctx.Reply("Skipping current song.")
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func shuffle(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if len(guild.Queue) > 0 {
			guild.ShuffleTracks()
			_, err := ctx.Reply("Shuffling.")
			return err
		}
		_, err := ctx.Reply("Queue empty.")
		return err
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func clear(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if len(guild.Queue) > 0 {
			guild.RemoveAllTracks()
			_, err := ctx.Reply("Clearing queue.")
			return err
		}
		_, err := ctx.Reply("Queue empty.")
		return err
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func remove(ctx *gommand.Context) error {
	index := ctx.Args[0].(int)
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if len(guild.Queue) > 0 {
			if index > -1 && index < len(guild.Queue) {
				guild.RemoveTrack(index)
				_, err := ctx.Reply("Cant find the guild.")
				return err
			} else {
				_, err := ctx.Reply("Index out of range.")
				return err
			}
		}
		_, err := ctx.Reply("Queue empty.")
		return err
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func leave(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		for _, vs := range guild.VoiceStates {
			if vs.UserID == ctx.Message.Author.ID {
				if guild.VoiceConnection != nil {
					guild.Stop()
					err := guild.VoiceConnection.Close()
					guild.UpdateVoiceConnection(nil)
					if err != nil {
						return err
					}
					ch, err := ctx.Session.GetChannel(context.Background(), vs.ChannelID)
					if err != nil {
						return err
					}
					_, err = ctx.Reply(fmt.Sprintf("Left channel - %s.", ch.Name))
					return err
				} else {
					_, err := ctx.Reply("Not on any voice channel.")
					return err
				}
			}
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func queue(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if len(guild.Queue) < 1 {
			_, err := ctx.Reply("Queue empty.")
			return err
		} else {
			desc := "```"
			if guild.CurrentTrack.URL != "" {
				desc += fmt.Sprintf("Now playing - %v\n", guild.CurrentTrack.Title)
			}
			for i, t := range guild.Queue {
				desc += fmt.Sprintf("%v - %v\n", i+1, t.Title)
			}
			desc += "```"
			_, err := ctx.Reply(embeds.Info(
				"Music - Queue",
				desc,
				"",
			))
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}

func now(ctx *gommand.Context) error {
	if tmp, ok := GuildPlayers.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*GuildPlayer)
		if guild.PlayerStatus.Current == PLAYER_STOP {
			_, err := ctx.Reply("Nothing playing.")
			return err
		} else {
			fields := []*disgord.EmbedField{
				embeds.Field("Title", guild.CurrentTrack.Title, false),
				embeds.Field("Requester", fmt.Sprintf("<@%s>", guild.CurrentTrack.Requester), true),
				embeds.Field("Source", guild.CurrentTrack.Source, true),
				nil,
				embeds.Field("URL", guild.CurrentTrack.URL, false),
			}
			if guild.CurrentTrack.IsLive {
				fields[3] = embeds.Field("Duration", "Live", true)
			} else {
				fields[3] = embeds.Field("Duration", guild.CurrentTrack.Duration.String(), true)
			}
			_, err := ctx.Reply(embeds.Info(
				"Music - Now playing",
				"",
				"",
				fields...,
			))
			return err
		}
	}
	_, err := ctx.Reply("Join a voice channel first.")
	return err
}
