package asmr_one

import "context"

type Client interface {
	GetWorkInfo(ctx context.Context, rj string) (*WorkInfoResponse, error)
}
