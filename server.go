package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pieterclaerhout/go-log"
)

const secret = "secret"

func login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "pieter" || password != "claerhout" {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Pieter Claerhout"
	claims["uuid"] = "9E98C454-C7AC-4330-B2EF-983765E00547"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	log.InfoDump(claims, "claims")
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)
	e.GET("/", accessible)

	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte(secret)))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":8080"))

}
