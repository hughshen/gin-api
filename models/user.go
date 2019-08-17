package models

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"gin-api/helpers"
	"strings"
	"time"
)

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username" form:"username"`
	Password    string `json:"password" form:"password"`
	ApiToken    string `json:"token"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (u User) Validate() string {
	if u.Username == "" {
		return "The username field is required"
	} else if u.ID == 0 && u.Password == "" {
		return "The password field is required"
	} else if nameExists := u.NameExists(); nameExists == true {
		return "The username field is exists"
	}

	return ""
}

func (lf LoginForm) Validate() (User, string) {
	var user User

	if lf.Username == "" || lf.Password == "" {
		return user, "Username and password is required"
	} else {

		lf.Password = helpers.Md5Hash(lf.Password)

		user.RowScan(db.QueryRow("select id, username, created_at, updated_at from user where username = ? and password = ?", lf.Username, lf.Password))
		if user.ID == 0 {
			return user, "Incorrect username or password"
		}

		// Update token
		token := UserCreateToken(16)
		_, err := db.Exec("update user set api_token = ? where id = ?", token, user.ID)
		if err != nil {
			ShowErrLog(err)
			return user, "Login failed"
		}

		user.ApiToken = token
		user.BuildUnixDate()

		return user, ""
	}
}

func (u User) NameExists() bool {
	var count int

	row := db.QueryRow("select count(id) from user where username = ? and id != ?", u.Username, u.ID)
	err := row.Scan(&count)
	if err != nil {
		ShowErrLog(err)
		return true
	}

	return count > 0
}

func (u *User) BuildUnixDate() {
	createdUnix := time.Unix(int64(u.CreatedAt), 0)
	updatedUnix := time.Unix(int64(u.UpdatedAt), 0)

	u.CreatedTime = createdUnix.Format("2006-01-02 15:04:05")
	u.UpdatedTime = updatedUnix.Format("2006-01-02 15:04:05")
}

func (u *User) RowScan(row *sql.Row) {
	err := row.Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil || err == sql.ErrNoRows {
		ShowErrLog(err)
	}
}

func (u *User) RowsScan(rows *sql.Rows) {
	err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil || err == sql.ErrNoRows {
		ShowErrLog(err)
	}
}

func (u User) Create() int64 {
	// Password
	u.Password = helpers.Md5Hash(u.Password)

	stmt, err := db.Prepare("insert into user set username = ?, password = ?, created_at = ?, updated_at = ?")
	if err != nil {
		ShowErrLog(err)
		return 0
	}
	defer stmt.Close()

	current := int32(time.Now().Unix())
	res, err := stmt.Exec(u.Username, u.Password, current, current)
	if err != nil {
		ShowErrLog(err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		ShowErrLog(err)
		return 0
	}

	return id
}

func (u User) Update() int64 {
	var err error
	var res sql.Result

	current := int32(time.Now().Unix())
	if u.Password == "" {
		res, err = db.Exec("update user set username = ?, updated_at = ? where id = ?", u.Username, current, u.ID)
	} else {
		u.Password = helpers.Md5Hash(u.Password)
		res, err = db.Exec("update user set username = ?, password = ?, updated_at = ? where id = ?", u.Username, u.Password, current, u.ID)
	}
	if err != nil {
		ShowErrLog(err)
		return 0
	}

	affect, err := res.RowsAffected()
	if err != nil {
		ShowErrLog(err)
		return 0
	}

	return affect
}

func UserValidateToken(token string) User {
	var user User
	user.RowScan(db.QueryRow("select id, username, created_at, updated_at from user where api_token = ?", token))

	return user
}

func UserCreateToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return strings.ToUpper(fmt.Sprintf("%x", b))
}

func UserDeleteToken(id int) bool {
	_, err := db.Exec("update user set api_token = '' where id = ?", id)
	if err != nil {
		ShowErrLog(err)
		return false
	}

	return true
}

func UserGetList() []User {
	users := make([]User, 0)

	rows, err := db.Query("select id, username, created_at, updated_at from user order by created_at desc, id desc")
	if err != nil {
		ShowErrLog(err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		user.RowsScan(rows)
		if user.ID == 0 {
			return users
		}

		user.BuildUnixDate()
		users = append(users, user)
	}

	return users
}

func UserGetById(id int64) User {
	var user User
	user.RowScan(db.QueryRow("select id, username, created_at, updated_at from user where id = ?", id))

	return user
}

func UserDeleteById(id int64) bool {
	stmt, err := db.Prepare("delete from user where id = ?")
	if err != nil {
		ShowErrLog(err)
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		ShowErrLog(err)
		return false
	}

	_, err = res.RowsAffected()
	if err != nil {
		ShowErrLog(err)
		return false
	}

	return true
}
