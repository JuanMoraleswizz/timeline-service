package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"uala.com/timeline-service/internal/domain"
)

type MockTimelineUseCase struct {
	mock.Mock
}

func (m *MockTimelineUseCase) GetTimeline(userID int32, forceSync bool) ([]domain.Timeline, error) {
	args := m.Called(userID, forceSync)
	return []domain.Timeline{}, args.Error(1)
}

func TestGetTimeline(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockTimelineUseCase)
	handler := NewTimelineHandler(mockUseCase)

	tests := []struct {
		name           string
		userId         string
		forceSync      string
		expectedStatus int
		expectedBody   string
		mockReturn     []domain.Timeline
		mockError      error
	}{
		{
			name:           "Valid request",
			userId:         "1",
			forceSync:      "true",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID": 1, "UserID": 1, "UserName": "John Doe", "Content": "Test content", "CreatedAt": "2021-01-01T00:00:00Z"}]`,
			mockReturn:     []domain.Timeline{{ID: 1, UserID: 1, UserName: "John Doe", Content: "Test content", CreatedAt: "2021-01-01T00:00:00Z"}},
			mockError:      nil,
		},
		{
			name:           "Invalid userId",
			userId:         "abc",
			forceSync:      "true",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"strconv.ParseInt: parsing \"abc\": invalid syntax"`,
			mockReturn:     nil,
			mockError:      nil,
		},
		{
			name:           "Invalid forceSync",
			userId:         "1",
			forceSync:      "notabool",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"strconv.ParseBool: parsing \"notabool\": invalid syntax"`,
			mockReturn:     nil,
			mockError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v2/timeline/%s/%s", tt.userId, tt.forceSync), nil)
			q := req.URL.Query()
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v2/timeline/:userId/:forceSync")
			c.SetParamNames("userId", "forceSync")
			c.SetParamValues(tt.userId, tt.forceSync)
			if tt.mockReturn != nil || tt.mockError != nil {
				userId, _ := strconv.ParseInt(tt.userId, 10, 32)
				forceSync, _ := strconv.ParseBool(tt.forceSync)
				mockUseCase.On("GetTimeline", int32(userId), forceSync).Return(tt.mockReturn, tt.mockError)
			}

			err := handler.GetTimeline(c)
			if err != nil {
				t.Errorf("handler returned an error: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
