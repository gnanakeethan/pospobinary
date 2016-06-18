package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gnanakeethan/print/escpos"
	"io"
	"log"
	"net/http"
	"os"
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
	//get variables of post
	print := r.FormValue("print")
	machine := r.FormValue("machine")
	printer := r.FormValue("printer")
	barcode := r.FormValue("barcode")
	code_format := r.FormValue("code_format")
	print = strings.Replace(print, "\\n", "\n", -1)
	fmt.Printf("%+v\n", r.Form)

	log.Println(barcode)
	log.Println(code_format)
	log.Println(print)

	printmachine := "\\\\" + machine + "\\" + printer
	f := bytes.NewBuffer([]byte(""))
	w := bufio.NewWriter(f)

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
	textmap := make(map[string]string)
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			textmap[key] = value
			fmt.Println(key, value)
		}
	}
	fmt.Printf("%+v\n", textmap)
	if print != "" {
		p.Text(textmap, print)
		//	p.Write(print)
	}
	p.End()
	w.Flush()

	log.Print(f)
	copyFileContents(f, printmachine)
}

func copyFileContents(in io.Reader, dst string) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
