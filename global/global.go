package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	LOG *zap.Logger
)
