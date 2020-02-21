package api

import (
	"fmt"
	// "net/http"

	// "github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"

	"chitchat0/model"
)

// type UserController struct {
// 	Ctx iris.Context
// }

func IndexHandler(ctx iris.Context) {
	threads, err := model.Threads()
	// for _, thread := range threads {
	// 	fmt.Println(thread.UserId, thread.Id, thread.Topic, "\n")
	// }
	if err != nil {
		danger(err, "cannot get threads")
	} else {
		_, err := session(ctx)
		if err != nil {
			generateHTML(ctx, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(ctx, threads, "layout", "private.navbar", "index")
		}
	}
}

func LoginHandler(ctx iris.Context) {
	t := parseTemplateFiles("login.layout", "login")
	t.Execute(ctx.ResponseWriter(), nil)
}

func AuthenticateHandler(ctx iris.Context) {
	email := ctx.FormValue("email")
	user, _ := model.UserByEmail(email)
	password := ctx.FormValue("password")
	if user.Password == model.Encrypt(password) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		ctx.RemoveCookie("_cookie")
		ctx.SetCookieKV("_cookie", session.Uuid)
		ctx.Redirect("/", 302)
	} else {
		fmt.Println("wrong password!")
		ctx.Redirect("/login", 302)
	}
}

func LogoutHandler(ctx iris.Context) {
	sess, _ := session(ctx)
	cookie := sess.Uuid
	if cookie != "" { // ? 为什么这里是 != 呢？
		session := model.Session{Uuid: cookie}
		session.DeleteByUUID()
	}
	ctx.Redirect("/", 302)
}

func SignupHandler(ctx iris.Context) {
	generateHTML(ctx, nil, "login.layout", "signup")
}

func Signup_accountHandler(ctx iris.Context) {
	user := model.User{
		Username:		ctx.FormValue("name"),
		Email:		ctx.FormValue("email"),
		Password:	model.Encrypt(ctx.FormValue("password")),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	ctx.Redirect("/login", 302)
}