package models

import (
	"CutVido_api/utils"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Name           string
	StartTime      utils.Time
	EndTime        utils.Time
	InputVideoPath string
	OutputDir      string
	StartCut       string
	Duration       string //持续的时间
}
