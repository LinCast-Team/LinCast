package queue

import (
	"time"

	"lincast/database"
	"lincast/podcasts"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

// Job returns a new job to be processed by a worker of an active UpdateQueue. The channel Job.Done can be used to know
// when that job has been processed and it shouldn't be used to send something, just to receive.
type Job struct {
	Podcast *podcasts.Podcast
	Done    chan struct{}
}

type UpdateQueue struct {
	dbInstance *database.Database
	q          chan Job
}

func NewUpdateQueue(db *database.Database, length int) (*UpdateQueue, error) {
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

func NewJob(p *podcasts.Podcast) *Job {
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

		eps, err := job.Podcast.GetEpisodes()
		if err != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Can't get the episodes of the podcast")

			continue
		}

		for _, e := range *eps {
			<-rateLimiter.C

			exists, err := q.dbInstance.EpisodeExists(e.GUID)
			if err != nil {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
					"error":       errorx.EnsureStackTrace(err),
				}).Error("Can't check if the episode already exists")

				continue
			}

			if exists {
				continue
			}

			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"episodeGUID": e.GUID,
			}).Debug("Episode is not in the database, storing")

			err = q.dbInstance.InsertEpisode(&e)
			if err != nil {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
					"error":       errorx.EnsureStackTrace(err),
				}).Error("Can't save the episode on the database")

				continue
			}
		}

		err = q.dbInstance.UpdatePodcastLastCheck(job.Podcast.ID, time.Now())
		if err != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Can't update the LastCheck time in the database")

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
