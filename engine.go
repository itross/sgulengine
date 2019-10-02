// Copyright 2019 Luca Stasio <joshuagame@gmail.com>
// Copyright 2019 IT Resources s.r.l.
//
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sgulengine defines the Sgul Engine structure and functionalities.
package sgulengine

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/itross/sgul"
)

// Engine is the sgul app engine main structure.
type Engine struct {
	// TODO: use a decoupled components registry
	components map[string]Component
	stopch     chan os.Signal
	logger     *sgul.Logger
}

// New returns a new sgul Engine instance.
func New() *Engine {
	e := &Engine{
		components: make(map[string]Component),
		stopch:     make(chan os.Signal),
		logger:     sgul.GetLogger(),
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
		if e.components[cname] != nil {
			e.logger.Warnf("component %s already registered", cname)
			continue
		}
		e.logger.Infof("registering %s component", cname)
		component.SetLogger(e.logger)
		e.components[cname] = component
	}

	return e
}

// ForEachComponent executes a function on each of the Engine components.
func (e *Engine) ForEachComponent(fn func(component Component) error) {
	for cname, component := range e.components {
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

	cconf := sgul.Get(strings.ToLower(cname))
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
	return e.components[name]
}
