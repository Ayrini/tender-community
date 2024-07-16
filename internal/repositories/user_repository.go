package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"tender/internal/models"
)

type UserRepository struct {
	Db *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("Пользователь с таким адресом электронной почты уже существует")
	ErrDuplicatePhone = errors.New("Пользователь с таким номером телефона уже существует")
	ErrNotFound       = func(errorMessage string) error {
		return errors.New(fmt.Sprintf("no client found with the given %s", errorMessage))
	}
)

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{Db: db}
}


func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {

	rows, err := r.Db.QueryContext(ctx, "SELECT id, name, last_name, email, phone, inn, balance, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Phone, &user.INN, &user.Balance, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SignUp(ctx context.Context, user models.User) error {

	var exists int
	emailCheckQuery := "SELECT COUNT(*) FROM users WHERE email= ?"
	phoneCheckQuery := "SELECT COUNT(*) FROM users WHERE phone IS NOT NULL AND phone = ? "

	err := r.Db.QueryRow(emailCheckQuery, user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 && user.Email != "" {
		return ErrDuplicateEmail
	}

	err = r.Db.QueryRow(phoneCheckQuery, user.Phone).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 && user.Phone != "" {
		return ErrDuplicatePhone
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	_, err = r.Db.ExecContext(ctx, "INSERT INTO users(name, last_name, email, phone, inn, password, balance) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Name, user.LastName, user.Email, user.Phone, user.INN, hashedPassword, 0)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return models.ErrDuplicateEmail
		}
		return err
	}

	//clientID, err := result.LastInsertId()
	//if err != nil {
	//	return err
	//}

	//convertedClientInfo, err := json.Marshal(clientID)
	//if err != nil {
	//	return err
	//}

	return nil
}
