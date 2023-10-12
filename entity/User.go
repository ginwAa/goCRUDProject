package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const salt = "<IAmSalt>%"

var SecretKey = []byte("<IAmKey>")

type User struct {
	Id        int64     `json:"id"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Status    int       `json:"status"` // 1 active 2 frozen
	Role      int       `json:"role"`   // 1 admin	2 normal
	Gender    int       `json:"gender"` // 1 male	2 female
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func (u User) ToString() string {
	return fmt.Sprintf("Id:%d Name:%s Pwd:%s Status:%d Role:%d", u.Id, u.Username, u.Password, u.Status, u.Role)
}

func (u *User) MarkUpdated() {
	u.UpdatedAt = time.Now()
}

func (u *User) MarkCreated() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) MarkDeleted() {
	u.DeletedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// md5Hash hash the string and return the hash value
func md5Hash(s *string) string {
	hasher := md5.New()
	hasher.Write([]byte(*s))
	*s = hex.EncodeToString(hasher.Sum(nil))
	return *s
}

func SaltMD5Hash(s string) string {
	s += salt
	md5Hash(&s)
	return s
}

func (u User) CheckPassword(s string) bool {
	s += salt
	return md5Hash(&s) == u.Password
}
