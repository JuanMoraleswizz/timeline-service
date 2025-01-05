package repository

import (
	"uala.com/timeline-service/internal/entity"
)

type TimeLineRepository interface {
	GetTimeline(userId int32, forceSync bool) ([]entity.Timeline, error)
}
