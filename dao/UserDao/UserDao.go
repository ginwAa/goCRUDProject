package UserDao

import (
	"project1/conf"
	"project1/entity"
	"strconv"
)

func PageByUsername(s string, page int, count int) ([]entity.User, error) {
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status FROM users WHERE username like CONCAT('%', $1::text, '%') LIMIT $2 OFFSET $3", s, count, page)
	if err != nil {
		return nil, err
	}
	var users []entity.User
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.Id, &u.Username, &u.UpdatedAt, &u.CreatedAt, &u.Status)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func Page(page int, count int) ([]entity.User, error) {
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status FROM users LIMIT $1 OFFSET $2", count, count*page)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var users []entity.User
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.Id, &u.Username, &u.UpdatedAt, &u.CreatedAt, &u.Status)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func Add(user entity.User) (int64, error) {
	var id int64
	err := conf.PGDB.QueryRow("INSERT INTO users (username, password, updated_at, created_at, status, granted) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Username, user.Password, user.UpdatedAt, user.CreatedAt, user.Status, user.Granted).Scan(&id)
	return id, err
}

func Sum() (int64, error) {
	res, err := conf.PGDB.Query("SELECT COUNT(1) FROM users")
	if err != nil {
		return 0, err
	}
	defer res.Close()
	res.Next()
	var sum int64
	err = res.Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, res.Err()
}

func Update(user entity.User) error {
	sql := "UPDATE users SET updated_at = $1"
	if user.Username != "" {
		sql += ",username='" + user.Username + "'"
	}
	if user.Password != "" {
		sql += ",password='" + user.Password + "'"
	}
	if user.Status != 0 {
		sql += ",status=" + strconv.Itoa(user.Status)
	}
	if user.Granted != 0 {
		sql += ",granted=" + strconv.Itoa(user.Granted)
	}
	sql += " WHERE id = $2"
	_, err := conf.PGDB.Exec(sql, user.UpdatedAt, user.Id)
	return err
}
