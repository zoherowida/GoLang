package crud

import (
	"api/models"
	"api/utils/channels"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// RepositoryUsersCRUD is the struct for the User CRUD
type repositoryTwittersCRUD struct {
	db *gorm.DB
}

func NewRepositoryTwittersCRUD(db *gorm.DB) *repositoryTwittersCRUD {
	return &repositoryTwittersCRUD{db}
}

func (r *repositoryTwittersCRUD) FindAll(offset uint32) ([]models.Twitter, error) {
	var err error
	if offset == 0|1 {
		offset = 0
	} else {
		offset = (offset - 1) * 20
	}
	twitters := []models.Twitter{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Twitter{}).Limit(20).Offset(offset).Find(&twitters).Error

		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return twitters, nil
	}
	return nil, err
}

func (r *repositoryTwittersCRUD) Save(twitter models.Twitter) (models.Twitter, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		err = r.db.Debug().Model(&models.Twitter{}).Create(&twitter).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return twitter, nil
	}
	return models.Twitter{}, err
}

func (r *repositoryTwittersCRUD) FindByID(uid uint32) (models.Twitter, error) {
	var err error
	twitter := models.Twitter{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Twitter{}).Where("id = ?", uid).Take(&twitter).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return twitter, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return models.Twitter{}, errors.New("User Not Found")
	}
	return models.Twitter{}, err
}

func (r *repositoryTwittersCRUD) Update(uid uint32, twitter models.Twitter) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Twitter{}).Where("id = ?", uid).Take(&models.Twitter{}).UpdateColumns(
			map[string]interface{}{
				"tweet":      twitter.Tweet,
				"updated_at": time.Now(),
			},
		)
		ch <- true
	}(done)

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}

		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}

func (r *repositoryTwittersCRUD) Delete(uid uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Twitter{}).Where("id = ?", uid).Take(&models.Twitter{}).Delete(&models.Twitter{})
		ch <- true
	}(done)

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}

		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}
