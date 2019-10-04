package sgulengine

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	chilogger "github.com/766b/chi-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/itross/sgul"
)

// APIComponent is the default Sgul Engine API component.
// It can be used to create Rest API endpoints for the Sguel Engine based app.
type APIComponent struct {
	BaseComponent
	config      sgul.API
	server      *http.Server
	router      chi.Router
	controllers []sgul.RestController
	middlewares []func(http.Handler) http.Handler
}

// Configure willl configure the api component with its internal server and router.
// All registered Rest Controller will be coupled to the relative route.
// Then the main server handler will be gained to the http server and started.
func (api *APIComponent) Configure(conf interface{}) error {
	api.config = conf.(sgul.API)
	api.configureRouter()
	return nil
}

func (api *APIComponent) configureRouter() {
	api.router = chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: api.config.Cors.Origin,
		AllowedMethods: api.config.Cors.Methods,
		AllowedHeaders: api.config.Cors.Headers,
	})
	api.middlewares = append(api.middlewares, cors.Handler)
	api.router.Use(api.middlewares...)
}

// Use sets middlewares for this API routes.
func (api *APIComponent) Use(middlewares ...func(http.Handler) http.Handler) *APIComponent {
	api.middlewares = append(api.middlewares, middlewares...)
}

func (api *APIComponent) registerRoutes() {
	if len(api.controllers) > 0 {
		// register controllers routes
		api.router.Route(api.config.Endpoint.BaseRoutingPath, func(r chi.Router) {
			for _, controller := range api.controllers {
				r.Mount(controller.BasePath(), controller.Router())
			}
		})

		// log out configured routes
		walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			route = strings.Replace(route, "/*/", "/", -1)
			api.logger.Infow("initialized", "method", method, "route", route)
			return nil
		}

		if err := chi.Walk(api.router, walker); err != nil {
			api.logger.Panicf("error: %s\n", err.Error())
		}
		api.logger.Info("all api routes set up")

		return
	}

	api.Logger().Warn("no controller registered with the APIComponent... no route installed!")
}

// Start willl start the API Server after initialization.
func (api *APIComponent) Start(e *Engine) error {
	api.registerRoutes()
	addr := fmt.Sprintf(":%d", api.config.Endpoint.Port)
	api.server = &http.Server{
		Addr:    addr,
		Handler: api.router,
	}

	go func() {
		api.server.ListenAndServe()
	}()

	return nil
}

// Shutdown will stop serving the API.
func (api *APIComponent) Shutdown() {
	if err := api.server.Shutdown(context.Background()); err != nil {
		api.Logger().Errorf("error shutting down API Component http server: %s", err.Error())
	}
}

// AddControllers adds a multiple Rest Controllers to the controllers list.
func (api *APIComponent) AddControllers(controllers ...sgul.RestController) *APIComponent {
	api.controllers = append(api.controllers, controllers...)
	return api
}

// AddController adds a single Rest Controller to the controllers list.
func (api *APIComponent) AddController(controller sgul.RestController) *APIComponent {
	return api.AddControllers(controller)
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

// NewAPIComponentWith returns a new API component instance initialized with the
// controllers list.
func NewAPIComponentWith(controllers ...sgul.RestController) *APIComponent {
	return NewAPIComponent().AddControllers(controllers...)
}

// NewDefaultAPIComponent returns a new API component instance configured
// with default middlewares.
func NewDefaultAPIComponent() *APIComponent {
	api := NewAPIComponent()
	api.Use(
		middleware.RequestID,
		middleware.RealIP,
		chilogger.NewZapMiddleware("router", sgul.GetLogger().Desugar()),
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.DefaultCompress,
	)
	return api
}

// NewDefaultAPIComponentWith returns a new API component instance configured
// with default middlewares and initialized with the controllers list.
func NewDefaultAPIComponentWith(controllers ...sgul.RestController) *APIComponent {
	return NewAPIComponent().AddControllers(controllers...)
}
