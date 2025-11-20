package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"ool/api/schema"
	"ool/internal/app/coupon"
	"ool/internal/app/product"
)

// Handler holds the application dependencies required by the HTTP handlers.
type Handler struct {
	ProductRepo *map[string]schema.Product
	CouponIndex *coupon.Index
	CouponCache *coupon.LRUCache
}

// NewHandler creates a new Handler instance with dependencies.
func NewHandler(
	repo *map[string]schema.Product,
	ci *coupon.Index,
	cc *coupon.LRUCache,
) *Handler {
	return &Handler{
		ProductRepo: repo,
		CouponIndex: ci,
		CouponCache: cc,
	}
}

// ListProducts handles GET /api/product
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	var productList []schema.Product
	for _, p := range *h.ProductRepo {
		productList = append(productList, p)
	}
	respondJSON(w, http.StatusOK, productList)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]
	if _, err := strconv.ParseInt(productID, 10, 64); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID supplied.")
		return
	}

	product, ok := (*h.ProductRepo)[productID]
	if !ok {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req schema.OrderReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	if len(req.Items) == 0 {
		respondError(w, http.StatusUnprocessableEntity, "Order must contain at least one item.")
		return
	}

	if req.CouponCode != "" {
		if !coupon.IsCouponValid(h.CouponIndex, h.CouponCache, req.CouponCode) {
			respondError(w, http.StatusUnprocessableEntity, "Invalid coupon code.")
			return
		}
	}

	var orderItems []schema.OrderItem
	var orderedProducts []schema.Product
	totalPrice := 0.0

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			respondError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Qantity for product %s must be greater than zero.", item.ProductID))
			return
		}

		product, ok := (*h.ProductRepo)[item.ProductID]
		if !ok {
			respondError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Product with ID '%s' not found.", item.ProductID))
			return
		}

		totalPrice += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, schema.OrderItem{ProductID: item.ProductID, Quantity: item.Quantity})
		orderedProducts = append(orderedProducts, product)
	}

	newOrder := schema.Order{
		ID:       uuid.New().String(),
		Items:    orderItems,
		Products: orderedProducts,
	}

	respondJSON(w, http.StatusOK, newOrder)
}

func GetProductRepo() *map[string]schema.Product {
	return &product.Repository
}
