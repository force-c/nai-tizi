package bootstrap

import (
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/platform/websocket"
)

// Bootstrap wires optional scaffold components without binding concrete business logic.
type Bootstrap struct{}

// New creates the bootstrap coordinator for optional runtime integrations.
func New(c container.Container) (*Bootstrap, error) {
	return &Bootstrap{}, nil
}

// ConfigureWebSocketHandler applies websocket runtime options.
func (b *Bootstrap) ConfigureWebSocketHandler(handler *websocket.Handler, cfg any, log logger.Logger) {
	// The scaffold keeps this hook for business projects to extend websocket behavior.
}
