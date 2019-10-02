package sgulengine

import (
	"net/http"

	"github.com/itross/sgul"
)

// APIComponent is the default Sgul Engine API component.
// It can be used to create Rest API endpoints for the Sguel Engine based app.
type APIComponent struct {
	BaseComponent
	server *http.Server
}

// NewAPIComponent returns a new API component instance.
func NewAPIComponent() *APIComponent {
	return &APIComponent{
		BaseComponent: BaseComponent{
			uniqueName: "api",
			logger:     sgul.GetLogger(),
		},
	}
}
