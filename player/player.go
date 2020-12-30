package player

import (
	"sync"
	"time"

	"lincast/database"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	currentProgress *CurrentProgress
	queue           *Queue
	db              *database.Database
	mutex           sync.RWMutex
}

type CurrentProgress struct {
	Progress    time.Duration
	EpisodeGUID string
	PodcastID   int
}

type Queue struct {
}

// New returns a new Player synchronized with the given database.
func New(db *database.Database) (*Player, error) {
	p := Player{
		currentProgress: new(CurrentProgress),
		queue:           new(Queue),
		db:              db,
	}

	err := p.initProgress()
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "error when trying to initialize the progress"+
			" on the database")
	}

	return &p, nil
}

// UpdateProgress updates the progress in the database and caches it internally.
func (p *Player) UpdateProgress(newProgress time.Duration, episodeGUID string, podcastID int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.currentProgress.Progress = newProgress
	p.currentProgress.EpisodeGUID = episodeGUID
	p.currentProgress.PodcastID = podcastID

	return p.updateProgressOnDB()
}

func (p *Player) GetProgress() CurrentProgress {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return *p.currentProgress
}

func (p *Player) updateProgressOnDB() error {
	iDB := p.db.GetInstance()
	query := `
UPDATE player_progress
SET progress = ?, episode_guid = ?, podcast_id = ?, user = ?
WHERE id = 0;
`

	r, err := iDB.Exec(query, p.currentProgress.Progress, p.currentProgress.EpisodeGUID, p.currentProgress.PodcastID, "")
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

func (p *Player) initProgress() error {
	initialized, err := p.isDBProgressInitialized()
	if err != nil {
		return err
	}

	iDB := p.db.GetInstance()

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
			&p.currentProgress.Progress,
			&p.currentProgress.EpisodeGUID,
			&p.currentProgress.PodcastID,
			&user,
		)
		if err != nil {
			return err
		}

		return nil
	}

	// If there isn't a row, we need to create one with empty values.
	query := "INSERT INTO player_progress (id, progress, episode_guid, podcast_id, user) VALUES (0, ?, ?, ?, ?)"

	r, err := iDB.Exec(query, p.currentProgress.Progress, p.currentProgress.EpisodeGUID, p.currentProgress.PodcastID, "")
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

func (p *Player) isDBProgressInitialized() (bool, error) {
	iDB := p.db.GetInstance()
	query := `
SELECT id FROM episodes
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
