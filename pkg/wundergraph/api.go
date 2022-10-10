package wundergraph

import "context"

type Wdg interface {
	ReloadConfig(ctx context.Context) error
	ReloadOpertionTs(ctx context.Context) error
	ReloadServerTs(ctx context.Context) error
}
