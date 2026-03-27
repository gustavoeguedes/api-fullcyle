package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gustavoeguedes/api-fullcycle/configs"
	_ "github.com/gustavoeguedes/api-fullcycle/docs"
	"github.com/gustavoeguedes/api-fullcycle/internal/entity"
	"github.com/gustavoeguedes/api-fullcycle/internal/infra/database"
	"github.com/gustavoeguedes/api-fullcycle/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title API Full Cycle
// @version 1.0
// @description API Full Cycle - Go Expert

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email email@email.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)

	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.GetProducts)
		r.Post("/", productHandler.Create)
		r.Get("//{id}", productHandler.GetProduct)
		r.Put("//{id}", productHandler.Update)
		r.Delete("//{id}", productHandler.Delete)
	})

	r.Post("/users", userHandler.Create)
	r.Post("/users/login", userHandler.GetJWT)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
