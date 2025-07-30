package review

import "time"

func NewRevisionFromReview(r *Review) ReviewRevision {
	return ReviewRevision{
		ReviewID:  r.ID,
		UserID:    r.UserID,
		Comment:   r.Comment,
		Semester:  r.Semester,
		Grade:     r.Grade,
		CreatedAt: time.Now(),
	}
}

func NewReview(courseID int, userID int, c *ReviewContent) Review {
	return Review{
		CourseID:  courseID,
		UserID:    userID,
		Comment:   c.Comment,
		Rating:    NewRating(c.Rating),
		Semester:  NewSemester(c.Semester),
		Grade:     c.Grade,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
