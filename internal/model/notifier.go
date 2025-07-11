package model

import "context"

type Notifier interface {
	Notify(ctx context.Context, room, event, data string) error
	Close() error
}
