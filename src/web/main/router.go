package main

import (
	"web/controllers"
	"fmt"
	"net/http"
)

func main() {
	home := controllers.Home {}
	http.HandleFunc("/", home.Index);
	
	fmt.Println("Serving requests.")
	err := http.ListenAndServe(":80", nil)
	if err != nil { fmt.Println(err); return }
}
