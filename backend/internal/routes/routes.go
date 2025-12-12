package routes

import (
	"net/http"

	"spectra-backend/internal/handlers"
	"spectra-backend/internal/middleware"
	"spectra-backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Router wraps http.ServeMux with additional features
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Handle registers a handler for a pattern
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

// HandleFunc registers a handler function for a pattern
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// SetupRoutes configures all API routes
func SetupRoutes(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	// Initialize repositories
	nodeRepo := repository.NewNodeRepository(pool)
	cableRepo := repository.NewCableRepository(pool)
	customerRepo := repository.NewCustomerRepository(pool)
	connectionRepo := repository.NewConnectionRepository(pool)

	// Initialize handlers
	nodeHandler := handlers.NewNodeHandler(nodeRepo)
	cableHandler := handlers.NewCableHandler(cableRepo)
	customerHandler := handlers.NewCustomerHandler(customerRepo)
	connectionHandler := handlers.NewConnectionHandler(connectionRepo)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","message":"SPECTRA API is running"}`))
	})

	// Node routes
	mux.HandleFunc("GET /api/nodes", nodeHandler.List)
	mux.HandleFunc("POST /api/nodes", nodeHandler.Create)
	mux.HandleFunc("GET /api/nodes/nearby", nodeHandler.GetNearby)
	mux.HandleFunc("GET /api/nodes/{id}", nodeHandler.GetByID)
	mux.HandleFunc("PUT /api/nodes/{id}", nodeHandler.Update)
	mux.HandleFunc("DELETE /api/nodes/{id}", nodeHandler.Delete)

	// Cable routes
	mux.HandleFunc("GET /api/cables", cableHandler.List)
	mux.HandleFunc("POST /api/cables", cableHandler.Create)
	mux.HandleFunc("GET /api/cables/{id}", cableHandler.GetByID)
	mux.HandleFunc("PUT /api/cables/{id}", cableHandler.Update)
	mux.HandleFunc("DELETE /api/cables/{id}", cableHandler.Delete)
	mux.HandleFunc("GET /api/cables/{id}/cores", cableHandler.GetCores)
	mux.HandleFunc("PUT /api/cables/{id}/cores/{coreId}", cableHandler.UpdateCore)

	// Connection routes
	mux.HandleFunc("GET /api/connections", connectionHandler.List)
	mux.HandleFunc("POST /api/connections", connectionHandler.Create)
	mux.HandleFunc("GET /api/connections/{id}", connectionHandler.GetByID)
	mux.HandleFunc("DELETE /api/connections/{id}", connectionHandler.Delete)
	mux.HandleFunc("GET /api/connections/matrix/{nodeId}", connectionHandler.GetSpliceMatrix)
	mux.HandleFunc("GET /api/connections/location/{nodeId}", connectionHandler.GetByLocation)

	// Customer routes
	mux.HandleFunc("GET /api/customers", customerHandler.List)
	mux.HandleFunc("POST /api/customers", customerHandler.Create)
	mux.HandleFunc("GET /api/customers/los", customerHandler.GetLOS)
	mux.HandleFunc("GET /api/customers/{id}", customerHandler.GetByID)
	mux.HandleFunc("PUT /api/customers/{id}", customerHandler.Update)
	mux.HandleFunc("DELETE /api/customers/{id}", customerHandler.Delete)
	mux.HandleFunc("PATCH /api/customers/{id}/status", customerHandler.UpdateStatus)

	// GeoJSON routes
	mux.HandleFunc("GET /api/geojson/nodes", nodeHandler.GetGeoJSON)
	mux.HandleFunc("GET /api/geojson/cables", cableHandler.GetGeoJSON)

	// Apply middleware
	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.Logger,
		middleware.CORS,
		middleware.ContentType,
	)

	return handler
}

// Note: Go 1.22+ ServeMux supports path parameters with {param} syntax
// Use r.PathValue("id") in handlers to extract path parameters
