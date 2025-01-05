package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"uala.com/timeline-service/internal/usescase"
)

type TimelineHandler struct {
	timeLineUsesCase usescase.GetTimeline
}

func NewTimelineHandler(timeLineUsesCase usescase.GetTimeline) TimelineHandler {
	return TimelineHandler{
		timeLineUsesCase: timeLineUsesCase,
	}
}

func (t *TimelineHandler) GetTimeline(c echo.Context) error {
	user, err := strconv.ParseInt(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	forceSync, err := strconv.ParseBool(c.Param("forceSync"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	timeline, err := t.timeLineUsesCase.GetTimeline(int32(user), forceSync)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, timeline)
}
