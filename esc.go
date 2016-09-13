package main

import (
	"bufio"
	"bytes"
	"github.com/gnanakeethan/print/escpos"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServeTLS(":8080", "server.pem", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(wr http.ResponseWriter, r *http.Request) {

	wr.Header().Add("Access-Control-Allow-Origin", "*")

	print := r.FormValue("print")
	machine := r.FormValue("machine")
	printer := r.FormValue("printer")
	barcode := r.FormValue("barcode")
	image := r.FormValue("image")
	code_format := r.FormValue("code_format")

	log.Println(image)
	ostype := runtime.GOOS
	printmachine := "testfile"
	switch ostype {
	case "windows":
		printmachine = "\\\\" + machine + "\\" + printer
	case "linux":
		printmachine = "/dev/" + printer
	}

	f := bytes.NewBuffer([]byte(""))
	w := bufio.NewWriter(f)

	p := escpos.New(w)
	p.Init()
	if barcode != "" {
		if code_format == "" {
			code_format = "2"
		}
		barcode = strings.ToUpper(barcode)
		intformat, _ := strconv.Atoi(code_format)
		p.PrintBarcode(barcode, intformat)
	}
	textmap := make(map[string]string)
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			textmap[key] = value
		}
	}
	if print != "" {
		p.Text(textmap, print)
	} else if image != "" {
		downloadFile("download.png", "http://localhost:8000"+image)
		cmd := exec.Command("convert", "download.png", "-resize", "500x500", "downloadr.png")
		cmd.Run()
		printcmd := exec.Command("png2pos -c -p downloadr.png >> " + printmachine)
		printcmd.Run()
	}
	p.End()
	w.Flush()

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

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
