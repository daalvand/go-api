package repositories

import (
	"github.com/daalvand/go-api/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *models.Event) error
	FindAll() ([]models.Event, error)
	FindByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(event *models.Event) error
	Register(event *models.Event, userID uint) error
	CancelRegistration(event *models.Event, userID uint) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (er *eventRepository) Create(event *models.Event) error {

	return er.db.Create(event).Error
}

func (er *eventRepository) FindAll() ([]models.Event, error) {
	var events []models.Event
	if err := er.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (er *eventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	if err := er.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (er *eventRepository) Update(event *models.Event) error {
	return er.db.Save(event).Error
}

func (er *eventRepository) Delete(event *models.Event) error {
	return er.db.Delete(event).Error
}

func (er *eventRepository) Register(event *models.Event, userID uint) error {
	registration := &models.Registration{
		EventID: event.ID,
		UserID:  userID,
	}
	return er.db.Create(registration).Error
}

func (er *eventRepository) CancelRegistration(event *models.Event, userID uint) error {
	return er.db.Delete(&models.Registration{}, "event_id = ? AND user_id = ?", event.ID, userID).Error
}
