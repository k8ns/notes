package main

import ginlambda "github.com/ksopin/notes/internal/lambda"

func main() {
	err := ginlambda.Run()
	if err != nil {
		panic(err)
	}
}
