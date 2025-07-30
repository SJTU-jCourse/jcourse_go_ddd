package review

type WriteReviewCommand struct {
	CourseID int
	ReviewContent
}

type UpdateReviewCommand struct {
	ReviewID int
	ReviewContent
}

type DeleteReviewCommand struct {
	ReviewID int
}
