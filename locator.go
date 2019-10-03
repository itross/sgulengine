package sgulengine

import "sync"

type (
	// componentRegistry is the type for the map of registred components.
	componentRegistry map[string]Component

	// ComponentLocator is the locator for registered components instances by name.
	ComponentLocator struct {
		sync.RWMutex
		cReg *componentRegistry
	}
)

// Get returns a component instance from the components registry.
func (locator *ComponentLocator) Get(cname string) Component {
	locator.RLock()
	defer locator.RUnlock()

	return (*locator.cReg)[cname]
}

// GetComponentLocator helper func to get the Engine components locator from the app context.
func GetComponentLocator() *ComponentLocator {
	return EngineContext.Value(CtxComponentLocator).(*ComponentLocator)

}
