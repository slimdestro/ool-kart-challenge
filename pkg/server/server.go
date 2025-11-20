package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"ool/api/schema"
	"ool/internal/app/coupon"
	"ool/internal/handler"
)

// NewRouter sets up the Gorilla Mux router with all application routes and middleware.
func NewRouter(ci *coupon.Index, cc *coupon.LRUCache, repo *map[string]schema.Product) *mux.Router {
	r := mux.NewRouter()
	h := handler.NewHandler(repo, ci, cc)
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/product", h.ListProducts).Methods("GET").Name("listProducts")
	apiRouter.HandleFunc("/product/{productId}", h.GetProduct).Methods("GET").Name("getProduct")
	apiRouter.Handle("/order", handler.APIKeyMiddleware(http.HandlerFunc(h.PlaceOrder))).Methods("POST").Name("placeOrder")

	return r
}

// NewServer creates a new configured http.Server instance.
func NewServer(port int, router http.Handler) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	return &http.Server{
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}
}
