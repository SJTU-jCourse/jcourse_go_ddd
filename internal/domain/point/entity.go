package point

import "time"

type UserPointRecord struct {
	ItemID      int
	UserID      int
	Point       int
	Description string
	CreatedAt   time.Time
}

type UserPoint struct {
	TotalPoint int
	Records    []UserPointRecord
}

func NewUserPointRecord(userID int, point int, description string) UserPointRecord {
	return UserPointRecord{
		UserID:      userID,
		Point:       point,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
