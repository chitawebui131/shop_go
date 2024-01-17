package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/shopspring/decimal"
	"github.com/chitawebui131/shop_go/categories"
	"github.com/chitawebui131/shop_go/user"
)

// Product представляє модель продукту
type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stockQuantity"`
	CategoryID    int       `json:"categoryID"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductWithCategoryWithoutDates struct {
	ProductID           int     `json:"product_id"`
	ProductName         string  `json:"product_name"`
	ProductDescription  string  `json:"product_description"`
	ProductPrice        float64 `json:"product_price"`
	StockQuantity       int     `json:"product_stockQuantity"`
	ProductCategoryID   int     `json:"product_category_id"`
	CategoryID          int     `json:"category_id"`
	CategoryName        string  `json:"category_name"`
	CategoryDescription string  `json:"category_description"`
}

// ProductService надає методи для роботи з продуктами
type ProductService struct {
	DB *sql.DB
}

// GetProducts повертає список усіх продуктів з пагінацією
//GET /api/products?page=1&limit=10
/*
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
	rows, err := s.DB.Query("SELECT * FROM products JOIN categories ON products.category_id = categories.id LIMIT ? OFFSET ?", limit, offset)
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
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CategoryID, &product.Created_at, &product.Updated_at ); err != nil {
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
	query := `
		SELECT products.id AS product_id, products.name AS product_name,
			   products.description AS product_description, products.price AS product_price,
			   products.stock_quantity AS product_stockQuantity, products.category_id AS product_category_id,
			   products.created_at AS created_at, products.updated_at AS updated_at,
			   categories.id AS category_id, categories.name AS category_name,
			   categories.description AS category_description,
			   categories.created_at AS category_created_at, categories.updated_at AS category_updated_at
		FROM products
		JOIN categories ON products.category_id = categories.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var productsWithCategories []ProductWithCategory

	// Зчитування результатів запиту
	for rows.Next() {
		var productWithCategory ProductWithCategory
		// Сканування результатів у структуру продукту з категорією
		if err := rows.Scan(
			&productWithCategory.ProductID,
			&productWithCategory.ProductName,
			&productWithCategory.ProductDescription,
			&productWithCategory.ProductPrice,
			&productWithCategory.StockQuantity,
			&productWithCategory.ProductCategoryID,
			&productWithCategory.CreatedAt,
			&productWithCategory.UpdatedAt,
			&productWithCategory.CategoryID,
			&productWithCategory.CategoryName,
			&productWithCategory.CategoryDescription,
			&productWithCategory.CategoryCreatedAt,
			&productWithCategory.CategoryUpdatedAt,
		); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsWithCategories = append(productsWithCategories, productWithCategory)
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
	if err := json.NewEncoder(w).Encode(productsWithCategories); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
*/

func (s *ProductService) GetProducts(w http.ResponseWriter, r *http.Request) {
	// ... (зберігаємо код пагінації та запиту з бази даних)
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
	query := `
		SELECT products.id AS product_id, products.name AS product_name, 
			   products.description AS product_description, products.price AS product_price, 
			   products.stock_quantity AS product_stockQuantity, products.category_id AS product_category_id,
			   categories.id AS category_id, categories.name AS category_name,
			   categories.description AS category_description
		FROM products
		LEFT JOIN categories ON products.category_id = categories.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var productsWithCategoriesWithoutDates []ProductWithCategoryWithoutDates

	// Зчитування результатів запиту
	for rows.Next() {
		var productWithCategoryWithoutDates ProductWithCategoryWithoutDates
		// Сканування результатів у структуру продукту з категорією без дат
		if err := rows.Scan(
			&productWithCategoryWithoutDates.ProductID,
			&productWithCategoryWithoutDates.ProductName,
			&productWithCategoryWithoutDates.ProductDescription,
			&productWithCategoryWithoutDates.ProductPrice,
			&productWithCategoryWithoutDates.StockQuantity,
			&productWithCategoryWithoutDates.ProductCategoryID,
			&productWithCategoryWithoutDates.CategoryID,
			&productWithCategoryWithoutDates.CategoryName,
			&productWithCategoryWithoutDates.CategoryDescription,
		); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsWithCategoriesWithoutDates = append(productsWithCategoriesWithoutDates, productWithCategoryWithoutDates)
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
	if err := json.NewEncoder(w).Encode(productsWithCategoriesWithoutDates); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
	//fmt.Println(row)
	// Зчитування результатів запиту
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CategoryID, &product.Created_at, &product.Updated_at)
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
		//fmt.Println(&newProduct)
		return
	}
	fmt.Println(newProduct)

	// Логіка додавання нового продукту до бази даних
	// result, err := s.DB.Exec("INSERT INTO products (name, description, price, stock_quantity, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
	// 	newProduct.Name, newProduct.Description, newProduct.Price, newProduct.StockQuantity, newProduct.CategoryID, time.Now(), time.Now())
	// if err != nil {
	// 	log.Println("Error inserting product into database:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	query := `
    INSERT INTO products (name, description, price, stock_quantity, category_id, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?)
`
	result, err := s.DB.Exec(query, newProduct.Name, newProduct.Description, newProduct.Price, newProduct.StockQuantity, newProduct.CategoryID, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting into database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	// result, err := s.DB.Exec("UPDATE products SET name=?, description=?, price=?, stock_quantity=?, category_id=?, updated_at=?   WHERE id=?",
	// 	updatedProduct.Name, updatedProduct.Description, updatedProduct.Price, updatedProduct.StockQuantity, updatedProduct.CategoryID, productID, time.Now())
	query := `
		UPDATE products
		SET
			name = ?,
			description = ?,
			price = ?,
			stock_quantity = ?,
			category_id = ?,
			updated_at = ?
		WHERE id = ?
	`
	result, err := s.DB.Exec(query,
		updatedProduct.Name,
		updatedProduct.Description,
		updatedProduct.Price,
		updatedProduct.StockQuantity,
		updatedProduct.CategoryID,
		time.Now(),
		productID,
	)

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
	db, err := sql.Open("mysql", "root:usbw@tcp(localhost:3306)/dbshopgo?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productService := &ProductService{DB: db}
	userSvc := &user.UserService{DB: db}
	catSvc := &categories.CatSetvices{DB: db}

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
	r.Route("/cat", func(r chi.Router) {
		r.Get("/", catSvc.GetCats)
		r.Get("/{id}", catSvc.GetCat)
		r.Post("/", catSvc.CreateCat)
		r.Put("/{id}", catSvc.UpdateCat)
		r.Delete("/{id}", catSvc.DeleteCat)
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
