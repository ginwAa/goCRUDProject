package UserHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"project1/constants"
	"project1/entity"
	"project1/entity/DTO"
	"project1/entity/Result"
	"project1/service/UserService"
	"strconv"
)

// Page 分页查询 GET return user list
func Page(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodGet {
		var user entity.User
		user.Username = r.URL.Query().Get("username")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))

		users, err := UserService.Page(user, page, size)
		if err != nil {
			log.Printf("user/page err: %s\n", err)
			res = Result.Error(http.StatusInternalServerError, constants.USER_SELECT_ERROR)
		} else {
			res = Result.Success(users)
		}
	} else if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
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
			log.Printf("user/add err1:%s\n", err)
			res = Result.Error(http.StatusBadRequest, constants.USER_DATA_ERROR)
			goto ret
		}
		id, err := UserService.Add(user)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_ADD_ERROR)
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
		username := r.URL.Query().Get("username")
		count, err := UserService.Count(username)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_SELECT_ERROR)
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
		if err != nil || user.Id == 0 {
			res = Result.Error(http.StatusBadRequest, constants.USER_DATA_ERROR)
			goto ret
		}
		err = UserService.Update(user)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_UPDATE_ERROR)
			goto ret
		}
		res = Result.Success(nil)
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Delete 删除用户 DELETE
func Delete(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodPost {
		var Id struct {
			Id int64 `json:"id"`
		}
		err := json.NewDecoder(r.Body).Decode(&Id)
		if err != nil || Id.Id == 0 {
			res = Result.Error(http.StatusBadRequest, constants.USER_DATA_ERROR)
			goto ret
		}
		err = UserService.Delete(Id.Id)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_DELETE_ERROR)
			goto ret
		}
		res = Result.Success(true)
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Login POST 登录
func Login(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodPost {
		var user DTO.UserLoginDTO
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_LOGIN_ERROR)
			goto ret
		}
		loginVO, err := UserService.Login(user)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_SELECT_ERROR)
		} else {
			res = Result.Success(loginVO)
			cookie := http.Cookie{
				Name:   "TKARL",
				Value:  loginVO.Token,
				Path:   "/",
				Domain: "localhost",
				//Expires:    time.Now().Add(time.Hour),
				RawExpires: "",
				MaxAge:     0,
				Secure:     false,
				HttpOnly:   true,
				SameSite:   0,
				Raw:        "",
				Unparsed:   nil,
			}
			http.SetCookie(w, &cookie)
		}
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// Register 添加用户 POST return user id
func Register(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodPost {
		var user entity.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			res = Result.Error(http.StatusBadRequest, constants.USER_DATA_ERROR)
			goto ret
		}
		user.Id = 0
		user.Status = 1
		user.Role = 2
		id, err := UserService.Add(user)
		if err != nil {
			log.Println(err)
			res = Result.Error(http.StatusInternalServerError, constants.USER_ADD_ERROR)
			goto ret
		}
		res = Result.Success(id)
	}
ret:
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

// AccountUnique 检查Account是否已被占用
func AccountUnique(w http.ResponseWriter, r *http.Request) {
	res := Result.BadRequest()

	if r.Method == http.MethodGet {
		account := r.URL.Query().Get("account")
		success, err := UserService.CheckAccountUnique(account)
		if err != nil {
			res = Result.Error(http.StatusInternalServerError, constants.USER_SELECT_ERROR)
		} else {
			res = Result.Success(success)
		}
	}
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}
