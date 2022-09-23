package main

import (
	"fmt"
	"log"
)

type TestTypeString string
type TestTypeMoreComplex moreComplex

type TestAliasString = string

func (t TestTypeString) f() {
	log.Println(t)
}

type moreComplex struct {
	a string
	b int
}

func main() {

	// 	t := TestTypeString("przekazuje stringa..ale nie do funkcji..")
	// to jest wywolanie kontruktora nowego typu

	// to chyba jest tak ze jak zdefiniuje typ inny niz struct
	// (czyli de facto stworze nowy typ na podstawie innego) to jest jakis defaultowy konstruktor
	// ktory sie wywoluje przez operator'()' i przkeazanie w nim typu

	// po co taka operacja w ogole mozliwa.. bo moze ten inny typ zwracaja mi metody
	// i on nie implementuje jakiegos interfejsu z ktorego chce skorzystac
	// a dla tego nowego typu moge zdefiniowac interfejs a potem go zaimplementowac korzystajac z typu ktory mam pod spodem

	// to nie jest asercja typu.
	// asercja by byla jakbym mial obiekt jakiegos interfejsu i chcial go zrzutowac na typ
	// poAsercjiTypu, czySiePowiodloCzyNie := obiektTypuInterfejsowego.(typ - np. implementujacy interfejs)
	//przyklad ponziej
	var a interface{} = "a"
	b, ok := a.(string)
	if !ok {
		log.Fatal("not succesfull type assertion")
	}
	fmt.Println(b)

	// to nie jest tez alias
	// bo jest tak
	// type TestTypeString string
	// a nie tak
	// type TestTypeString = string
	// alias nie tworzy nowego typu
	var alias TestAliasString = "string as alias"
	fmt.Println(alias)
	fmt.Printf("alias type: %T i.e. not type 'defined' by alias\n", alias)

	example1 := TestTypeString("underlying type object")
	fmt.Printf("t type: %T\n", example1)
	log.Println(example1)
	example1.f()

	example2 := TestTypeMoreComplex(TestTypeMoreComplex{
		a: "a",
		b: 2,
	})
	fmt.Printf("%+v", example2)
}
