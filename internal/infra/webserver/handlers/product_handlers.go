package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/QuatroQuatros/go-API/internal/dto"
	"github.com/QuatroQuatros/go-API/internal/entity"
	"github.com/QuatroQuatros/go-API/internal/infra/database"
	"github.com/go-chi/chi/v5"

	entityPkg "github.com/QuatroQuatros/go-API/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary			Create Product
// @Description		Create Product
// @Tags 			products
// @Accept			json
// @Produce			json
// @Param			request		body		dto.CreateProductInput 	true 	"product request"
// @Success			201
// @Failure			404			{object}	Error
// @Failure			500			{object}	Error
// @Router			/products [post]
// @Security		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// List Products godoc
// @Summary			List Products
// @Description		Get All Products
// @Tags 			products
// @Accept			json
// @Produce			json
// @Param			page		query		string 	false 	"page number"
// @Param			limit		query		string 	false 	"limit"
// @Success			200			{array}     entity.Product
// @Failure			404			{object}	Error
// @Failure			500			{object}	Error
// @Router			/products [get]
// @Security		ApiKeyAuth
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Get Product godoc
// @Summary			Get a product
// @Description		Get a product
// @Tags 			products
// @Accept			json
// @Produce			json
// @Param			id			path		string 	true 	"product ID" Format(uuid)
// @Success			200			{object}    entity.Product
// @Failure			404         {object}	Error
// @Failure			500			{object}	Error
// @Router			/products/{id} [get]
// @Security		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update Product godoc
// @Summary			Update a product
// @Description		Update a product
// @Tags 			products
// @Accept			json
// @Produce			json
// @Param			id			path		string 					true 	"product ID" Format(uuid)
// @Param			request		body		dto.CreateProductInput 	true 	"product request"
// @Success			200			{object}    entity.Product
// @Failure			404         {object}	Error
// @Failure			500			{object}	Error
// @Router			/products/{id} [put]
// @Security		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}

// Delete Product godoc
// @Summary			Delete a product
// @Description		Delete a product
// @Tags 			products
// @Accept			json
// @Produce			json
// @Param			id			path		string 	true 	"product ID" Format(uuid)
// @Success			200
// @Failure			404         {object}	Error
// @Failure			500			{object}	Error
// @Router			/products/{id} [delete]
// @Security		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(product)
}
