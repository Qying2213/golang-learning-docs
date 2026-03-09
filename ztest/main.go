package main

import (
	"fmt"
	"testing"
)

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello,world"

	if got != want {
		t.Errorf("got '%q' want '%q' ", got, want)
	} else {
		fmt.Println("test funcHelloWorld pass")
	}
}

func Hello() string {
	return "Hello,world"
}

func main() {
	TestHello(&testing.T{})

}
