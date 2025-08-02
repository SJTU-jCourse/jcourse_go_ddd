package handler

import (
	"context"
	"testing"
	
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/event"
	"jcourse_go/pkg/apperror"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPointCommandService is a mock implementation of the point command service
type MockPointCommandService struct {
	mock.Mock
}

func (m *MockPointCommandService) CreatePoint(commonCtx *common.CommonContext, userID int, amount int, reason string) error {
	args := m.Called(commonCtx, userID, amount, reason)
	return args.Error(0)
}

func (m *MockPointCommandService) Transaction(commonCtx *common.CommonContext, fromUserID int, toUserID int, amount int, reason string) error {
	args := m.Called(commonCtx, fromUserID, toUserID, amount, reason)
	return args.Error(0)
}

func (m *MockPointCommandService) AwardPointsForReview(commonCtx *common.CommonContext, userID int, reviewID int) error {
	args := m.Called(commonCtx, userID, reviewID)
	return args.Error(0)
}

func TestPointEventHandler_HandleReviewCreated(t *testing.T) {
	// Setup
	mockService := new(MockPointCommandService)
	handler := NewPointEventHandler(mockService)
	
	// Create test event
	payload := &event.ReviewPayload{
		ReviewID: 123,
		UserID:   456,
		CourseID: 789,
		Rating:   5,
		Content:  "Great course!",
		Action:   "created",
	}
	
	event := event.NewBaseEvent(event.TypeReviewCreated, payload)
	
	// Expect the AwardPointsForReview method to be called with correct parameters
	mockService.On("AwardPointsForReview", mock.MatchedBy(func(ctx *common.CommonContext) bool {
		return ctx.User.UserID == 0 && ctx.User.Role == common.RoleAdmin // System user
	}), 456, 123).Return(nil)
	
	// Execute
	err := handler.Handle(context.Background(), event)
	
	// Assert
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// MockInvalidPayload is a mock payload that doesn't match the expected type
type MockInvalidPayload struct{}

func (m *MockInvalidPayload) Type() event.Type {
	return event.Type(999) // Invalid type
}

func TestPointEventHandler_HandleInvalidPayload(t *testing.T) {
	// Setup
	mockService := new(MockPointCommandService)
	handler := NewPointEventHandler(mockService)
	
	// Create event with invalid payload type
	invalidPayload := &MockInvalidPayload{}
	event := event.NewBaseEvent(event.TypeReviewCreated, invalidPayload)
	
	// Execute
	err := handler.Handle(context.Background(), event)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload type")
	mockService.AssertNotCalled(t, "AwardPointsForReview")
}

func TestPointEventHandler_HandleReviewModified(t *testing.T) {
	// Setup
	mockService := new(MockPointCommandService)
	handler := NewPointEventHandler(mockService)
	
	// Create test event for modified review
	payload := &event.ReviewPayload{
		ReviewID: 123,
		UserID:   456,
		CourseID: 789,
		Rating:   4,
		Content:  "Updated review",
		Action:   "modified",
	}
	
	event := event.NewBaseEvent(event.TypeReviewModified, payload)
	
	// Execute
	err := handler.Handle(context.Background(), event)
	
	// Assert
	assert.NoError(t, err)
	// AwardPointsForReview should not be called for modified reviews
	mockService.AssertNotCalled(t, "AwardPointsForReview")
}

func TestPointEventHandler_HandleSaveError(t *testing.T) {
	// Setup
	mockService := new(MockPointCommandService)
	handler := NewPointEventHandler(mockService)
	
	// Create test event
	payload := &event.ReviewPayload{
		ReviewID: 123,
		UserID:   456,
		CourseID: 789,
		Rating:   5,
		Content:  "Great course!",
		Action:   "created",
	}
	
	event := event.NewBaseEvent(event.TypeReviewCreated, payload)
	
	// Expect the AwardPointsForReview method to be called and return an error
	mockService.On("AwardPointsForReview", mock.AnythingOfType("*common.CommonContext"), 456, 123).Return(apperror.ErrDB)
	
	// Execute
	err := handler.Handle(context.Background(), event)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockService.AssertExpectations(t)
}