package auth

import (
	"database/sql"
	"errors"
)

type Repository struct {
	users []User
	db    *sql.DB
}

func NewRepository() *Repository {
	return &Repository{
		users: []User{
			{
				ID:           1,
				Matricule:    "20250001",
				Name:         "Olivier Mulamba",
				Email:        "olivier@example.com",
				PasswordHash: "$2a$10$example",
				Role:         "STUDENT",
			},
		},
	}
}

// NewRepositoryWithDB crée un repository utilisant une base SQL (SQLite/Postgres via database/sql).
func NewRepositoryWithDB(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByMatricule(matricule string) (*User, error) {
	if r.db != nil {
		var user User
		row := r.db.QueryRow(
			"SELECT id, matricule, name, email, password_hash, role FROM users WHERE matricule = ?",
			matricule,
		)
		if err := row.Scan(&user.ID, &user.Matricule, &user.Name, &user.Email, &user.PasswordHash, &user.Role); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("user not found")
			}
			return nil, err
		}

		return &user, nil
	}

	for _, user := range r.users {
		if user.Matricule == matricule {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	if r.db != nil {
		var user User
		row := r.db.QueryRow(
			"SELECT id, matricule, name, email, password_hash, role FROM users WHERE email = ?",
			email,
		)
		if err := row.Scan(&user.ID, &user.Matricule, &user.Name, &user.Email, &user.PasswordHash, &user.Role); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("user not found")
			}
			return nil, err
		}

		return &user, nil
	}

	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *Repository) Create(user User) (*User, error) {
	if r.db != nil {
		res, err := r.db.Exec(
			"INSERT INTO users (matricule, name, email, password_hash, role) VALUES (?, ?, ?, ?, ?)",
			user.Matricule,
			user.Name,
			user.Email,
			user.PasswordHash,
			user.Role,
		)
		if err != nil {
			return nil, err
		}

		id, err := res.LastInsertId()
		if err == nil {
			user.ID = int(id)
		}

		return &user, nil
	}

	for _, u := range r.users {
		if u.Matricule == user.Matricule || u.Email == user.Email {
			return nil, errors.New("user already exists")
		}
	}

	user.ID = len(r.users) + 1
	r.users = append(r.users, user)

	return &r.users[len(r.users)-1], nil
}
