package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gustavoeguedes/api-fullcycle/internal/dto"
	"github.com/gustavoeguedes/api-fullcycle/internal/entity"
	"github.com/gustavoeguedes/api-fullcycle/internal/infra/database"
	entityPkg "github.com/gustavoeguedes/api-fullcycle/pkg/entity"
)

type ProductHandler struct {
	ProductDb database.ProductInterface
}

func NewProductHandler(productDb database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDb: productDb}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = h.ProductDb.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.ProductDb.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)

}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productFound, err := h.ProductDb.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if product.Name == "" {
		product.Name = productFound.Name
	}

	if product.Price == 0 {
		product.Price = productFound.Price
	}

	err = h.ProductDb.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.ProductDb.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.ProductDb.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDb.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

}
