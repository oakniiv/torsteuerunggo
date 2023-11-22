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

var gateGpioMap = make(map[string]int)

type jsonBody struct {
	Gate  string `json:"gate"`
	Email string `json:"userEmail"` // kommt vom Frontend json, nur überprüfen ob mail string mit b-ite endet
}

func initialiseGPIO() error {

	cmd := exec.Command("gpio", "mode", "22", "OUT", "&&", "gpio", "mode", "23", "OUT", "&&", "gpio", "mode", "24", "OUT", "&&", "gpio", "mode", "25", "OUT") //fix für toggle nach dem Neustart
	return cmd.Run()
}

func toggleGPIO(pin int) error {
	initialiseGPIO()

	cmd := exec.Command("gpio", "toggle", fmt.Sprintf("%d", pin)) //int, int8 etc.: %d
	return cmd.Run()
}

func toggleGate(gate string) {

	gpio, ok := gateGpioMap[gate]

	if !ok {
		fmt.Print("NOT OK")
		return
	}

	fmt.Print("BUTTON PRESS ")
	toggleGPIO(gpio)
	time.Sleep(time.Second * 1)
	toggleGPIO(gpio)
	fmt.Print("BUTTON RELASE ")
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

		var payload = new(jsonBody)

		payloadBindError := c.Bind(payload)

		if payloadBindError != nil {
			fmt.Print("payloadBindError")
			return echo.NewHTTPError(http.StatusBadRequest, payloadBindError)
		}

		if !strings.HasSuffix(payload.Email, "@b-ite.de") && !strings.HasSuffix(payload.Email, "@b-ite.com") && !strings.HasSuffix(payload.Email, "@b-ite.net") {
			fmt.Print("NOT BITE")
			return c.NoContent(http.StatusForbidden)
		}

		if payload.Gate == "" {
			fmt.Print("payload gate empty")
			return c.NoContent(http.StatusBadRequest)
		}

		go toggleGate(payload.Gate)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
