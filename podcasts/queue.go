package podcasts

import (
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"time"
)

type Job struct {
	Podcast *Podcast
}

type UpdateQueue struct {
	dbInstance *Database
	q          chan Job
}

func NewUpdateQueue(db *Database, length int) (*UpdateQueue, error) {
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

/*
func NewJob(p *Podcast) *Job {

}
*/
func (q *UpdateQueue) worker(id int) {
	log.WithField("worker", id).Debug("Starting worker")

	for {
		job := <-q.q
		receivedTime := time.Now()

		log.WithFields(log.Fields{
			"worker":      id,
			"podcastID":   job.Podcast.ID,
			"podcastFeed": job.Podcast.FeedLink,
		}).Info("New job received. Updating podcast")

		eps, err := job.Podcast.GetEpisodes()
		if err != nil {
			log.WithFields(log.Fields{
				"worker":      id,
				"podcastID":   job.Podcast.ID,
				"podcastFeed": job.Podcast.FeedLink,
				"error":       errorx.Decorate(err, "error when getting episodes of the podcast"),
			}).Error("Can't get the episodes of the podcast")

			continue
		}

		for _, e := range *eps {
			exists, err := q.dbInstance.EpisodeExists(e.GUID)
			if err != nil {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
					"error":       errorx.Decorate(err, "error when trying to check if an episode exists"),
				}).Error("Can't get the episodes of the podcast")

				continue
			}

			if exists {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
				}).Debug("Episode already on the database, skipping")

				continue
			}

			err = q.dbInstance.InsertEpisode(&e)
			if err != nil {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"episodeGUID": e.GUID,
					"error":       errorx.Decorate(err, "error when trying to save the episode on the database"),
				}).Error("Can't save the episode on the database")

				continue
			}

			err = q.dbInstance.UpdatePodcastLastCheck(job.Podcast.ID, time.Now())
			if err != nil {
				log.WithFields(log.Fields{
					"worker":      id,
					"podcastID":   job.Podcast.ID,
					"podcastFeed": job.Podcast.FeedLink,
					"error": errorx.Decorate(err, "error when trying to update the column "+
						"'last_check' in the database"),
				}).Error("Can't update the LastCheck time in the database")

				continue
			}

			log.WithFields(log.Fields{
				"worker":         id,
				"podcastID":      job.Podcast.ID,
				"podcastFeed":    job.Podcast.FeedLink,
				"updateDuration": time.Since(receivedTime).String(),
			}).Info("Job finished. Podcast updated correctly")
		}
	}
}
