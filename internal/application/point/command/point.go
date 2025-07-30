package command

import "context"

type PointCommandService interface {
	CreatePoint(ctx context.Context) error
	Transaction(ctx context.Context) error
}
