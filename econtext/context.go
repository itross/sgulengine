package econtext

import (
	"context"

	"github.com/itross/sgulengine"
)

type ctxComponentLocatorKey int

// CtxComponentLocator is the key to set or get components locator from the Engine Context
const CtxComponentLocator ctxComponentLocatorKey = iota

// EngineContext is the global shared application context.
var EngineContext context.Context

// ComponentLocator helper func to get the Engine components locator from the app context.
func ComponentLocator() (*sgulengine.ComponentLocator, bool) {
	locator, ok := econtext.EngineContext.Value(econtext.CtxComponentLocator)
	if ok {
		return locator.(*sgulengine.ComponentLocator), ok
	}
	return &sgulengine.ComponentLocator{}, false
}
