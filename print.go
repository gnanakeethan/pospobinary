package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	//"strings"
)

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServeTLS(":8080", "server.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("print")
	machine := r.FormValue("machine")
	printer := r.FormValue("printer")

	log.Print(p)
	log.Print(machine)
	log.Print(printer)

	ioutil.WriteFile("output.txt", []byte(p), 0644)
	b := "copy output.txt \\\\" + machine + "\\" + printer
	ioutil.WriteFile("run.cmd", []byte(b), 0644)

	cmd := exec.Command("run.cmd")
	cmd.Run()
}
