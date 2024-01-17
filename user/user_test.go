package user_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// Імпорт вашого пакету user та інших необхідних залежностей
	"github.com/chitawebui131/shop_go/user"
)

// DBQueryer описує метод Query бази даних
type DBQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// MockDB - мок бази даних, що реалізує інтерфейс DBQueryer
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Мокання методу Query
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

func TestGetUsers(t *testing.T) {
	// Створення інстанції мок-бази даних
	mockDB := new(MockDB)

	// Створення інстанції UserService з мок-базою даних
	userService := &user.UserService{DB: mockDB}

	// Параметри тестового запиту
	req, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	// Використання httptest для створення запису відповіді
	rr := httptest.NewRecorder()

	// Мокання методу Query бази даних
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "created_at", "modified_at"}).
		AddRow(1, "John", "Doe", "john@example.com", "password", time.Now(), time.Now())
	mockDB.On("Query", mock.Anything, mock.Anything).Return(rows, nil)

	// Виклик функції обробки HTTP-запиту
	userService.GetUsers(rr, req)

	// Перевірка статус-коду відповіді
	assert.Equal(t, http.StatusOK, rr.Code)

	// Розкодування JSON та перевірка результатів
	var users []user.User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "John", users[0].FirstName)
	// Додайте інші перевірки за необхідності

	// Перевірка викликів моканого методу Query
	mockDB.AssertExpectations(t)
}
