package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	 "shop_go/user/user"
)

// Product представляє модель продукту
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Count int    `json:"count"`
}

// ProductService надає методи для роботи з продуктами
type ProductService struct {
	DB *sql.DB
}

// GetProducts повертає список усіх продуктів з пагінацією
//GET /api/products?page=1&limit=10

func (s *ProductService) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Отримання значень параметрів пагінації
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // За замовчуванням 10 елементів на сторінці
	}

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit

	// Вибірка продуктів з бази даних з пагінацією
	rows, err := s.DB.Query("SELECT * FROM products LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var products []Product

	// Зчитування результатів запиту
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Count); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	// Перевірка наявності помилок під час зчитування
	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetProduct повертає інформацію про конкретний продукт за ID
// GetProduct повертає інформацію про конкретний продукт за ID
func (s *ProductService) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	productID := chi.URLParam(r, "id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Вибірка конкретного продукту з бази даних за ID
	row := s.DB.QueryRow("SELECT * FROM products WHERE id=?", productID)

	// Створення змінної для зберігання результатів
	var product Product

	// Зчитування результатів запиту
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Count)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Відправлення відповіді у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодуємо та виводимо дані у відповідь
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


// CreateProduct додає новий продукт
func (s *ProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Отримання даних про новий продукт з тіла запиту (JSON)
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(newProduct)

	// Логіка додавання нового продукту до бази даних
	result, err := s.DB.Exec("INSERT INTO products (name, price, count) VALUES (?, ?, ?)",
		newProduct.Name, newProduct.Price, newProduct.Count)
	if err != nil {
		log.Println("Error inserting product into database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отримання ID новоствореного продукту
	newProductID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newProduct.ID = int(newProductID)

	// Відправлення відповіді у форматі JSON з новоствореним продуктом та статусом 201 (Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newProduct); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


// UpdateProduct оновлює інформацію про продукт за ID
func (s *ProductService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	productID := chi.URLParam(r, "id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання нових даних про продукт з тіла запиту (JSON)
	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Логіка оновлення інформації про продукт в базі даних за ID
	result, err := s.DB.Exec("UPDATE products SET name=?, price=?, count=? WHERE id=?",
		updatedProduct.Name, updatedProduct.Price, updatedProduct.Count, productID)
	if err != nil {
		log.Println("Error updating product in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Перевірка, чи існує продукт за вказаним ID
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		// Якщо немає відповідного продукту, відправити HTTP статус 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Відправлення відповіді у форматі JSON з оновленим продуктом
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


// DeleteProduct видаляє продукт за ID
func (s *ProductService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Отримання ID продукту з URL-параметра
	productID := chi.URLParam(r, "id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Логіка видалення продукту з бази даних за ID
	result, err := s.DB.Exec("DELETE FROM products WHERE id=?", productID)
	if err != nil {
		log.Println("Error deleting product from database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Перевірка, чи існує продукт за вказаним ID
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		// Якщо немає відповідного продукту, відправити HTTP статус 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Відправлення відповіді з підтвердженням видалення та статусом 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}


func main() {
	// Ініціалізація роутера
	r := chi.NewRouter()

	// Ініціалізація сервісу продуктів з підключенням до бази даних
	db, err := sql.Open("mysql", "root:usbw@tcp(localhost:3306)/dbshopgo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productService := &ProductService{DB: db}
	userSvc := &user.UserService{DB: db} 

	// Додавання middleware для логування запитів
	r.Use(middleware.Logger)

	// Додавання роутів
	r.Get("/products", productService.GetProducts)
	r.Get("/products/{id}", productService.GetProduct)
	r.Post("/products", productService.CreateProduct)
	r.Put("/products/{id}", productService.UpdateProduct)
	r.Delete("/products/{id}", productService.DeleteProduct)
 r.Route("/users", func(r chi.Router) {
        r.Get("/", userSvc.GetUsers)
        r.Get("/{id}", userSvc.GetUser)
        r.Post("/", userSvc.CreateUser)
        r.Put("/{id}", userSvc.UpdateUser)
        r.Delete("/{id}", userSvc.DeleteUser)
    })
	
	// Запуск сервера на порту 8080
	port := getPort()
	fmt.Printf("Server is running on :%s...\n", port)
	http.ListenAndServe(":"+port, r)
}

// getPort повертає номер порту для веб-сервера
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000" // За замовчуванням використовуємо 8080
	}
	return port
}
