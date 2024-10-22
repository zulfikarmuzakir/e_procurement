package domain

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAllByRole(role string) ([]*User, error)
	Update(user *User) error
	Delete(id int64) error
}

type UserUsecase interface {
	Register(user *User) error
	Login(email, password string) (string, string, error)
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAllByRole(role string) ([]*User, error)
	Update(user *User) error
	Delete(id int64) error
	RefreshToken(refreshToken string) (string, error)
	ApproveVendor(id int64) error
	RejectVendor(id int64) error
}
