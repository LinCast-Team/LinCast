package update

import (
	"errors"
	"time"

	"lincast/models"
	"lincast/podcasts"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Job returns a new job to be processed by a worker of an active UpdateQueue. The channel Job.Done can be used to know
// when that job has been processed and it shouldn't be used to send something, just to receive.
type Job struct {
	Podcast *models.Podcast
	Done    chan struct{}
}

type UpdateQueue struct {
	dbInstance *gorm.DB
	q          chan Job
}

func NewUpdateQueue(db *gorm.DB, length int) (*UpdateQueue, error) {
	if length < 1 {
		return nil, errorx.IllegalArgument.New("the length of the queue should be at least 1")
	}

	if db == nil {
		return nil, errorx.IllegalState.New("the instance of the database is nil")
	}

	q := UpdateQueue{
		q:          make(chan Job),
		dbInstance: db,
	}

	for i := 0; i < length; i++ {
		go q.worker(i)
	}

	return &q, nil
}

func (q *UpdateQueue) Send(job *Job) {
	q.q <- *job
}

func NewJob(p *models.Podcast) *Job {
	j := Job{
		Podcast: p,
		Done:    make(chan struct{}),
	}

	return &j
}

func (q *UpdateQueue) worker(id int) {
	log.WithField("worker", id).Debug("Worker started")

	// Limit the frequency by which each episode is processed.
	rateLimiter := time.NewTicker(time.Millisecond * 300)
	defer rateLimiter.Stop()

	for {
		job := <-q.q
		receivedTime := time.Now()

		log.WithFields(log.Fields{
			"worker":      id,
			"podcastID":   job.Podcast.ID,
			"podcastFeed": job.Podcast.FeedLink,
		}).Info("New job received")

		_, feed, err := podcasts.GetPodcastData(job.Podcast.FeedLink)
		if err != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Error when trying to obtain the feed")

			continue
		}

		eps, err := podcasts.GetEpisodes(feed)
		if err != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Error on episodes parsing")

			continue
		}

		for _, e := range *eps {
			<-rateLimiter.C

			// Check if the episode is already on the table.
			result := q.dbInstance.Where("guid = ?", e.GUID).First(&models.Episode{})
			if result.Error != nil {
				// The only error that we expect to get here is one of type `gorm.ErrRecordNotFound` (which means
				// basically that the episode is not stored on the database). So, if we get another type of error
				// we should log it and skip the rest of the iteration.
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					log.WithFields(log.Fields{
						"worker":      id,
						"podcastID":   job.Podcast.ID,
						"podcastFeed": job.Podcast.FeedLink,
						"episodeGUID": e.GUID,
						"error":       errorx.EnsureStackTrace(err),
					}).Error("Can't check if the episode already exists")

					continue
				}
			} else {
				// If there are no errors, then the episode with the given GUID exists so we should skip storage.
				continue
			}

			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"episodeGUID": e.GUID,
			}).Debug("Episode is not in the database, storing")

			// Set the ID of the parent podcast before store the episode (if is not already on the db).
			e.PodcastID = job.Podcast.ID

			result = q.dbInstance.Create(&e)
			if result.Error != nil || result.RowsAffected == 0 {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
					"error":       errorx.EnsureStackTrace(result.Error),
				}).Error("The new episode can't be stored")

				continue
			}
		}

		result := q.dbInstance.Model(job.Podcast).Update("last_check", time.Now())
		if result.Error != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.EnsureStackTrace(result.Error),
			}).Error("The last_check time of the podcast can't be updated")

			continue
		}

		log.WithFields(log.Fields{
			"worker":      id,
			"podcastID":   job.Podcast.ID,
			"podcastFeed": job.Podcast.FeedLink,
		}).Debug("Job finished, sending notification through the channel Job.Done")

		// Notify that the job has been processed without blocking.
		select {
		case job.Done <- struct{}{}:
		default:
		}

		log.WithFields(log.Fields{
			"worker":         id,
			"podcastID":      job.Podcast.ID,
			"podcastFeed":    job.Podcast.FeedLink,
			"updateDuration": time.Since(receivedTime).String(),
		}).Info("Podcast updated correctly")
	}
}
