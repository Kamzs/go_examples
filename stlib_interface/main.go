package main

import (
	"errors"
	"fmt"
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

	//jak tworze jakis error to implementuje interfejs builtin.Error - tj. mam metode error() string
	someErr := errors.New("test")
	fmt.Println(someErr)
	//jak tworze jakis typ ktory ma metode string to to implementuje interfejs fmt.Stringer - przez co np. fmt.Print
	//druknie mi zamiast defaultowej implementacji to to co w moim String() string
	//uwaga: nalezy pamietac ze String powinien byc na value (a nie pointer) receiverze
	//bo jak bedzie na pointer to wtedy jak sie da fmt.Print(x) (zamiast fmt.Print(&x)) to i tak da defaultowy print

	fmt.Println(someType{b: "hahaha"})

}

type someType struct {
	b string
}

func (value someType) String() string {
	return value.b
}
