package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	/*
	   builtin.Error
	   fmt.Stringer
	   io.Reader
	   io.Writer
	   io.ReadWriteCloser
	   http.ResponseWriter
	   http.Handler
	*/

	//runBuiltinErrorExample()

	//runFmtStringerExample()
	//bufferExample()
	ioReaderWriterOwnImplementation()
}

type someType struct {
	b string
}

func (value someType) String() string {
	return value.b
}

func (value someType) Read(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func runBuiltinErrorExample() {
	//builtin.Error
	//jak tworze jakis error to implementuje interfejs builtin.Error - tj. mam metode error() string
	someErr := errors.New("test")
	fmt.Println(someErr)
}

func runFmtStringerExample() {

	// fmt.Stringer
	//jak tworze jakis typ ktory ma metode string to to implementuje interfejs fmt.Stringer - przez co np. fmt.Print
	//druknie mi zamiast defaultowej implementacji to to co w moim String() string
	//uwaga: nalezy pamietac ze String powinien byc na value (a nie pointer) receiverze
	//bo jak bedzie na pointer to wtedy jak sie da fmt.Print(x) (zamiast fmt.Print(&x)) to i tak da defaultowy print

	fmt.Println(someType{b: "test2"})
}

// https://programmer.group/writer-and-reader-in-go-language.html
func bufferExample() {

	//Define zero-valued Buffer type variable b
	var b bytes.Buffer
	//Write string using Write method
	b.Write([]byte("Hello"))
	//This is to stitch a string into Buffer
	fmt.Fprint(&b, ",", "http://www.flysnow.org")
	//Print Buffer to Terminal Console
	b.WriteTo(os.Stdout)
	data, err := ioutil.ReadAll(&b)
	fmt.Println(string(data), err)
	fmt.Println()
}

type ioReaderWriter struct {
	storage string
}

func (rw *ioReaderWriter) Writer(p []byte) (int, error) {
	rw.storage = string(p)
	return len(rw.storage), nil
}
func (rw *ioReaderWriter) Reader(p []byte) (int, error) {
	//Reader interface implementation is tricky
	//append on p would modify only copy of slice passed here ---> _, _ = rw.Reader(by)
	//in order to see modified original slice underlying array has to be modified
	//it can be done as below
	for i, v := range rw.storage {
		p[i] = byte(v)
	}
	return len(rw.storage), nil
}
func ioReaderWriterOwnImplementation() {
	//Read reads smth and stores it in []byte p
	//Write reads []byte p and do smth to this data

	rw := ioReaderWriter{}
	_, _ = rw.Writer([]byte("some bytes"))
	by := make([]byte, 10)
	_, _ = rw.Reader(by)

	fmt.Println(string(by))

}
