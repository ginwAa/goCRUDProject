package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const salt string = "<IAmSalt>%"

type User struct {
	Id        int64     `json:"id"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Status    int       `json:"status"`
	Role      int       `json:"role"`
	Gender    int       `json:"gender"`
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
}

// md5Hash hash the string and return the hash value
func md5Hash(s *string) string {
	hasher := md5.New()
	hasher.Write([]byte(*s))
	*s = hex.EncodeToString(hasher.Sum(nil))
	return *s
}

func (u *User) SaltMD5HashPassword() {
	u.Password += salt
	md5Hash(&u.Password)
	fmt.Println("hash done: " + u.Password)
}

func (u User) CheckPassword(s string) bool {
	s += salt
	return md5Hash(&s) == u.Password
}
