package main

import (
	"github.com/Spankyduck/minyr/yr"
	"bufio"
    "fmt"
    "os"
)	

func main() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Vennligst velg convert, exit eller average.\n")
	for scanner.Scan() {
		input = scanner.Text()
		if input == "q" || input == "exit" {
			fmt.Println("exit")
			os.Exit(0)
		} else if input == "convert" {
			fmt.Println("Konverterer alle maalingene gitt i grader Celsius til grader Fahrenheit.")

		yr.ConvTemperature()
		} else if input == "average" {
		yr.AverageTemp()
		}

	}
}
