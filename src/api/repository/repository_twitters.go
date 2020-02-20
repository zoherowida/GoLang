package repository

import (
	"api/models"
)

// UserRepository is the interface User CRUD
type TwitterRepository interface {
	Save(models.Twitter) (models.Twitter, error)
	FindAll(uint32) ([]models.Twitter, error)
	FindByID(uint32) (models.Twitter, error)
	Update(uint32, models.Twitter) (int64, error)
	Delete(uint32) (int64, error)
}
