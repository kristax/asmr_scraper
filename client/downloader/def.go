package downloader

import "context"

type Client interface {
	Download(ctx context.Context, url string) ([]byte, error)
}
