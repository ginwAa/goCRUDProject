package UserService

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"project1/constants"
	"project1/entity"
	"project1/entity/DTO"
	"project1/entity/VO"
	"time"
)
import "project1/dao/UserDao"

func Page(user entity.User, page int, size int) ([]entity.User, error) {
	if user.Username != "" {
		return UserDao.PageByUsername(user.Username, page, size)
	} else {
		return UserDao.Page(page, size)
	}
}

func Add(user entity.User) (int64, error) {
	if user.Account == "" || user.Password == "" || user.Username == "" ||
		user.Role != 1 && user.Role != 2 || user.Status != 1 && user.Status != 2 || user.Gender != 1 && user.Gender != 2 {
		return 0, fmt.Errorf(constants.USER_DATA_ERROR)
	}
	user.Password = entity.SaltMD5Hash(user.Password)
	user.MarkCreated()
	id, err := UserDao.Add(user)
	if err != nil {
		log.Printf("user/add err2: %s\n", err)
		return 0, err
	}
	return id, nil
}

func Count(username string) (int64, error) {
	if username == "" {
		return UserDao.CountAll()
	} else {
		return UserDao.CountByUsernameLike(username)
	}
}

func Update(user entity.User) error {
	user.MarkUpdated()
	if user.Password != "" {
		user.Password = entity.SaltMD5Hash(user.Password)
		res, err := UserDao.UpdatePassword(user)
		if err != nil {
			log.Printf("user/update err1: %s\n", err)
			return err
		}
		if res != 1 {
			log.Printf("user/update err2: empty update\n")
			return fmt.Errorf("")
		}
		return nil
	}
	if user.Username != "" {
		res, err := UserDao.UpdateUsername(user)
		if err != nil {
			log.Printf("user/update err3: %s\n", err)
			return err
		}
		if res != 1 {
			log.Printf("user/update err4: empty update\n")
			return fmt.Errorf("")
		}
		return nil
	}
	return fmt.Errorf("no handler, bad request!")
}

func Login(userQ DTO.UserLoginDTO) (VO.LoginVO, error) {
	if userQ.Account == "" || userQ.Password == "" {
		return VO.LoginVO{}, fmt.Errorf(constants.USER_DATA_ERROR)
	}
	userQ.Password = entity.SaltMD5Hash(userQ.Password)
	user, err := UserDao.SelectOneByLogin(userQ.Account, userQ.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return VO.LoginVO{
			Success: false,
		}, nil
	}
	if err != nil || user.Role == 0 {
		log.Printf("/user/login err1: %s\n", err)
		return VO.LoginVO{
			Success: false,
		}, err
	}
	roleStr := "normal"
	if user.Role == 2 {
		roleStr = "admin"
	}
	expirationTime := time.Now().Add(time.Hour)
	claims := &jwt.MapClaims{
		"role": roleStr,
		"exp":  expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, err := token.SignedString(entity.SecretKey)
	if err != nil {
		log.Printf("/user/login err2: %s\n", err)
		return VO.LoginVO{
			Success: false,
		}, err
	}
	loginVO := VO.LoginVO{
		Token:   tokenStr,
		Role:    roleStr,
		Success: true,
	}
	return loginVO, nil
}

func Delete(userId int64) error {
	if userId == 0 {
		return fmt.Errorf(constants.USER_DATA_ERROR)
	}
	user := entity.User{}
	user.Id = userId
	user.MarkDeleted()
	res, err := UserDao.Delete(user)
	if err != nil {
		log.Printf("/user/delete err: %s\n", err)
		return err
	}
	if res != 1 {
		return fmt.Errorf("")
	} else {
		return nil
	}
}

func CheckAccountUnique(account string) (bool, error) {
	if account == "" {
		return false, fmt.Errorf(constants.USER_DATA_ERROR)
	}
	res, err := UserDao.CountByAccount(account)
	if err != nil {
		fmt.Errorf("user/checkAccountUnique err: %s\n", err)
		return false, err
	} else {
		return res == 0, nil
	}
}
