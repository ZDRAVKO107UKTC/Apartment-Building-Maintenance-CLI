package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

type GormIssueRepository struct {
	db *gorm.DB
}

func NewGormIssueRepository(db *gorm.DB) *GormIssueRepository {
	return &GormIssueRepository{db: db}
}

func (r *GormIssueRepository) Create(issue *model.Issue) error {
	return r.db.Create(issue).Error
}

func (r *GormIssueRepository) FindAll() ([]model.Issue, error) {
	var issues []model.Issue
	err := r.db.Find(&issues).Error
	return issues, err
}

func (r *GormIssueRepository) FindByID(id uint) (*model.Issue, error) {
	var issue model.Issue
	if err := r.db.First(&issue, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrIssueNotFound
		}
		return nil, err
	}
	return &issue, nil
}

func (r *GormIssueRepository) Update(issue *model.Issue) error {
	return r.db.Save(issue).Error
}

func (r *GormIssueRepository) Delete(id uint) error {
	return r.db.Delete(&model.Issue{}, id).Error
}
