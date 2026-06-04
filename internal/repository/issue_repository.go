package repository

import (
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/database"
	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

func CreateIssue(issue *model.Issue) error {
	return database.DB.Create(issue).Error
}

func FindAllIssues() ([]model.Issue, error) {
	var issues []model.Issue
	err := database.DB.Find(&issues).Error
	return issues, err
}

func FindIssueByID(id uint) (*model.Issue, error) {
	var issue model.Issue
	if err := database.DB.First(&issue, id).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func UpdateIssue(issue *model.Issue) error {
	return database.DB.Save(issue).Error
}

func DeleteIssue(id uint) error {
	return database.DB.Delete(&model.Issue{}, id).Error
}
