package main

import (
	"fmt"
	"net/http"

	//"github.com/warthog618/gpiod"
	"os/exec"
)

// PIN_CH1     = 25
// PIN_CH2     = 24
// PIN_CH3     = 23
// PIN_CH4     = 22

func toggleGPIO(pin int) error {
	cmd := exec.Command("gpio", "toggle", fmt.Sprintf("%d", pin))
	return cmd.Run()
}

func buttonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		toggleGPIO(25)
		fmt.Fprintf(w, `<html><body><form action="/" method="post"><input type="submit" value="Toggle GPIO"></form></body></html>`)
	} else {
		fmt.Fprintf(w, `<html><body><form action="/" method="post"><input type="submit" value="Toggle GPIO"></form></body></html>`)
	}
}

func main() {
	http.HandleFunc("/", buttonHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
