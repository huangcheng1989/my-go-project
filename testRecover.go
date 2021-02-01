package main

import (
	"errors"
	"fmt"
)

func main() {
	{
		err := release()
		fmt.Println(err)
	}

	{
		err := correctRelease()
		fmt.Println(err)
	}
}

func release() error {
	defer func() error {
		return errors.New("error abc")
	}()

	return nil
}

func correctRelease() (err error) {
	defer func() {
		err = errors.New("error 123")
	}()
	return nil
}
