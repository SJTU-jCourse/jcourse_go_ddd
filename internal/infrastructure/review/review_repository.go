package review

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/review"
	"jcourse_go/internal/infrastructure/entity"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) review.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Get(ctx context.Context, id int) (*review.Review, error) {
	var reviewEntity entity.Review
	result := r.db.WithContext(ctx).
		Preload("Course").
		Preload("User").
		First(&reviewEntity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get review: %w", result.Error)
	}
	return r.toDomainReview(&reviewEntity), nil
}

func (r *reviewRepository) FindBy(ctx context.Context, filter review.ReviewFilter) ([]review.Review, error) {
	var reviewEntitys []entity.Review
	query := r.db.WithContext(ctx).Preload("Course").Preload("User")

	if filter.ReviewID != nil {
		query = query.Where("id = ?", *filter.ReviewID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.CourseID != nil {
		query = query.Where("course_id = ?", *filter.CourseID)
	}
	if filter.MainTeacherID != nil {
		query = query.Joins("JOIN courses ON reviews.course_id = courses.id").
			Where("courses.main_teacher_id = ?", *filter.MainTeacherID)
	}
	if filter.Semester != nil {
		query = query.Where("semester = ?", *filter.Semester)
	}
	if filter.Rating != nil {
		query = query.Where("rating = ?", *filter.Rating)
	}

	result := query.Find(&reviewEntitys)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find reviews: %w", result.Error)
	}

	reviews := make([]review.Review, len(reviewEntitys))
	for i, reviewEntity := range reviewEntitys {
		reviews[i] = *r.toDomainReview(&reviewEntity)
	}

	return reviews, nil
}

func (r *reviewRepository) Save(ctx context.Context, review *review.Review, revision *review.ReviewRevision) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		reviewEntity := r.toORMReview(review)

		if review.ID == 0 {
			if err := tx.Create(reviewEntity).Error; err != nil {
				return fmt.Errorf("failed to create review: %w", err)
			}
			review.ID = reviewEntity.ID
		} else {
			if err := tx.Save(reviewEntity).Error; err != nil {
				return fmt.Errorf("failed to update review: %w", err)
			}
		}

		if revision != nil {
			revisionEntity := r.toORMReviewRevision(revision)
			revisionEntity.ReviewID = review.ID
			if err := tx.Create(revisionEntity).Error; err != nil {
				return fmt.Errorf("failed to create review revision: %w", err)
			}
		}

		return nil
	})
}

func (r *reviewRepository) Delete(ctx context.Context, filter review.ReviewFilter) error {
	query := r.db.WithContext(ctx)
	if filter.ReviewID != nil {
		query = query.Where("id = ?", *filter.ReviewID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.CourseID != nil {
		query = query.Where("course_id = ?", *filter.CourseID)
	}

	result := query.Delete(&entity.Review{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete reviews: %w", result.Error)
	}

	return nil
}

func (r *reviewRepository) SaveReviewAction(ctx context.Context, action *review.ReviewAction) error {
	actionEntity := r.toORMReviewAction(action)
	result := r.db.WithContext(ctx).Create(actionEntity)
	if result.Error != nil {
		return fmt.Errorf("failed to save review action: %w", result.Error)
	}
	return nil
}

func (r *reviewRepository) DeleteReviewAction(ctx context.Context, actionID int) error {
	result := r.db.WithContext(ctx).Delete(&entity.ReviewAction{}, actionID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete review action: %w", result.Error)
	}
	return nil
}

func (r *reviewRepository) GetReviewRevisions(ctx context.Context, reviewID int) ([]review.ReviewRevision, error) {
	var revisionEntitys []entity.ReviewRevision
	result := r.db.WithContext(ctx).Where("review_id = ?", reviewID).Find(&revisionEntitys)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get review revisions: %w", result.Error)
	}

	revisions := make([]review.ReviewRevision, len(revisionEntitys))
	for i, revisionEntity := range revisionEntitys {
		revisions[i] = *r.toDomainReviewRevision(&revisionEntity)
	}
	return revisions, nil
}

// Helper methods to convert between domain and ORM models
func (r *reviewRepository) toDomainReview(reviewEntity *entity.Review) *review.Review {
	return &review.Review{
		ID:       reviewEntity.ID,
		UserID:   reviewEntity.UserID,
		CourseID: reviewEntity.CourseID,
		Rating:   review.NewRating(reviewEntity.Rating),
		Semester: review.NewSemester(reviewEntity.Semester),
		Comment:  reviewEntity.Content,
	}
}

func (r *reviewRepository) toORMReview(review *review.Review) *entity.Review {
	return &entity.Review{
		ID:       review.ID,
		UserID:   review.UserID,
		CourseID: review.CourseID,
		Rating:   int(review.Rating),
		Semester: string(review.Semester),
		Content:  review.Comment,
		Category: "general", // Default category
		IsPublic: true,      // Default to public
	}
}

func (r *reviewRepository) toORMReviewRevision(revision *review.ReviewRevision) *entity.ReviewRevision {
	return &entity.ReviewRevision{
		ID:       revision.ID,
		ReviewID: revision.ReviewID,
		Content:  revision.Comment,
	}
}

func (r *reviewRepository) toDomainReviewRevision(revisionEntity *entity.ReviewRevision) *review.ReviewRevision {
	return &review.ReviewRevision{
		ID:       revisionEntity.ID,
		ReviewID: revisionEntity.ReviewID,
		Comment:  revisionEntity.Content,
		Rating:   3,           // Default rating
		Semester: "2023-Fall", // Default semester
		Grade:    "A",         // Default grade
	}
}

func (r *reviewRepository) toORMReviewAction(action *review.ReviewAction) *entity.ReviewAction {
	return &entity.ReviewAction{
		ID:          0, // Auto-generated
		ReviewID:    action.ReviewID,
		UserID:      action.UserID,
		Action:      action.ActionType,
		Description: "", // Default empty description
	}
}
