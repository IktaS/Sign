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
	"github.com/unidoc/unipdf/v3/common/license"

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

	err = license.SetMeteredKey("32f87af5351df5bdbc0c0568a0d0bbada1141c3d23ac85fe666e7210e8f826e4")
	if err != nil {
		app.Logger.Fatal(err)
		return
	}

	signRepo, err := sqlite.NewSQLiteDB(db)
	if err != nil {
		app.Logger.Fatal(err)
		return
	}

	signService := service.NewSignService(signRepo, os.Getenv("VERIFY_PATH"))

	app.Static("/public", "public")

	app.GET("/", handler.IndexHandler)
	app.POST("/verify", handler.VerifyFileHandler(signService))
	app.GET("/verify", handler.VerifyHandler(signService))
	app.POST("/sign", handler.SignHandler(signService))

	app.Logger.Fatal(app.Start(":4000"))
}
