package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func main() {
	fmt.Println(nomain("1"))
	fmt.Println(nomain(1))
	fmt.Println(nomain(""))
	fmt.Println(nomain(0.1232424452))
	fmt.Println(nomain(0))


}
func nomain [T constraints.](a T) T {
	zero := *new(T) 
	if  a > zero {return a
	} else
	 { return zero }

}