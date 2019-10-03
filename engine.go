// Copyright 2019 Luca Stasio <joshuagame@gmail.com>
// Copyright 2019 IT Resources s.r.l.
//
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sgulengine defines the Sgul Engine structure and functionalities.
package sgulengine

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/itross/sgul"
	"github.com/itross/sgulengine/econtext"
)

type (
	// componentRegistry is the type for the map of registred components.
	componentRegistry map[string]Component

	// ComponentLocator is the locator for registered components instances by name.
	ComponentLocator struct {
		sync.RWMutex
		cReg *componentRegistry
	}

	// Engine is the sgul app engine main structure.
	Engine struct {
		// TODO: use a decoupled components registry
		cReg    componentRegistry
		locator *ComponentLocator
		stopch  chan os.Signal
		logger  *sgul.Logger
		ctx     context.Context
	}
)

// Get returns a component instance from the components registry.
func (locator *ComponentLocator) Get(cname string) Component {
	locator.RLock()
	defer locator.RUnlock()

	return (*locator.cReg)[cname]
}

// New returns a new sgul Engine instance.
func New() *Engine {
	registry := make(componentRegistry)
	e := &Engine{
		cReg:    registry,
		locator: &ComponentLocator{cReg: &registry},
		stopch:  make(chan os.Signal),
		logger:  sgul.GetLogger(),
		ctx:     context.Background(),
	}
	// set up os signal notifications
	signal.Notify(e.stopch, syscall.SIGTERM)
	signal.Notify(e.stopch, syscall.SIGINT)

	// start a go-func to trigger os signals and gently shutdown the Machinery
	go func() {
		sig := <-e.stopch
		e.logger.Infof("caught sig: %+v", sig)

		e.logger.Info("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)

		e.Shutdown()
		e.logger.Info("now Engine is down")
		e.logger.Info("Bye!")
		os.Exit(0)
	}()

	return e
}

// With registers one or more sgul Components with the sgul Engine.
func (e *Engine) With(components ...Component) *Engine {
	var cname string

	for _, component := range components {
		cname = component.Name()
		if e.cReg[cname] != nil {
			e.logger.Warnf("component %s already registered", cname)
			continue
		}
		e.logger.Infof("registering %s component", cname)
		component.SetLogger(e.logger)
		e.cReg[cname] = component
	}

	return e
}

// ForEachComponent executes a function on each of the Engine components.
func (e *Engine) ForEachComponent(fn func(component Component) error) {
	for cname, component := range e.cReg {
		if err := fn(component); err != nil {
			e.logger.Errorf("error on component %s: %s", cname, err.Error())
			panic(err)
		}
	}
}

// Configure configure each registered component.
func (e *Engine) Configure() {
	e.logger.Info("configuring Engine components")
	e.ForEachComponent(e.configureComponent)
}

// Start starts each registered components.
func (e *Engine) Start() {
	e.logger.Info("starting Engine components")
	e.ForEachComponent(e.startComponent)
}

// Run will starts up the Engine. All starts here!
func (e *Engine) Run() {
	e.Configure()
	e.Start()
	econtext.EngineContext = context.WithValue(e.ctx, econtext.CtxComponentLocator, &ComponentLocator{cReg: &e.cReg})
	e.logger.Info("component locator has been set into app context")
}

// RunAndWait wil starts up the Engine and wait for shutdown.
func (e *Engine) RunAndWait() {
	e.Run()
	select {}
}

// shutdown a single component.
func (e *Engine) shutdownComponent(component Component) error {
	e.logger.Infof("shutting down %s component", component.Name())
	component.Shutdown()
	return nil
}

// configure a single component.
func (e *Engine) configureComponent(component Component) error {
	cname := component.Name()
	e.logger.Infof("configuring %s component", cname)

	cconf := sgul.GetComponentConfig(cname)
	if cconf == nil {
		return errors.New("no configuration found")
	}

	return component.Configure(cconf)
}

func (e *Engine) startComponent(component Component) error {
	cname := component.Name()
	e.logger.Infof("starting %s component", cname)

	return component.Start(e)
}

// Shutdown will call shutdown func on each registered component.
func (e *Engine) Shutdown() {
	e.logger.Info("shutting down the Engine")
	e.ForEachComponent(e.shutdownComponent)
}

// Component returns a component instance.
func (e *Engine) Component(name string) Component {
	return e.cReg[name]
}
