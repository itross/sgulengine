package sgulengine

import (
	"context"
)

type ctxComponentLocatorKey int

// CtxComponentLocator is the key to set or get components locator from the Engine Context
const CtxComponentLocator ctxComponentLocatorKey = iota

// EngineContext is the global shared application context.
var EngineContext context.Context
