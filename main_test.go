package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	r := chi.NewRouter()
	productService := &mockProductService{} // Вам слід замінити це на ваш власний мок ProductService
	userService := &mockUserService{}       // Вам слід замінити це на ваш власний мок UserService
	catService := &mockCatService{}         // Вам слід замінити це на ваш власний мок CatService

	r.Use(middleware.Logger)

	r.Get("/products", productService.GetProducts)
	r.Get("/products/{id}", productService.GetProduct)
	r.Post("/products", productService.CreateProduct)
	r.Put("/products/{id}", productService.UpdateProduct)
	r.Delete("/products/{id}", productService.DeleteProduct)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userService.GetUsers)
		r.Get("/{id}", userService.GetUser)
		r.Post("/", userService.CreateUser)
		r.Put("/{id}", userService.UpdateUser)
		r.Delete("/{id}", userService.DeleteUser)
	})

	r.Route("/cat", func(r chi.Router) {
		r.Get("/", catService.GetCats)
		r.Get("/{id}", catService.GetCat)
		r.Post("/", catService.CreateCat)
		r.Put("/{id}", catService.UpdateCat)
		r.Delete("/{id}", catService.DeleteCat)
	})

	// Тестуємо роут /products
	t.Run("TestGetProducts", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Тут ви можете провести інші перевірки для відповіді, якщо потрібно
	})

	// Тестуємо роут /products/{id}
	t.Run("TestGetProduct", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products/123", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Тут ви можете провести інші перевірки для відповіді, якщо потрібно
	})

	// Тестуємо роут /products (POST)
	t.Run("TestCreateProduct", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/products", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Тут ви можете провести інші перевірки для відповіді, якщо потрібно
	})

	// Тут можна додати інші тести для інших роутів
}

// Ось приклад моків для ProductService, UserService, CatService (вам слід реалізувати свої моки)
type mockProductService struct{}
type mockUserService struct{}
type mockCatService struct{}

func (m *mockProductService) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Мок для GetProducts
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "GetProducts mock response"}`))
}

func (m *mockProductService) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Мок для GetProduct
}

func (m *mockProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Мок для CreateProduct
}

func (m *mockProductService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Мок для UpdateProduct
}

func (m *mockProductService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Мок для DeleteProduct
}

// Аналогічно реалізуйте моки для UserService і CatService

func (m *mockUserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Мок для GetUsers
}

func (m *mockUserService) GetUser(w http.ResponseWriter, r *http.Request) {
	// Мок для GetUser
}

func (m *mockUserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Мок для CreateUser
}

func (m *mockUserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Мок для UpdateUser
}

func (m *mockUserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Мок для DeleteUser
}

func (m *mockCatService) GetCats(w http.ResponseWriter, r *http.Request) {
	// Мок для GetCats
}

func (m *mockCatService) GetCat(w http.ResponseWriter, r *http.Request) {
	// Мок для GetCat
}

func (m *mockCatService) CreateCat(w http.ResponseWriter, r *http.Request) {
	// Мок для CreateCat
}

func (m *mockCatService) UpdateCat(w http.ResponseWriter, r *http.Request) {
	// Мок для UpdateCat
}

func (m *mockCatService) DeleteCat(w http.ResponseWriter, r *http.Request) {
	// Мок для DeleteCat
}
