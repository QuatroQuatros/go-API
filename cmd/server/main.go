package main

import (
	"fmt"
	"net/http"

	"github.com/QuatroQuatros/go-API/configs"
	_ "github.com/QuatroQuatros/go-API/docs"
	"github.com/QuatroQuatros/go-API/internal/entity"
	"github.com/QuatroQuatros/go-API/internal/infra/database"
	"github.com/QuatroQuatros/go-API/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           GO Example API
// @version         1.0
// @description     Simple API with authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Guilherme Cesar
// @contact.url    http://www.swagger.io/support
// @contact.email  guicossalter@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:8000
// @BasePath  /
// @securityDefinitions.apiKey  ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)

	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))
	//r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.GetAllProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	fmt.Println("Servidor rodando no endereço http://127.0.0.1:8000")

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(" http://127.0.0.1:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

// func LogRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Request: %s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }
