package models

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Name           string
	StartTime      string
	EndTime        string
	InputVideoPath string
	OutputDir      string
	StartCut       string
	Duration       string //持续的时间
}
