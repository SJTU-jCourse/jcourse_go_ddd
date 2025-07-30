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
