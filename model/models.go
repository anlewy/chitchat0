package model

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	// "github.com/jinzhu/gorm"
)

var Db *sql.DB

var Models = []interface{}{
	&User{}, &Session{}, &Thread{}, &Post{},
}

type Model struct {
	Id int `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type Thread struct {
	Model
	Uuid 		string			`gorm:"size:512;unique" json:"uuid" form:"uuid"`
	Topic		string			`gorm:"size:128" json:"topic" form:"topic"`
	UserId		int				`gorm: "size:128" json:"userid" form:"userid"`
	CreateTime 	int64			`json:"createTime" form:"createTime"`
	Posts 		[]Post
}

type Post struct {
	Model
	Uuid		string			`gorm:"size:512;unique" json:"uuid" form:"uuid"`
	Body		string			`gorm:"size:2048" json:"body" form:"body"`
	UserId		int				`gorm: "size:128" json:"userid" form:"userid"`
	ThreadId	int				`gorm: "size:128" json:"threadid" form:"threadid"`
	CreateTime	int64			`json:"createTime" form:"createTime"`
}


// User ...
type User struct {
	Model
	Uuid		string			`gorm:"size:512;unique" json:"uuid" form:"uuid"`
	Username	string 			`gorm:"size:32;unique;" json:"username" form:"username"`
	Email		string			`gorm:"size:128;unique;" json:"email" form:"email"`
	Password	string			`gorm:"size:512" json:"password" form:"password"`
	CreateTime	int64			`json:"createTime" form:"createTime"`
}

// Session ...
type Session struct {
	Model
	Uuid		string			`gorm:"size:512;unique" json:"uuid" form:"uuid"`
	Email		string			`gorm:"size:128;" json:"email" form:"email"`
	UserId		int				`gorm: "size:128" json:"userid" form:"userid"`
	CreateTime	int64			`json:"createTime" form:"createTime"`
}



// CreateSession ...
func (user *User) CreateSession() (session Session, err error) {
	session = Session{
		Uuid:		createUUID(),
		Email:		user.Email,
		UserId:		user.Id,
	}
	db.Create(&session)
	db.NewRecord(session)
	return
}

// Create ...
func (user *User) Create() (err error) {
	user.Uuid = createUUID()
	db.Create(&user)
	db.NewRecord(user)
	err = nil
	return
}

// Check ...
func (session *Session) Check() (valid bool, err error) {
	if err := db.First(&session, "uuid = ?", session.Uuid).Error; err!= nil {
		valid = false
	} else {
		valid = true
	}
	return
}

// DeleteByUUID ...
// 在sessions表中根据uuid删除一条记录，这个函数是由
// 一个session对象调用的
func (session *Session) DeleteByUUID() (err error) {
	err = db.Delete(Session{}, "uuid = ?", session.Uuid).Error
	return
}

func (session *Session) User() (user User, err error) {
	user = User{}
	err = db.First(&user, "id = ?", session.UserId).Error
	return
}

// UserByEmail ...
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = db.First(&user, "email = ?", email).Error
	return
}

// CreateThread ...
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	conv = Thread{
		Uuid:		createUUID(),
		Topic:		topic,
		UserId:		user.Id,
	}
	db.Create(&conv)
	db.NewRecord(conv)
	return
}

// Threads ...
func Threads() (threads []Thread, err error) {
	err = db.Find(&threads).Error
	return
}


// ThreadByUUID ...
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = db.First(&conv, "uuid = ?", uuid).Error
	return
}


// CreatePost ...
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	post = Post{
		Uuid:		createUUID(),
		Body:		body,
		UserId:		user.Id,
		ThreadId:	conv.Id,
	}
	db.Create(&post)
	db.NewRecord(post)
	return
}

// GetPosts ...
func (thread *Thread) GetPosts() (err error) {
	rows := []Post{}
	err = db.Where("thread_id = ?", thread.Id).Find(&rows).Error
	thread.Posts = rows
	return
}


// createUUID ...
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	u[8] = (u[8] | 0x40) & 0x7f
	u[6] = (u[6] & 0xf) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Encrypt ...
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}