package main

import (
    "fmt"
    "net/http"
    "github.com/warthog618/gpiod"
)

// PIN_CH1     = 25
// PIN_CH2     = 24
// PIN_CH3     = 23
// PIN_CH4     = 22

func toggleGPIO() {
    chip, err := gpiod.NewChip("gpiochip0")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer chip.Close()

    line, err := chip.RequestLine(25, gpiod.AsOutput(0)) // PIN_CH4     = 22
    if err != nil {
        fmt.Println(err)
        return
    }
    defer line.Close()

    val, err := line.Value()
    if err != nil {
        fmt.Println(err)
        return
    }

    if val == 0 {
        line.SetValue(1)
        fmt.Println("SANITY CHECK 1")
    } else {
        line.SetValue(0)
        fmt.Println("SANITY CHECK 2")
    }
}

func buttonHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        toggleGPIO()
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
