package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var secretCheckKey = "ultra-geheim"
var gateGpioMap = make(map[string]int)

type jsonBody struct {
	Secret string `json:"secret"`
	Gate   string `json:"gate"`
	Email  string `json:"userEmail"` // kommt vom Frontend json, nur überprüfen ob mail string mit b-ite endet oder so
}

func toggleGPIO(pin int) error {
	cmd := exec.Command("gpio", "toggle", fmt.Sprintf("%d", pin))
	return cmd.Run()
}

func toggleGate(gate string) {
	fmt.Print("TOGGLE GATE ENTER")
	gpio, ok := gateGpioMap[gate]

	if !ok {
		fmt.Print("NOT OK")
		return
	}

	toggleGPIO(gpio)
	time.Sleep(time.Second * 1)
	toggleGPIO(gpio)
	fmt.Print("TOGGLEGATE EXIT")
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.HideBanner = true

	gateGpioMap["open1"] = 22
	gateGpioMap["open2"] = 25
	gateGpioMap["open3"] = 24
	gateGpioMap["open4"] = 23

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, ":)")
	})

	e.POST("/api/toggle", func(c echo.Context) error {
		fmt.Print("SANITY CHECK 1")
		var payload = new(jsonBody)

		payloadBindError := c.Bind(payload)

		if payloadBindError != nil {
			fmt.Print("payloadBindError")
			return echo.NewHTTPError(http.StatusBadRequest, payloadBindError)
		}

		if payload.Secret != secretCheckKey {
			fmt.Print("INSIDE SECRET")
			return c.NoContent(http.StatusForbidden)
		}

		if !strings.HasSuffix(payload.Email, "@b-ite.de") { // bite.net auch?
			return c.NoContent(http.StatusUnauthorized)
		}

		if !strings.HasSuffix(payload.Email, "@b-ite.com") {
			return c.NoContent(http.StatusUnauthorized)
		}

		if payload.Gate == "" {
			fmt.Print("payload gate empty")
			return c.NoContent(http.StatusBadRequest)
		}
		fmt.Print("SANITY CHECK 2")
		go toggleGate(payload.Gate)
		fmt.Print("SANITY CHECK 3")
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
