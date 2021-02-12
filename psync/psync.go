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
	len     int
}

// QueueEpisode represents an episode of the queue.
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

// GetQueue returns the actual queue of the player.
func (s *Synchronizer) GetQueue() Queue {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return *s.queue
}

// SetQueue overwrites the entire player's queue with the given content.
func (s *Synchronizer) SetQueue(eps *[]QueueEpisode) error {
	// Clean the queue to later add the new episodes.
	err := s.CleanQueue()
	if err != nil {
		return err
	}

	// Get the real instance of the database (*sql.DB) to execute queries directly on it.
	sqlDB := s.db.GetInstance()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store each episode on the table `player_queue`.
	for _, ep := range *eps {
		query := "INSERT INTO player_queue (podcast_id, episode_id, position) VALUES (?, ?, ?);"

		result, err := sqlDB.Exec(query, ep.PodcastID, ep.EpisodeID, ep.Position)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// If there are no rows affected, then, for some reason, the query has made no effect.
		if rowsAffected == 0 {
			return errorx.InternalError.New("no rows have been affected")
		}
	}

	// Get the newly stored queue.
	epsFromDB, err := s.getQueueEpsFromDB()
	if err != nil {
		return err
	}

	// Finally, update the queue in memory with the values that are in the database. This will give us the ID of each
	// episode.
	s.queue.Content = *epsFromDB

	s.queue.len = len(s.queue.Content)

	return nil
}

// CleanQueue removes the entire queue of the player, deleting the contents from memory and database.
func (s *Synchronizer) CleanQueue() error {
	query := "DELETE FROM player_queue;"

	sqlDB := s.db.GetInstance()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := sqlDB.Exec(query)
	if err != nil {
		return err
	}

	s.queue.Content = []QueueEpisode{}

	s.queue.len = 0

	return nil
}

// AddToQueue adds the given QueueEpisode to the actual queue. The parameter `atBeginning` defines if that QueueEpisode
// should be added with the first position or the last one. If there is an error, the ID returned will be -1.
func (s *Synchronizer) AddToQueue(e QueueEpisode, atBeginning bool) (id int, err error) {
	sqlDB := s.db.GetInstance()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.queue.len == 0 {
		/*
		** Case:  the queue is empty, so the episode should be simply added without taking care of the value
		** of the variable `atBeginning`.
		 */

		// Make sure that the position of the episode is 0.
		e.Position = 0

		err := s.insertEpInQueue(e)
		if err != nil {
			return -1, err
		}
	} else if s.queue.len != 0 && atBeginning {
		/*
		** Case: the queue is not empty and the episode should be added at the beginning of it. To do that,
		** the position of the rest of the episodes should be updated (new position = old position + 1).
		 */

		upQuery := "UPDATE player_queue SET position = position + 1;"

		// Update the position of the episodes adding 1 to each one.
		result, err := sqlDB.Exec(upQuery)
		if err != nil {
			return -1, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return -1, err
		}

		if rowsAffected == 0 && s.queue.len != 0 {
			return -1, errorx.InternalError.New("no rows have been affected")
		}

		// Update the queue on memory
		for i := range s.queue.Content {
			s.queue.Content[i].Position = s.queue.Content[i].Position + 1
		}

		// Make sure that the position of the episode is 0.
		e.Position = 0

		err = s.insertEpInQueue(e)
		if err != nil {
			return -1, err
		}
	} else {
		/*
		** Case: the queue is not empty and the episode should be added at the end of it. We don't need to
		** touch the position of the rest of the episodes.
		 */

		// Get the higher position.
		var maxPosition int
		for _, ep := range s.queue.Content {
			if ep.Position > maxPosition {
				maxPosition = ep.Position
			}
		}

		e.Position = maxPosition + 1

		err := s.insertEpInQueue(e)
		if err != nil {
			return -1, err
		}
	}

	idQuery := "SELECT id FROM player_queue ORDER BY id DESC LIMIT 1;"
	rows, err := sqlDB.Query(idQuery)
	if err != nil {
		log.WithError(err).Panic("error when tring to get the ID of the last row in the table 'player_queue'")
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.WithError(err).Error("error when trying to close rows")
		}
	}()

	if !rows.Next() {
		log.WithField("addedEp", e).Panic("there should be at least one row to scan")
	}

	err = rows.Scan(&id)
	if err != nil {
		log.WithError(err).Panic("error when trying to scan the returned ID")
	}

	// Set the obtained ID
	e.ID = id

	// After adding the episode to the database, it should be added to the cached queue.
	s.queue.Content = append(s.queue.Content, e)

	// Update the variable that contains the length.
	s.queue.len++

	return id, nil
}

func (s *Synchronizer) insertEpInQueue(e QueueEpisode) error {
	sqlDB := s.db.GetInstance()
	insertQuery := "INSERT INTO player_queue (podcast_id, episode_id, position) VALUES (?, ?, ?);"

	result, err := sqlDB.Exec(insertQuery, e.PodcastID, e.EpisodeID, e.Position)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.InternalError.New("no rows have been affected")
	}

	return nil
}

// RemoveFromQueue removes the episode with the passed `id` from the queue.
func (s *Synchronizer) RemoveFromQueue(id int) error {
	return nil
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

	var q []QueueEpisode

	for row.Next() {
		var e QueueEpisode
		err = row.Scan(&e.ID, &e.PodcastID, &e.EpisodeID, &e.Position)
		if err != nil {
			return err
		}

		q = append(q, e)
	}

	s.queue.Content = q

	s.queue.len = len(q)

	return nil
}

func (s *Synchronizer) getQueueEpsFromDB() (*[]QueueEpisode, error) {
	query := "SELECT * FROM player_queue;"

	sqlDB := s.db.GetInstance()

	rows, err := sqlDB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.WithError(err).Error("error when trying to close rows")
		}
	}()

	var eps []QueueEpisode

	for rows.Next() {
		var ep QueueEpisode

		err = rows.Scan(&ep.ID, &ep.PodcastID, &ep.EpisodeID, &ep.Position)
		if err != nil {
			return nil, err
		}

		eps = append(eps, ep)
	}

	return &eps, nil
}
