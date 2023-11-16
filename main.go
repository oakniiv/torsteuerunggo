package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/labstack/echo/v4"
)

var secretCheckKey = "ultra-geheim"
var gateGpioMap = make(map[string]int)

type jsonBody struct {
	Secret string `json:"secret"`
	Gate   string `json:"gate"`
}

func toggleGPIO(pin int) error {
	cmd := exec.Command("gpio", "toggle", fmt.Sprintf("%d", pin)) //https://www.geeksforgeeks.org/fmt-sprintf-function-in-golang-with-examples/
	return cmd.Run()
}

func toggleGate(gate string) {
	gpio, ok := gateGpioMap[gate]

	if !ok {
		return
	}

	toggleGPIO(gpio)
	time.Sleep(time.Second * 1)
	toggleGPIO(gpio)
}

func main() {
	e := echo.New()
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
			return echo.NewHTTPError(http.StatusBadRequest, payloadBindError)
		}

		if payload.Secret != secretCheckKey {
			return c.NoContent(http.StatusForbidden)
		}

		if payload.Gate == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		go toggleGate(payload.Gate)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
