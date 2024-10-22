package converter

import "github.com/zulfikarmuzakir/e_procurement/internal/domain"

func UserToUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
