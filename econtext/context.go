package econtext

import "context"

type ctxEngineComponentsKey int

// CtxComponents is the key to set or get components list from the Engine Context
const CtxComponents ctxEngineComponentsKey = iota

// EngineContext is the global shared application context.
var EngineContext context.Context
