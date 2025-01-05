package usescase

import (
	"uala.com/timeline-service/internal/domain"
	"uala.com/timeline-service/internal/repository"
)

type GetTimeLineImpl struct {
	timeLineRepository repository.TimeLineRepository
}

func NewGetTimeLineImpl(timeLineRepository repository.TimeLineRepository) GetTimeline {
	return &GetTimeLineImpl{
		timeLineRepository: timeLineRepository,
	}
}

func (g *GetTimeLineImpl) GetTimeline(userId int32, forceSync bool) ([]domain.Timeline, error) {
	var timeline []domain.Timeline
	result, err := g.timeLineRepository.GetTimeline(userId, forceSync)

	for _, item := range result {
		timeline = append(timeline, domain.Timeline{
			ID:        item.ID,
			UserID:    item.UserID,
			UserName:  item.UserName,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		})
	}
	return timeline, err
}
