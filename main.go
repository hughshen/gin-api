package main

import (
	"flag"
	"gin-api/handlers"
	"gin-api/helpers"
	"gin-api/middlewares"
	"gin-api/models"
	"github.com/gin-gonic/gin"
	"html/template"
)

func main() {
	help := flag.Bool("h", false, "Help message")
	host := flag.String("host", "0.0.0.0", "Server host")
	port := flag.String("port", "8080", "Server port")
	dbHost := flag.String("db_host", "localhost", "DB host")
	dbPort := flag.String("db_port", "3306", "DB port")
	dbUser := flag.String("db_user", "root", "DB username")
	dbPass := flag.String("db_pass", "root", "DB password")
	dbName := flag.String("db_name", "test", "DB name")

	flag.Parse()

	if *help == true {
		flag.Usage()
		return
	}

	dsn := *dbUser + ":" + *dbPass + "@tcp(" + *dbHost + ":" + *dbPort + ")/" + *dbName + "?charset=utf8"
	db, err := models.InitDb(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := InitRouter()
	router.Run(*host + ":" + *port)
}

func InitRouter() *gin.Engine {
	// Gin config
	//gin.DisableConsoleColor()

	router := gin.New()

	// Global middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.HeaderMiddleware())

	// Func map
	router.SetFuncMap(template.FuncMap{
		"raw":   helpers.FuncRaw,
		"add":   helpers.FuncAdd,
		"minus": helpers.FuncMinus,
	})

	// Site
	router.LoadHTMLGlob("./views/*.html")
	router.GET("/", handlers.SiteHome)

	// 404
	router.NoRoute(handlers.SiteNotFound)

	// Auth
	router.POST("/api/login", handlers.AuthLogin)

	api := router.Group("/api").Use(middlewares.TokenAuthMiddleware())
	{
		api.GET("/hello", handlers.AuthHello)
		api.PUT("/logout", handlers.AuthLogout)
		// User
		api.GET("/users", handlers.UserIndex)
		api.GET("/users/:id", handlers.UserShow)
		api.POST("/users", handlers.UserCreate)
		api.PUT("/users/:id", handlers.UserUpdate)
		api.DELETE("/users/:id", handlers.UserDelete)
	}

	return router
}
