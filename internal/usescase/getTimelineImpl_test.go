package usescase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"uala.com/timeline-service/internal/entity"
)

// Mock repository
type MockTimeLineRepository struct {
	mock.Mock
}

func (m *MockTimeLineRepository) GetTimeline(userId int32, forceSync bool) ([]entity.Timeline, error) {
	args := m.Called(userId, forceSync)
	return args.Get(0).([]entity.Timeline), args.Error(1)
}

func TestGetTimeline(t *testing.T) {
	mockRepo := new(MockTimeLineRepository)
	service := NewGetTimeLineImpl(mockRepo)

	timelineData := []entity.Timeline{
		{
			ID:        1,
			UserID:    123,
			UserName:  "John Doe",
			Content:   "Test content",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetTimeline", int32(123), false).Return(timelineData, nil)

		result, err := service.GetTimeline(123, false)

		assert.NoError(t, err)
		assert.Equal(t, len(timelineData), len(result))
		mockRepo.AssertExpectations(t)
	})

}
