## Sgul ENGINE

Example:
```
package main

import (
	"github.com/itross/sgul"
	"github.com/itross/sgulengine/sgulengine"
)

// APIComponent .
type APIComponent struct {
	sgulengine.BaseComponent
}

// NewAPIComponent .
func NewAPIComponent() *APIComponent {
	ic := &APIComponent{BaseComponent: sgulengine.BaseComponent{}}
	ic.SetName("api")
	ic.SetLogger(sgul.GetLogger())
	return ic
}

// Configure .
func (ic *APIComponent) Configure(conf interface{}) error {
	ic.Logger().Info("* api configured")
	return nil
}

// Start .
func (ic *APIComponent) Start(e *sgulengine.Engine) error {
	ic.Logger().Info("* api started")
	return nil
}

// Shutdown .
func (ic *APIComponent) Shutdown() {
	ic.Logger().Info("* api shutted down")
}

func main() {
	logger := sgul.GetLogger()

	e := sgulengine.New()
	e.With(NewAPIComponent())

	go func() {
		logger.Info("Started goroutine")
		api := (e.Component("api"))
		logger.Infof("using %s-component", api.Name())
	}()

	e.RunAndWait()
}
```