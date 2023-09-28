package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const salt string = "<IAmSalt>%"

type User struct {
	Id        int64
	Username  string
	Password  string
	Status    int
	Granted   int
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (u User) ToString() string {
	return fmt.Sprintf("Id:%d Name:%s Pwd:%s Status:%d Granted:%d", u.Id, u.Username, u.Password, u.Status, u.Granted)
}

func (u *User) MarkUpdated() {
	u.UpdatedAt = time.Now()
}

func (u *User) MarkCreated() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
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

// GenerateRandomString securely generated random string.
//
//	func GenerateRandomString() (string, error) {
//		b := make([]byte, 64)
//		_, err := rand.Read(b)
//		if err != nil {
//			return "", err
//		}
//		return base64.URLEncoding.EncodeToString(b), err
//	}
