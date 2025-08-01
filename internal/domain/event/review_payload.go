package event

type ReviewPayload struct {
	ReviewID int    `json:"review_id"`
	UserID   int    `json:"user_id"`
	CourseID int    `json:"course_id"`
	Rating   int    `json:"rating"`
	Content  string `json:"content"`
	Action   string `json:"action"` // "created" or "modified"
}

func (p *ReviewPayload) Type() Type {
	if p.Action == "modified" {
		return TypeReviewModified
	}
	return TypeReviewCreated
}
