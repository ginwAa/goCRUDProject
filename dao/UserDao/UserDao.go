package UserDao

import (
	"project1/conf"
	"project1/entity"
)

func PageByUsername(s string, page int, size int) ([]entity.User, error) {
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status, account, gender, phone, email, role FROM users WHERE username like CONCAT('%', $1::text, '%') and deleted_at is NULL ORDER BY id LIMIT $2 OFFSET $3", s, size, (page-1)*size)
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

func Page(page int, size int) ([]entity.User, error) {
	rows, err := conf.PGDB.Query("SELECT id, username, updated_at, created_at, status, account, gender, phone, email, role FROM users WHERE deleted_at is NULL ORDER BY id LIMIT $1 OFFSET $2", size, (page-1)*size)
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
	row := conf.PGDB.QueryRow("INSERT INTO users (username, password, updated_at, created_at, status, role, account, gender, phone, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", user.Username, user.Password, user.UpdatedAt, user.CreatedAt, user.Status, user.Role, user.Account, user.Gender, user.Phone, user.Email)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SelectOneByLogin(account string, password string) (entity.User, error) {
	res := conf.PGDB.QueryRow("SELECT id, username, status, role, account, gender, phone, email FROM users WHERE deleted_at is NULL and account = $1 and password = $2", account, password)
	err := res.Err()
	var ret entity.User
	if err != nil {
		return ret, err
	}
	err = res.Scan(&ret.Id, &ret.Username, &ret.Status, &ret.Role, &ret.Account, &ret.Gender, &ret.Phone, &ret.Email)
	return ret, err
}

func CountByAccount(account string) (int64, error) {
	res, err := conf.PGDB.Query("SELECT COUNT(1) FROM users WHERE account = $1 and deleted_at is NULL", account)
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

func CountByUsernameLike(username string) (int64, error) {
	res, err := conf.PGDB.Query("SELECT COUNT(1) FROM users WHERE username like CONCAT('%', $1::text, '%')and deleted_at is NULL", username)
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

func UpdateUsername(user entity.User) (int64, error) {
	res, err := conf.PGDB.Exec("UPDATE users SET updated_at = $1, username = $2 WHERE id = $3", user.UpdatedAt, user.Username, user.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdatePassword(user entity.User) (int64, error) {
	res, err := conf.PGDB.Exec("UPDATE users SET updated_at = $1, password = $2 WHERE id = $3", user.UpdatedAt, user.Password, user.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Delete(user entity.User) (int64, error) {
	res, err := conf.PGDB.Exec("UPDATE users SET deleted_at = $1, updated_at = $2 WHERE id = $3", user.DeletedAt, user.UpdatedAt, user.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
