package psync

import (
	"sync"
	"time"

	"lincast/database"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

// Synchronizer has in charge the synchronization of the different capabilities of the player (current progress,
// queue, etc) across the different clients.
type Synchronizer struct {
	currentProgress *CurrentProgress
	queue           *Queue
	db              *database.Database
	mutex           sync.RWMutex
}

// CurrentProgress is the structure used to store and parse the information related with the episode that is being
// currently playing on the player.
type CurrentProgress struct {
	Progress    time.Duration `json:"progress"`
	EpisodeGUID string        `json:"episode_id"`
	PodcastID   int           `json:"podcast_id"`
}

// Queue is the structure that represents the queue of the player (located on the client), and is used for its storage,
// manipulation and synchronization across clients.
type Queue struct {
	Content []QueueEpisode
}

type QueueEpisode struct {
	ID        int    `json:"id"`
	PodcastID int    `json:"podcast_id"`
	EpisodeID string `json:"episode_id"`
	Position  int    `json:"position"`
}

// New returns a new Synchronizer.
func New(db *database.Database) (*Synchronizer, error) {
	if db == nil {
		return nil, errorx.IllegalState.New("the db used can't be nil")
	}

	s := Synchronizer{
		currentProgress: new(CurrentProgress),
		queue:           new(Queue),
		db:              db,
	}

	err := s.initProgress()
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "error when trying to initialize the progress"+
			" on the database")
	}

	err = s.initQueue()
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "error when trying to get the queue from the database")
	}

	return &s, nil
}

// UpdateProgress updates the progress of the player in the database and caches it internally.
func (s *Synchronizer) UpdateProgress(newProgress time.Duration, episodeGUID string, podcastID int) error {
	if episodeGUID == "" {
		return errorx.IllegalArgument.New("the episodeGUID should not be empty")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.currentProgress.Progress = newProgress
	s.currentProgress.EpisodeGUID = episodeGUID
	s.currentProgress.PodcastID = podcastID

	return s.updateProgressOnDB()
}

// GetProgress returns the current progress of the player.
func (s *Synchronizer) GetProgress() CurrentProgress {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return *s.currentProgress
}

func (s *Synchronizer) updateProgressOnDB() error {
	iDB := s.db.GetInstance()
	query := `
UPDATE player_progress
SET progress = ?, episode_guid = ?, podcast_id = ?, user = ?
WHERE id = 0;
`

	r, err := iDB.Exec(query, s.currentProgress.Progress, s.currentProgress.EpisodeGUID, s.currentProgress.PodcastID, "")
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalState.New("the row to store the progress of the player apparently" +
			" does not exist")
	}

	return nil
}

func (s *Synchronizer) initProgress() error {
	initialized, err := s.isDBProgressInitialized()
	if err != nil {
		return err
	}

	iDB := s.db.GetInstance()

	// If there is a row, we need to get those values to set the cache.
	if initialized {
		query := "SELECT * FROM player_progress WHERE id = 0"

		row, err := iDB.Query(query)
		if err != nil {
			return err
		}

		defer func() {
			err = row.Close()
			if err != nil {
				log.Error(errorx.Decorate(err, "error when trying to close rows"))
			}
		}()

		// We know that there is one row, so there is no need to check the returned value.
		_ = row.Next()

		var id int      // Not used
		var user string // Not used
		err = row.Scan(
			&id,
			&s.currentProgress.Progress,
			&s.currentProgress.EpisodeGUID,
			&s.currentProgress.PodcastID,
			&user,
		)
		if err != nil {
			return err
		}

		return nil
	}

	// If there isn't a row, we need to create one with empty values.
	query := "INSERT INTO player_progress (id, progress, episode_guid, podcast_id, user) VALUES (0, ?, ?, ?, ?)"

	r, err := iDB.Exec(query, s.currentProgress.Progress, s.currentProgress.EpisodeGUID, s.currentProgress.PodcastID, "")
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.AssertionFailed.New("for some reason the default row wasn't added")
	}

	return nil
}

func (s *Synchronizer) isDBProgressInitialized() (bool, error) {
	iDB := s.db.GetInstance()
	query := `
SELECT id FROM player_progress
WHERE id = 0;
`
	// Check if the row exists.
	row, err := iDB.Query(query)
	if err != nil {
		return false, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	return row.Next(), nil
}

func (s *Synchronizer) initQueue() error {
	query := "SELECT * FROM player_queue"

	db := s.db.GetInstance()

	row, err := db.Query(query)
	if err != nil {
		return err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	var q Queue

	for row.Next() {
		var e QueueEpisode
		err = row.Scan(&e.ID, &e.PodcastID, &e.EpisodeID)
		if err != nil {
			return err
		}

		q.Content = append(q.Content, e)
	}

	s.queue = &q

	return nil
}
