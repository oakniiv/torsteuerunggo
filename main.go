package main

import (
	"fmt"
	"log"
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

//TODO: MODE OUT nur nach Neustart?

func initGPIO(pin int) error {

	cmd := exec.Command("gpio", "mode", fmt.Sprintf("%d", pin), "OUT") //int, int8 etc.: %d
	return cmd.Run()
}

// 'gpio toggle' macht nur probleme
func toggleGPIO(pin int) error {

	out, err := exec.Command("gpio", "read", fmt.Sprintf("%d", pin)).Output()
	if err != nil {
		log.Fatal(err)
	}
	state := string(out)

	fmt.Printf("ZUSTAND: %s\n", state)

	if state == "1\n" { //AUS

		time.Sleep(time.Second * 1)
		fmt.Printf("BUTTON PRESS \n")
		err := exec.Command("gpio", "write", fmt.Sprintf("%d", pin), "0").Run() //EIN
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 2)

		err = exec.Command("gpio", "write", fmt.Sprintf("%d", pin), "1").Run() //AUS
		if err != nil {
			return err
		}
		fmt.Printf("BUTTON RELEASE \n")
		time.Sleep(time.Second * 10)
		

	} else {
		fmt.Printf("WAR SCHON EIN, SOLLTE NICHT PASSIEREN \n")
		time.Sleep(time.Second * 1)
		err = exec.Command("gpio", "write", fmt.Sprintf("%d", pin), "1").Run() //AUS
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil //todo?
}

func toggleGate(gate string) {

	gpio, ok := gateGpioMap[gate]

	if !ok {
		fmt.Print("NOT OK \n")
		return
	}

	initGPIO(gpio)
	toggleGPIO(gpio)
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
			fmt.Print("payloadBindError \n")
			return echo.NewHTTPError(http.StatusBadRequest, payloadBindError)
		}

		if !strings.HasSuffix(payload.Email, "@b-ite.de") && !strings.HasSuffix(payload.Email, "@b-ite.com") && !strings.HasSuffix(payload.Email, "@b-ite.net") {
			fmt.Print("NOT BITE \n")
			return c.NoContent(http.StatusForbidden)
		}

		if payload.Gate == "" {
			fmt.Print("payload gate empty \n")
			return c.NoContent(http.StatusBadRequest)
		}

		go toggleGate(payload.Gate)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
