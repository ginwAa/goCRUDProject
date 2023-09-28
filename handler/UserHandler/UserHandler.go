package UserHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"project1/entity"
	"project1/entity/Result"
	"project1/error"
	"project1/service/UserService"
	"strconv"
)

// Page 分页查询 GET return user list
func Page(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodGet {
		var user entity.User
		user.Username = r.URL.Query().Get("name")
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		count, _ := strconv.ParseInt(r.URL.Query().Get("count"), 10, 64)

		users, err := UserService.Page(user, int(page), int(count))
		if err != nil {
			log.Println(err)
			res = Result.Error(http.StatusInternalServerError, error.USER_SELECT_ERROR)
		} else {
			res = Result.Success(users)
		}
	}
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Add 添加用户 POST return user id
func Add(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodPost {
		var user entity.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			res = Result.Error(http.StatusBadRequest, error.USER_DATA_ERROR)
			goto ret
		}
		id, err := UserService.Add(user)
		if err != nil {
			log.Println(err)
			res = Result.Error(http.StatusInternalServerError, error.USER_ADD_ERROR)
			goto ret
		}
		res = Result.Success(id)
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Count 统计用户数 GET return sum
func Count(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodGet {
		count, err := UserService.Count(entity.User{})

		if err != nil {
			res = Result.Error(http.StatusInternalServerError, error.USER_SELECT_ERROR)
		} else {
			res = Result.Success(count)
		}
	}
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Update 修改用户信息 POST
func Update(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodPost {
		var user entity.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			res = Result.Error(http.StatusBadRequest, error.USER_DATA_ERROR)
			goto ret
		}
		err = UserService.Update(user)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, error.USER_UPDATE_ERROR)
			goto ret
		}
		res = Result.Success(nil)
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}
