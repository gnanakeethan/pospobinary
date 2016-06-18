package main

import (
	"bufio"
	"github.com/gnanakeethan/print/escpos"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
)

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServeTLS(":8080", "server.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(wr http.ResponseWriter, r *http.Request) {
	os.Remove("output")

	//get variables of post
	print := r.FormValue("print")
	machine := r.FormValue("machine")
	printer := r.FormValue("printer")
	barcode := r.FormValue("barcode")
	code_format := r.FormValue("code_format")
	print = strings.Replace(print, "\\n", "\n", -1)

	log.Println(barcode)
	log.Println(code_format)
	log.Println(print)
	//print command
	b := "copy output \\\\" + machine + "\\" + printer
	ioutil.WriteFile("run.cmd", []byte(b), 0644)
	cmd := exec.Command("run.cmd")
	defer cmd.Run()

	path := "output"
	// detect if file exists
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		// create file if not exists
		var file, _ = os.Create(path)
		defer file.Close()
	}

	f, err := os.OpenFile("output", O_WRONLY, 0644)
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := escpos.New(w)
	p.Init()
	if barcode != "" {
		if code_format == "" {
			code_format = "1"
		}
		barcode = strings.ToUpper(barcode)
		intformat, _ := strconv.Atoi(code_format)
		p.PrintBarcode(barcode, intformat)
	}
	if print != "" {
		p.Write(print)
	}

	p.End()

	w.Flush()
}
