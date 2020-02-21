package api

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	// "net/http"
	// "strings"

	"github.com/kataras/iris/v12"

	"chitchat0/model"
)

var logger *log.Logger

// // error_message ...
// func error_message(w http.ResponseWriter, r *http.Request, msg string) {
// 	url := []string{"/err?msg=", msg}
// 	http.Redirect(w, r, strings.Join(url, ""), 302)
// }


// 首先检查有没有会话，如果没有则直接返回错误，否则检查数据库中是否有
// 这个会话记录（使用check()函数），如果检查不对则返回错误，否则返回
// 会话和为nil值的err
func session(ctx iris.Context) (sess model.Session, err error) {
	cookie := ctx.GetCookie("_cookie")
	err = nil
	sess = model.Session{Uuid: cookie}
	if ok, _ := sess.Check(); !ok {
		err = errors.New("Invalid session")
		fmt.Println("maybe here")
	}
	return
}

func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}


func generateHTML(c iris.Context, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./templates/%s.html", file))
	}
	// fmt.Println(files)
	templates := template.Must(template.ParseFiles(files...))  // 为什么这里要加三个点呢？
	templates.ExecuteTemplate(c.ResponseWriter(), "layout", data)
}


func warning(args ...interface{}) {
	logger.SetPrefix("WARNING")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}