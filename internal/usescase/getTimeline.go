package usescase

import (
	"uala.com/timeline-service/internal/domain"
)

type GetTimeline interface {
	GetTimeline(userId int32, forceSync bool) ([]domain.Timeline, error)
}
