package main

import (
	"database/sql"
	"os"

	"github.com/IktaS/sign/handler"
	"github.com/IktaS/sign/service"
	"github.com/IktaS/sign/sqlite"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := echo.New()
	err := godotenv.Load()
	if err != nil {
		app.Logger.Fatal(err)
		return
	}

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		app.Logger.Fatal(err)
		return
	}

	defer db.Close()

	signRepo, err := sqlite.NewSQLiteDB(db)
	if err != nil {
		app.Logger.Fatal(err)
		return
	}

	signService := service.NewSignService(signRepo, os.Getenv("VERIFY_PATH"))

	app.Static("/public", "public")

	app.GET("/", handler.IndexHandler)
	app.POST("/verify", handler.VerifyFileHandler(signService))
	app.GET("/verify", handler.VerifyHandler(os.Getenv("OWNER_NAME"), signService))
	app.POST("/sign", handler.SignHandler(signService))

	app.Logger.Fatal(app.Start(":4000"))
}
