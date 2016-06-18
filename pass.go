package main

import (
	"bufio"
	"github.com/gnanakeethan/print/escpos"
	"os"
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

	print := "\x1dk\x04123123123\x00"
	path := "output"
	// detect if file exists
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		// create file if not exists
		var file, _ = os.Create(path)
		defer file.Close()
	}

	f, err := os.OpenFile("output", O_RDWR, 0644)
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := escpos.New(w)
	p.Init()
	p.WriteRaw([]byte(print))
	p.End()

	w.Flush()
}
