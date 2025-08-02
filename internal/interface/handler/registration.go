package handler

import (
	"jcourse_go/internal/application/point/command"
	"jcourse_go/internal/domain/event"
)

// RegisterEventHandlers registers all event handlers with the event bus
func RegisterEventHandlers(eventBus event.EventBusPublisher, pointService command.PointCommandService) error {
	reviewHandler := NewReviewEventHandler()
	pointHandler := NewPointEventHandler(pointService)
	statsHandler := NewStatisticsEventHandler()

	if err := eventBus.Register(event.TypeReviewCreated, reviewHandler); err != nil {
		return err
	}
	if err := eventBus.Register(event.TypeReviewModified, reviewHandler); err != nil {
		return err
	}
	if err := eventBus.Register(event.TypeReviewCreated, pointHandler); err != nil {
		return err
	}
	if err := eventBus.Register(event.TypeReviewModified, statsHandler); err != nil {
		return err
	}

	return nil
}