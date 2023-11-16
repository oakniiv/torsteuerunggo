package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

func toggleGPIO(pin int) error {
	cmd := exec.Command("gpio", "toggle", fmt.Sprintf("%d", pin)) //https://www.geeksforgeeks.org/fmt-sprintf-function-in-golang-with-examples/
	return cmd.Run()
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlContent)
	})

	e.POST("/", func(c echo.Context) error {
		buttonValue := c.FormValue("button")

		if buttonValue == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		switch buttonValue {
		case "open1":
			toggleGPIO(25) // Pin 25 TOR 2
		case "open2":
			toggleGPIO(24) // Pin 24 TOR 3
		case "open3":
			toggleGPIO(23) // Pin 23 TOR 4
		case "open4":
			toggleGPIO(22) // Pin 22 TOR 1
		}

		return c.HTML(http.StatusOK, htmlContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

const htmlContent = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Torsteuerung</title>
<style>
    body {
        background: #FFFFFF; 
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        margin: 0;
        font-family: Arial, sans-serif;
        color: #333;
    }

    .container {
        text-align: center;
        padding: 20px;
        max-width: 300px;
    }

    .logo {
        max-width: inherit;
    }


    .neu-button {
        background: #e0e5ec; 
        border: none;
        padding: 15px 30px; 
        margin-bottom: 25px; 
        border-radius: 25px; 
        box-shadow: 8px 8px 15px #a3b1c6, -8px -8px 15px #ffffff; 
        font-size: 16px;
        color: #3291ff; /* BITE Blau */
        cursor: pointer;
        outline: none;
        transition: box-shadow 0.3s ease-in-out, transform 0.1s;
        width: 100%; 
        box-sizing: border-box; 
    }

    .neu-button:hover {
        box-shadow: 4px 4px 10px #a3b1c6, -4px -4px 10px #ffffff;
    }

    .neu-button:active {
        box-shadow: inset 2px 2px 5px #a3b1c6, inset -2px -2px 5px #ffffff;
        transform: translateY(2px);
    }
</style>
</head>
<body>
    <div class="container">
    <img class="logo" src="https://getlogo.net/wp-content/uploads/2020/05/bite-gmbh-logo-vector.png" alt="BITE GmbH Logo">
        <form method="post">
            <button class="neu-button" type="submit" name="button" value="open1"><strong>Tor 1 Öffnen</strong></button>
            <button class="neu-button" type="submit" name="button" value="open2"><strong>Tor 2 Öffnen</strong></button>
            <button class="neu-button" type="submit" name="button" value="open3"><strong>Tor 3 Öffnen</strong></button>
            <button class="neu-button" type="submit" name="button" value="open4"><strong>Tor 4 Öffnen</strong></button>
        </form>
    </div>
</body>
</html>`
