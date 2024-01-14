package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	s "strconv"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

// User представляє структуру користувача
type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
}

// UserService надає методи для роботи з користувачами
type UserService struct {
	DB *sql.DB
}

// GetUsers повертає список усіх користувачів з пагінацією
func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Отримання значень параметрів пагінації
	page := getQueryParamInt(r, "page", 1)
	limit := getQueryParamInt(r, "limit", 10)

	// Розрахунок зсуву (offset) для пагінації
	offset := (page - 1) * limit

	// Вибірка користувачів з бази даних з пагінацією
	rows, err := s.DB.Query("SELECT * FROM users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Error querying database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Створення слайсу для зберігання результатів
	var users []User

	// Зчитування результатів запиту
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.ModifiedAt); err != nil {
			log.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
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
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetUser повертає інформацію про конкретного користувача за ID
func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Вибірка конкретного користувача з бази даних за ID (з параметром)
	row := s.DB.QueryRow("SELECT * FROM users WHERE id=?", userID)

	// Створення змінної для зберігання результатів
	var user User

	// Зчитування результатів запиту
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.ModifiedAt)
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
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Create
func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Додавання нового користувача до бази даних
	result, err := s.DB.Exec("INSERT INTO users (first_name, last_name, email, password, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?)",
		newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting user into database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отримання ID новоствореного користувача
	userID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отримання повнішої інформації про новоствореного користувача
	newUser.ID = int(userID)
	err = s.DB.QueryRow("SELECT * FROM users WHERE id=?", userID).Scan(&newUser.ID, &newUser.FirstName, &newUser.LastName, &newUser.Email, &newUser.Password, &newUser.CreatedAt, &newUser.ModifiedAt)
	if err != nil {
		log.Println("Error querying new user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON з повною інформацією про нового користувача
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateUser
func (s *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отримання старої інформації про користувача
	var oldUser User
	err := s.DB.QueryRow("SELECT * FROM users WHERE id=?", userID).Scan(&oldUser.ID, &oldUser.FirstName, &oldUser.LastName, &oldUser.Email, &oldUser.Password, &oldUser.CreatedAt, &oldUser.ModifiedAt)
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
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Оновлення інформації про користувача в базі даних
	_, err = s.DB.Exec("UPDATE users SET first_name=?, last_name=?, email=?, password=?, modified_at=? WHERE id=?",
		updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.Password, time.Now(), userID)
	if err != nil {
		log.Println("Error updating user in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді у форматі JSON з оновленою інформацією про користувача
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Println("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteUser видаляє користувача за ID
func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Отримання ID користувача з URL-параметра
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Видалення користувача з бази даних за ID
	result, err := s.DB.Exec("DELETE FROM users WHERE id=?", userID)
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

// Введення допоміжної функції для отримання значення параметра запиту як ціле число
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

