package identifier

import "context"

type SnowflakeService interface {
	Generate(ctx context.Context) (id int64, err error)
}
