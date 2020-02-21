package main

import (
	// "fmt"
	"net/http"
	"os"
	"os/signal"
	// "strings"
	"syscall"

	// "github.com/go-resty/resty/v2"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	// "github.com/kataras/iris/v12/mvc"
	// "github.com/mlogclub/simple"
	log "github.com/sirupsen/logrus"

	"chitchat0/common/config"
	"chitchat0/controllers/api"
	"chitchat0/model"
)

// InitIris kkk
func initIris() {
	app := iris.New()
	app.Logger().SetLevel("warn")
	app.Use(recover.New())
	app.Use(logger.New())
	// cors 是一个处理http跨域请求的中间件
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))
	app.AllowMethods(iris.MethodOptions)

	// // 这个是不能照搬的，如果要这段发挥作用，则需要改变它
	// app.OnAnyErrorCode(func(ctx iris.Context) {
	// 	path := ctx.Path()
	// 	var err error
	// 	if strings.Contains(path, "/api/admin/") {
	// 		_, err = ctx.JSON(simple.JsonErrorCode(ctx.GetStatusCode(), "Http error"))
	// 	}
	// 	if err != nil {
	// 		logrus.Error(err)
	// 	}
	// })

	

	// api
	// mvc.Configure(app.Party("/"), func(m *mvc.Application) {
	// 	m.Party("/topic").Handle(new(api.TopicController))
	// 	m.Party("/article").Handle(new(api.ArticleController))
	// 	m.Party("/project").Handle(new(api.ProjectController))
	// })

	// thread
	// user_p := app.Party("/user", userHandler)
	app.Any("/", api.IndexHandler)
	app.Any("/index", api.IndexHandler)
	app.Get("/login", api.LoginHandler)
	app.Post("/authenticate", api.AuthenticateHandler)
	app.Get("/logout", api.LogoutHandler)
	app.Get("/signup", api.SignupHandler)
	app.Post("/signup_account", api.Signup_accountHandler)

	app.Get("/thread/new", api.ThreadNewHandler)
	app.Post("/thread/create", api.ThreadCreateHandler)
	app.Get("/thread/read", api.ThreadReadHandler)
	app.Post("/thread/post", api.ThreadPostHandler)



	server := &http.Server{Addr: ":" + config.Conf.Port}
	handleSignal(server)
	err := app.Run(iris.Server(server), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	}))
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Infof("got signal [%s], exiting now", s)
		if err := server.Close(); nil != err {
			log.Errorf("server close failed: " + err.Error())
		}

		model.CloseDB()

		log.Infof("Exited")
		os.Exit(0)
	}()
}
