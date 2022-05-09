package main

import (
	"fmt"
	"strconv"
)

// Flags: Multiple String Values
type ArrayFlagString []string

func (i *ArrayFlagString) String() string {
	return fmt.Sprint(*i)
}
func (i *ArrayFlagString) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Flags: Multiple Int Values
type ArrayFlagInt []int

func (i *ArrayFlagInt) String() string {
	return fmt.Sprint(*i)
}
func (i *ArrayFlagInt) Set(value string) error {
	cValue, _ := strconv.Atoi(value)
	*i = append(*i, cValue)
	return nil
}
