package podcasts

import (
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

type Job struct {
	PodcastID int

	// TODO
}

type Queue struct {
	q chan Job
}

func NewQueue(length int) (*Queue, error) {
	if length < 1 {
		return nil, errorx.IllegalArgument.New("the length of the queue should be at least 1")
	}

	q := &Queue{
		q: make(chan Job),
	}

	for i := 0; i < length; i++ {
		go worker(i, q.q)
	}

	return q, nil
}

func (q *Queue) Send(job *Job) {
	q.q <- *job
}

func worker(id int, q chan Job) {
	log.Debugln("Starting worker", id)

	for {
		job := <-q
		log.Debugf("[worker %d] New job received. Updating podcast %d\n", id, job.PodcastID)

		// TODO Update feed
	}
}
