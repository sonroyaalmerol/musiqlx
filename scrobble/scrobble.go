package scrobble

import (
	"time"

	"github.com/sonroyaalmerol/musiqlx/db"
)

type Track struct {
	Track                string
	Artist               string
	Album                string
	AlbumArtist          string
	TrackNumber          uint
	Duration             time.Duration
	MusicBrainzID        string
	MusicBrainzReleaseID string
}

type Scrobbler interface {
	IsUserAuthenticated(user db.User) bool
	Scrobble(user db.User, track Track, stamp time.Time, submission bool) error
}
