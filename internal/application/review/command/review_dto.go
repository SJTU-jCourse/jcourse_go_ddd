package command

import (
	"jcourse_go/internal/domain/review"
)

type WriteReviewCommand struct {
	CourseID int `json:"course_id"`
	review.ReviewContent
}

type UpdateReviewCommand struct {
	ReviewID int `json:"review_id"`
	review.ReviewContent
}

type DeleteReviewCommand struct {
	ReviewID int `json:"review_id"`
}
