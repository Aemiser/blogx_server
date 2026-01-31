package user_service

import "blogx_server/models"

type UserService struct {
	models.UserModel
}

func NewUserService(userModel *models.UserModel) *UserService {
	return &UserService{
		UserModel: *userModel,
	}
}
