// Package handler contains event handlers for domain events.
// This package has been refactored to follow proper DDD architecture:
// - Event handlers are split into individual files
// - Handlers use application services instead of repositories directly
// - Each handler has a single responsibility
//
// See individual files:
// - review_event_handler.go: Handles review-related events
// - point_event_handler.go: Handles point awarding events
// - statistics_event_handler.go: Handles statistics update events
// - noop_handler.go: No-operation handler
// - registration.go: Event handler registration
package handler
