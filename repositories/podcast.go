package repositories

import (
	"errors"
	"lincast/models"

	"gorm.io/gorm"
)

type PodcastRepository interface {
	GetById(id uint) (*models.Podcast, error)
	GetByFeed(feedUrl string) (*models.Podcast, error)
	Create(podcast models.Podcast) error
	Update(podcast models.Podcast) error
	Delete(id uint) error
	UpdateSubscriptionStatus(userID uint, podcastID uint, subscribed bool) error
}

type podcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) PodcastRepository {
	return &podcastRepository{
		db,
	}
}

func (pr *podcastRepository) GetById(id uint) (*models.Podcast, error) {
	var p models.Podcast

	if err := pr.db.First(&p, id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *podcastRepository) Create(podcast models.Podcast) error {
	if err := pr.db.Create(podcast).Error; err != nil {
		return err
	}

	return nil
}

func (pr *podcastRepository) Update(podcast models.Podcast) error {
	if err := pr.db.Save(podcast).Error; err != nil {
		return nil
	}

	return nil
}

func (pr *podcastRepository) Delete(id uint) error {
	if err := pr.db.Delete(&models.Podcast{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (pr *podcastRepository) GetByFeed(feedUrl string) (*models.Podcast, error) {
	var p = models.Podcast{
		FeedLink: feedUrl,
	}

	// TODO check if this does work. If not, make use of the Where clause
	if err := pr.db.First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *podcastRepository) UpdateSubscriptionStatus(userID uint, podcastID uint, subscribed bool) error {
	return errors.New("not implemented")

	/*
	// TODO Check if this does work
	// Update the subscription status of a user to a podcast, by using the variable subscribed to give the value
	if err := pr.db.Model(&models.User{Model: gorm.Model{ID: userID}}).Preload("SubscribedTo").Delete("subscribedTo", "podcast_id = ?", podcastID).Error; err != nil {
		return err
	}

	return nil*/
}
