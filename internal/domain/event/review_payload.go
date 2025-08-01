package event

type ReviewPayload struct {
	ReviewID string `json:"review_id"`
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
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
