package categories

import (
	"database/sql"
	"log"
	"time"

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

func GetCategoriesFromDB(db *sql.DB) ([]Category, error) {
	var categories []Category

	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		log.Println("Error querying categories:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err)
		return nil, err
	}

	return categories, nil
}

// GetCategoryByIDFromDB - отримання категорії з бази даних за ID
func GetCategoryByIDFromDB(db *sql.DB, categoryID int) (*Category, error) {
	var category Category

	row := db.QueryRow("SELECT * FROM categories WHERE id = ?", categoryID)
	if err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}

	return &category, nil
}

// CreateCategoryInDB - створення нової категорії в базі даних
func CreateCategoryInDB(db *sql.DB, newCategory Category) (int, error) {
	result, err := db.Exec("INSERT INTO categories (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)",
		newCategory.Name, newCategory.Description, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting category:", err)
		return 0, err
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted ID:", err)
		return 0, err
	}

	return int(categoryID), nil
}

// UpdateCategoryInDB - оновлення інформації про категорію в базі даних за ID
func UpdateCategoryInDB(db *sql.DB, updatedCategory Category, categoryID int) error {
	_, err := db.Exec("UPDATE categories SET name=?, description=?, updated_at=? WHERE id=?",
		updatedCategory.Name, updatedCategory.Description, time.Now(), categoryID)
	if err != nil {
		log.Println("Error updating category:", err)
		return err
	}

	return nil
}

// DeleteCategoryFromDB - видалення категорії з бази даних за ID
func DeleteCategoryFromDB(db *sql.DB, categoryID int) error {
	_, err := db.Exec("DELETE FROM categories WHERE id=?", categoryID)
	if err != nil {
		log.Println("Error deleting category:", err)
		return err
	}

	return nil
}
