package youtube

import (
	"context"
	"fmt"

	"github.com/lrstanley/go-ytdlp"
)

type YouTubeStreamer struct {
	VideoId string
}

func InitializeStreamer(id string) *YouTubeStreamer {
	return &YouTubeStreamer{
		VideoId: id,
	}
}

func (streamer *YouTubeStreamer) Stream() error {
	dl := ytdlp.New().
		SponsorblockRemove("music_offtopic").
		ExtractAudio().
		AudioFormat("best").
		AudioQuality("0").
		Format("ba[ext=m4a]").
		Output(fmt.Sprintf("%s-cache.%%(ext)s", streamer.VideoId))

	_, err := dl.Run(context.TODO(), fmt.Sprintf("https://music.youtube.com/watch?v=%s", streamer.VideoId))
	if err != nil {
		return err
	}

	return nil
}
