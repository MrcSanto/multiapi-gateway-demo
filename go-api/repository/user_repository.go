package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{
		connection: connection,
	}
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {

	query, err := ur.connection.Prepare("SELECT * FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}

	var user model.User
	err = query.QueryRow(email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println(err)
		return nil, err
	}
	query.Close()

	return &user, nil
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {
	var id int
	query, err := ur.connection.Prepare("INSERT INTO users" +
		"(name, email, password_hash)" +
		" VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return id, nil
}
