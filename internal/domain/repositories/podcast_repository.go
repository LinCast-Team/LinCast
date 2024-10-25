package repositories

import (
	"lincast/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PodcastRepository interface {
	GetById(id uint) (*entities.Podcast, error)
	GetByFeed(feedUrl string) (*entities.Podcast, error)
	Create(podcast entities.Podcast) error
	Update(podcast entities.Podcast) error
	Delete(id uint) error
	UpdateSubscriptionStatus(userID uuid.UUID, podcastID uint, subscribed bool) error
}

type podcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) PodcastRepository {
	return &podcastRepository{
		db,
	}
}

func (pr *podcastRepository) GetById(id uint) (*entities.Podcast, error) {
	var p entities.Podcast

	if err := pr.db.First(&p, id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *podcastRepository) Create(podcast entities.Podcast) error {
	if err := pr.db.Create(podcast).Error; err != nil {
		return err
	}

	return nil
}

func (pr *podcastRepository) Update(podcast entities.Podcast) error {
	if err := pr.db.Save(podcast).Error; err != nil {
		return nil
	}

	return nil
}

func (pr *podcastRepository) Delete(id uint) error {
	if err := pr.db.Delete(&entities.Podcast{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (pr *podcastRepository) GetByFeed(feedUrl string) (*entities.Podcast, error) {
	var p = entities.Podcast{
		FeedLink: feedUrl,
	}

	if err := pr.db.First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *podcastRepository) UpdateSubscriptionStatus(userID uuid.UUID, podcastID uint, subscribed bool) error {
	user := entities.User{ID: userID}
	podcast := []entities.Podcast{
		{
			Model: gorm.Model{ID: podcastID},
		},
	}

	association := pr.db.Model(&user).Association("SubscribedTo")

	if subscribed {
		return association.Delete(&podcast)
	} else {
		return association.Append(&podcast)
	}
}
