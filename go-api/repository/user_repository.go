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

func (ur *UserRepository) GetUsers() ([]model.User, error) {
	query, err := ur.connection.Prepare("SELECT id, name, email, username, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.User,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetUserById(idUser int) (*model.User, error) {

	query, err := ur.connection.Prepare("SELECT id, name, email, username, created_at, updated_at FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var user model.User
	err = query.QueryRow(idUser).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.User,
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

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {

	query, err := ur.connection.Prepare("SELECT id, name, email, username, password, created_at, updated_at FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var user model.User
	err = query.QueryRow(email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.User,
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

	return &user, nil
}

func (ur *UserRepository) GetUserByUserName(username string) (*model.User, error) {

	query, err := ur.connection.Prepare("SELECT id, name, email, username, password, created_at, updated_at FROM users WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var user model.User
	err = query.QueryRow(username).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.User,
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

	return &user, nil
}

func (ur *UserRepository) DeleteUser(idUser int) (int64, error) {
	query, err := ur.connection.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return 0, err
	}
	defer query.Close()

	result, err := query.Exec(idUser)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (ur *UserRepository) UpdateUser(idUser int, user model.User) error {
	query, err := ur.connection.Prepare(
		"UPDATE users SET name = $1, email = $2, username = $3, updated_at = NOW() WHERE id = $4",
	)
	if err != nil {
		return err
	}
	defer query.Close()

	result, err := query.Exec(user.Name, user.Email, user.User, idUser)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario nao encontrado")
	}

	return nil
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {
	var id int
	query, err := ur.connection.Prepare("INSERT INTO users " +
		"(name, email, username, password) " +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer query.Close()

	err = query.QueryRow(user.Name, user.Email, user.User, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}
