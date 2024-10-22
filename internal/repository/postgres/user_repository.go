package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zulfikarmuzakir/e_procurement/internal/domain"

	postgres "github.com/zulfikarmuzakir/e_procurement/internal/repository/postgres/sqlc"
)

type userRepository struct {
	q *postgres.Queries
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{q: postgres.New(db)}
}

// Create implements domain.UserRepository.
func (u *userRepository) Create(user *domain.User) error {
	ctx := context.Background()
	_, err := u.q.CreateUser(ctx, postgres.CreateUserParams{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Status:   user.Status,
	})

	return err
}

// GetByID implements domain.UserRepository.
func (u *userRepository) GetByID(id int64) (*domain.User, error) {
	ctx := context.Background()
	dbUser, err := u.q.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        int64(dbUser.ID),
		Name:      dbUser.Name,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Role:      dbUser.Role,
		Status:    dbUser.Status,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetByEmail implements domain.UserRepository.
func (u *userRepository) GetByEmail(email string) (*domain.User, error) {
	ctx := context.Background()
	dbUser, err := u.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        int64(dbUser.ID),
		Name:      dbUser.Name,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Role:      dbUser.Role,
		Status:    dbUser.Status,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetAllByRole implements domain.UserRepository.
func (u *userRepository) GetAllByRole(role string) ([]*domain.User, error) {
	ctx := context.Background()
	dbUsers, err := u.q.GetAllByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = &domain.User{
			ID:        int64(dbUser.ID),
			Name:      dbUser.Name,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			Role:      dbUser.Role,
			Status:    dbUser.Status,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

func (u *userRepository) Update(user *domain.User) error {
	ctx := context.Background()
	return u.q.UpdateUser(ctx, postgres.UpdateUserParams{
		ID:       int32(user.ID),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Status:   user.Status,
	})
}

// Delete implements domain.UserRepository.
func (u *userRepository) Delete(id int64) error {
	ctx := context.Background()
	return u.q.DeleteUser(ctx, int32(id))
}
