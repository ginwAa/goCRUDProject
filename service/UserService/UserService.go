package UserService

import "project1/entity"
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

func Count(user entity.User) (int64, error) {
	return UserDao.Sum()
}

func Update(user entity.User) error {
	user.MarkUpdated()
	user.SaltMD5HashPassword()
	return UserDao.Update(user)
}
