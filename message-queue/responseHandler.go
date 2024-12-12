package main

import "fmt"

func ResponseHandler(ch string) {
	switch ch {
	case "free\n":
		fmt.Println("Parsing URL")
	case "busy\n":
		fmt.Println("Storaging URL")
	case "inurl\n":
		fmt.Println("Invalid URL")
	}
}
