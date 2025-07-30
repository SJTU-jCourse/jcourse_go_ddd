package point

import (
	"context"
)

type UserPointRepository interface {
	GetUserAllPoints(ctx context.Context, userID int) (*UserPoint, error)
	GetPointRecord(ctx context.Context, itemID int) (*UserPointRecord, error)
	Save(ctx context.Context, point *UserPointRecord) error
}
