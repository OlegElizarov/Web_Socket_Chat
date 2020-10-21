package main

import (
	"Web_Socket_Chat/delivery"
	"Web_Socket_Chat/infrastructure"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Logger())
	db, err := infrastructure.InitDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
	h := delivery.Handler{Db: db}
	e.Use(middleware.Recover())
	e.Static("/", "static")
	e.GET("/ws", h.Send)
	e.POST("/send", h.Receive)
	e.Logger.Fatal(e.Start(":8080"))
	//e.Start(":8080")
}
