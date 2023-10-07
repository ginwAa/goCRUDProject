package UserDao

import (
	"project1/conf"
	"project1/entity"
	"strconv"
)

func PageByUsername(s string, page int, count int) ([]entity.User, error) {
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status, account, gender, phone, email, role FROM users WHERE username like CONCAT('%', $1::text, '%') and deleted_at is NULL ORDER BY id LIMIT $2 OFFSET $3", s, count, (page-1)*count)
	if err != nil {
		return nil, err
	}
	var users []entity.User
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.Id, &u.Username, &u.UpdatedAt, &u.CreatedAt, &u.Status, &u.Account, &u.Gender, &u.Phone, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func Page(page int, count int) ([]entity.User, error) {
	//rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status, account, gender, phone, email FROM users LIMIT $1 OFFSET $2", count, count*page)
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status, account, gender, phone, email, role FROM users  WHERE deleted_at is NULL ORDER BY id LIMIT $1 OFFSET $2", count, (page-1)*count)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var users []entity.User
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.Id, &u.Username, &u.UpdatedAt, &u.CreatedAt, &u.Status, &u.Account, &u.Gender, &u.Phone, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func Add(user entity.User) (int64, error) {
	var id int64
	err := conf.PGDB.QueryRow("INSERT INTO users (username, password, updated_at, created_at, status, role, account, gender, phone, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", user.Username, user.Password, user.UpdatedAt, user.CreatedAt, user.Status, user.Role, user.Account, user.Gender, user.Phone, user.Email).Scan(&id)
	return id, err
}

func CountAll() (int64, error) {
	res, err := conf.PGDB.Query("SELECT COUNT(1) FROM users WHERE deleted_at is NULL")
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

func CountByUserLike(user entity.User) (int64, error) {
	res, err := conf.PGDB.Query("SELECT COUNT(1) FROM users WHERE username like CONCAT('%', $1::text, '%') and account like CONCAT('%', $2::text, '%') and deleted_at is NULL", user.Username, user.Account)
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

func CountByUser(user entity.User) (int64, error) {
	sql := "SELECT COUNT(1) FROM users WHERE deleted_at is NULL"
	if user.Username != "" {
		sql += " and username='" + user.Username + "'"
	}
	if user.Account != "" {
		sql += " and account='" + user.Account + "'"
	}
	res, err := conf.PGDB.Query(sql)
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
	if user.Role != 0 {
		sql += ",granted=" + strconv.Itoa(user.Role)
	}
	if user.Gender != 0 {
		sql += ",gender=" + strconv.Itoa(user.Gender)
	}
	if user.Account != "" {
		sql += ",account=" + user.Account + "'"
	}
	if user.Phone != "" {
		sql += ",phone=" + user.Phone + "'"
	}
	if user.Email != "" {
		sql += ",email=" + user.Email + "'"
	}
	sql += " WHERE id = $2"
	_, err := conf.PGDB.Exec(sql, user.UpdatedAt, user.Id)
	return err
}

func Delete(user entity.User) error {
	_, err := conf.PGDB.Exec("UPDATE users SET deleted_at = $1 WHERE id = $2", user.DeletedAt, user.Id)
	return err
}
