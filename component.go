// Copyright 2019 Luca Stasio <joshuagame@gmail.com>
// Copyright 2019 IT Resources s.r.l.
//
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sgulengine defines the Sgul Engine structure and functionalities.
package sgulengine

import "github.com/itross/sgul"

// Component is the main sgul engine component interface.
type Component interface {
	Name() string
	SetName(name string)
	Logger() *sgul.Logger
	SetLogger(*sgul.Logger)
	Configure(config interface{}) error
	Start(*Engine) error
	Shutdown()
}

// BaseComponent .
type BaseComponent struct {
	uniqueName string
	logger     *sgul.Logger
}

// Name returns the unique component name within the Engine.
func (c *BaseComponent) Name() string {
	return c.uniqueName
}

// SetName sets the unique component name.
func (c *BaseComponent) SetName(name string) {
	c.uniqueName = name
}

// SetLogger sets the logger dependency during the component registration phase.
func (c *BaseComponent) SetLogger(logger *sgul.Logger) {
	c.logger = logger
}

// Logger returns the component Logger instance.
func (c *BaseComponent) Logger() *sgul.Logger {
	return c.logger
}
