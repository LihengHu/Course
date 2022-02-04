package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var (
	DB             *gorm.DB
	LOG            *zap.Logger
	RDB            *redis.Client
	CTX            context.Context
	MutexMembers   sync.Mutex
	MutexCourses   sync.Mutex
	MutexSchedules sync.Mutex
)
