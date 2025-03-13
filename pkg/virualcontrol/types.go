package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type Options struct {
	BaseOptions
	OnHandler    OnHandler[string]
	DefaultValue string
}

type BaseOptions struct {
	Queries *db.Queries
	Client  wb.ClientInterface
	Device  string
	Name    string
	Meta    control.Meta
}

type OnHandler[T comparable] func(payload OnHandlerPayload[T])

type OnHandlerPayload[T comparable] struct {
	Set   func(value T)
	Value T
	Error error
}
