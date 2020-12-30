package player

import (
	"sync"
	"time"
)

type Player struct {
	Progress PlayerProgress
	Queue PlayerQueue

	mutex sync.RWMutex
}

type PlayerProgress struct {
	Progress time.Duration
	EpisodeGUID string
	PodcastID int
}

type PlayerQueue struct {

}

//func NewSynchronizedPlayer(db *Database) *Player {
//	p := Player{}
//	// TODO
//	return &p
//}
