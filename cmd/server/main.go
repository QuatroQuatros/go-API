package main

import (
	"fmt"
	"net/http"

	"github.com/QuatroQuatros/go-API/configs"
	"github.com/QuatroQuatros/go-API/internal/entity"
	"github.com/QuatroQuatros/go-API/internal/infra/database"
	"github.com/QuatroQuatros/go-API/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)

	fmt.Println("Servidor rodando no endereço http://127.0.0.1:8000")
	http.ListenAndServe(":8000", r)
}
