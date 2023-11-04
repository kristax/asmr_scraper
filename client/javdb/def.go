package javdb

import "context"

type Client interface {
	Get(ctx context.Context, avCode, lang string) (*Detail, error)
}
