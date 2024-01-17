package categories

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	s "strconv"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

// User представляє структуру користувача
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CatSetvices struct {
	DB *sql.DB
}

func getQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	result, err := s.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}

func (s *CatSetvices) GetCats(w http.ResponseWriter, r *http.Request) {
	// Отримання значень параметрів пагінації
	page := getQueryParamInt(r, "page", 1)
	limit := getQueryParamInt(r, "limit", 10)

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit

	// Вибірка користувачів з бази даних з пагінацією
	rows, err := s.DB.Query("SELECT * FROM categories LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var cats []Category

	// Зчитування результатів запиту
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cats = append(cats, cat)
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
	if err := json.NewEncoder(w).Encode(cats); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Get category
func (s *CatSetvices) GetCat(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	catID := chi.URLParam(r, "id")
	if catID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Вибірка конкретного користувача з бази даних за ID (з параметром)
	row := s.DB.QueryRow("SELECT * FROM categories WHERE id=?", catID)

	// Створення змінної для зберігання результатів
	var cat Category

	// Зчитування результатів запиту
	err := row.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
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
	if err := json.NewEncoder(w).Encode(cat); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Create
func (s *CatSetvices) CreateCat(w http.ResponseWriter, r *http.Request) {
	var newCat Category
	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Додавання нового користувача до бази даних
	result, err := s.DB.Exec("INSERT INTO categories (name, description, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?)",
		newCat.Name, newCat.Description, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting user into database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отримання ID новоствореного користувача
	catID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отримання повнішої інформації про новоствореного користувача
	newCat.ID = int(catID)
	err = s.DB.QueryRow("SELECT * FROM categories WHERE id=?", catID).Scan(&newCat.ID, &newCat.Name, &newCat.Description, &newCat.CreatedAt, &newCat.UpdatedAt)
	if err != nil {
		log.Println("Error querying new user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON з повною інформацією про нового користувача
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newCat); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateUser
func (s *CatSetvices) UpdateCat(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	catID := chi.URLParam(r, "id")
	if catID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання старої інформації про користувача
	var oldCat Category
	err := s.DB.QueryRow("SELECT * FROM categories WHERE id=?", catID).Scan(&oldCat.ID, &oldCat.Name, &oldCat.Description, &oldCat.CreatedAt, &oldCat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println("Error querying user:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Отримання нових даних про користувача з тіла запиту (JSON)
	var updatedCat Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCat); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Оновлення інформації про користувача в базі даних
	_, err = s.DB.Exec("UPDATE categories SET name=?, description=?,  updated_at=? WHERE id=?",
		updatedCat.Name, updatedCat.Description, time.Now(), catID)
	if err != nil {
		log.Println("Error updating user in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON з оновленою інформацією про користувача
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedCat); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteUser видаляє користувача за ID
func (s *CatSetvices) DeleteCat(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	catID := chi.URLParam(r, "id")
	if catID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Видалення користувача з бази даних за ID
	result, err := s.DB.Exec("DELETE FROM users WHERE id=?", catID)
	if err != nil {
		log.Println("Error deleting user from database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Перевірка, чи існує користувач за вказаним ID
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		// Якщо немає відповідного користувача, відправити HTTP статус 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Відправлення відповіді з підтвердженням видалення та статусом 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
