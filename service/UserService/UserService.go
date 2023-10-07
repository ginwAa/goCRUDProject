package UserService

import (
	"project1/entity"
	"project1/entity/DTO"
)
import "project1/dao/UserDao"

func Page(user entity.User, page int, count int) ([]entity.User, error) {
	if user.Username != "" {
		return UserDao.PageByUsername(user.Username, page, count)
	} else {
		return UserDao.Page(page, count)
	}
}

func Add(user entity.User) (int64, error) {
	user.SaltMD5HashPassword()
	user.MarkCreated()
	return UserDao.Add(user)
}

func Count(userCountDTO DTO.UserCountDTO) (int64, error) {
	var user entity.User
	user.Username = userCountDTO.Username
	user.Account = userCountDTO.Account
	if userCountDTO.Like != "true" {
		return UserDao.CountByUser(user)
	} else {
		return UserDao.CountByUserLike(user)
	}
}

func Update(user entity.User) error {
	if !user.DeletedAt.IsZero() {
		return UserDao.Delete(user)
	}
	user.MarkUpdated()
	if user.Password != "" {
		user.SaltMD5HashPassword()
	}
	return UserDao.Update(user)
}
