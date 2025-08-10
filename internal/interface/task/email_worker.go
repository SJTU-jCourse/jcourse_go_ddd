package task

import (
	"context"
	"log"
	"time"

	"jcourse_go/internal/app"
)

// EmailWorker handles email sending tasks
type EmailWorker struct {
	serviceContainer *app.ServiceContainer
}

func NewEmailWorker(serviceContainer *app.ServiceContainer) *EmailWorker {
	return &EmailWorker{
		serviceContainer: serviceContainer,
	}
}

func (w *EmailWorker) Start(ctx context.Context) {
	log.Println("Email worker started")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker stopped")
			return
		case <-ticker.C:
			// Process pending email tasks
			// Check for pending email tasks and send them using the email service
		}
	}
}
