package task

import (
	"context"
	"log"
	"time"

	"jcourse_go/internal/app"
	"jcourse_go/internal/application/statistics/service"
)

const (
	StatisticsWorkerTicker = 24 * time.Hour
	StatisticsWorkerDelay  = 5 * time.Second
	DateFormat            = "2006-01-02"
)

// StatisticsWorker handles periodic statistics calculation
type StatisticsWorker struct {
	serviceContainer  *app.ServiceContainer
	dailyStatsService service.DailyStatisticsService
}

func NewStatisticsWorker(serviceContainer *app.ServiceContainer) *StatisticsWorker {
	return &StatisticsWorker{
		serviceContainer:  serviceContainer,
		dailyStatsService: serviceContainer.DailyStatisticsService,
	}
}

func (w *StatisticsWorker) Start(ctx context.Context) {
	log.Println("Statistics worker started")

	// Run daily statistics calculation at midnight
	ticker := time.NewTicker(StatisticsWorkerTicker)
	defer ticker.Stop()

	// Calculate statistics for today on startup
	go w.calculateDailyStatistics(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Statistics worker stopped")
			return
		case <-ticker.C:
			w.calculateDailyStatistics(ctx)
		}
	}
}

func (w *StatisticsWorker) calculateDailyStatistics(ctx context.Context) {
	log.Println("Calculating daily statistics...")

	// Calculate statistics for yesterday (complete day)
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

	err := w.dailyStatsService.CalculateAndSaveDailyStatistics(ctx, yesterday)
	if err != nil {
		log.Printf("Failed to calculate daily statistics for %v: %v", yesterday.Format(DateFormat), err)
		return
	}

	log.Printf("Successfully calculated daily statistics for %v", yesterday.Format(DateFormat))
}
