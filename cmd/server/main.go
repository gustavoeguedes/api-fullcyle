package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gustavoeguedes/api-fullcycle/configs"
	"github.com/gustavoeguedes/api-fullcycle/internal/entity"
	"github.com/gustavoeguedes/api-fullcycle/internal/infra/database"
	"github.com/gustavoeguedes/api-fullcycle/internal/infra/webserver/handlers"
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
	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProduct(db)

	productHandler := handlers.NewProductHandler(productDB)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/products", productHandler.GetProducts)
	r.Post("/products", productHandler.Create)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.Update)
	r.Delete("/products/{id}", productHandler.Delete)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
