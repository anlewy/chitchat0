package api

import (
	"fmt"
	// "net/http"

	// "github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"

	"chitchat0/model"
)

// type ThreadController struct {
// 	Ctx iris.Context
// }


// 获取新建thread的页面
// 即在这个页面可以编辑新建的thread的内容
// 这个页面提交后作为post create
func ThreadNewHandler(ctx iris.Context) {
	_, err := session(ctx)
	if err != nil {
		fmt.Println("you should login!")
		ctx.Redirect("/login", 302)
	} else {
		generateHTML(ctx, nil, "layout", "private.navbar", "new.thread")
	}
}


// 创建新的thread
func ThreadCreateHandler(ctx iris.Context) {
	sess, err := session(ctx)
	log.Info("here is ThreadCreateHandler()")
	log.Infof("here is ThreadCreateHandler(): %v", err)
	if err != nil {
		ctx.Redirect("/login", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		topic := ctx.FormValue("topic")
		if _, err := user.CreateThread(topic); err != nil {
			danger(err, "Cannot create thread")
		}
		ctx.Redirect("/index", 302) // 创建结束之后到主页去
	}
}


// 阅读特定thread
func ThreadReadHandler(ctx iris.Context) {
	uuid := ctx.FormValue("id")
	thread, err := model.ThreadByUUID(uuid)
	if err != nil {
		// error_message(w, r, "Cannot read thread")
		fmt.Println("nothing")
	} else {
		// 现在想要为一个thread查询到所有他的posts
		thread.GetPosts()
		_, err := session(ctx)
		if err != nil {
			generateHTML(ctx, &thread, "layout", "public.navbar", "public.thread")
		} else {
			generateHTML(ctx, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}


// 对特定Thread添加评论
func ThreadPostHandler(ctx iris.Context) {
	sess, err := session(ctx)
	if err != nil {
		ctx.Redirect("/login", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := ctx.FormValue("body")
		uuid := ctx.FormValue("uuid")
		thread, err := model.ThreadByUUID(uuid)
		if err != nil {
			danger(err, "cannot read thread")
		}
		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		ctx.Redirect(url, 302)
	}
}
