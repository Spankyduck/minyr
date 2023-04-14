package yr

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// Convert Celsius to Fahrenheit
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func Convert() error {
	// sjekker hvis output filen finnes
	if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); !os.IsNotExist(err) {
		// Output file finnes allerede, prompt user to regenerate
		var regenerate string
		fmt.Print("Output file already exists. Regenerate? (y/n): ")
		fmt.Scanln(&regenerate)
		if regenerate != "y" && regenerate != "Y" {
			fmt.Println("Exiting without generating new file.")
			return nil
		}
	}

	// åpner input CSV file
	inputFile, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
	}
	defer inputFile.Close()

	// CLager en ny scanner = input CSV file
	inputScanner := bufio.NewScanner(inputFile)

	// lager en ny CSV writer som skriver til output csv file
	outputFile, err := os.Create("kjevik-temp-fahr-20220318-20230318.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
	}
	defer outputFile.Close()

	outputWriter := csv.NewWriter(outputFile)
	defer outputWriter.Flush()

	// Prints out the first line of the input CSV file
	if inputScanner.Scan() {
		firstLine := inputScanner.Text()
		if err = outputWriter.Write(strings.Split(firstLine, ";")); err != nil {
			fmt.Println("Error writing first line:", err)
		}
	}

	// Loop through each line of the input CSV file
	lineNo := 2 // Starter på linje 2 siden linje 1 allerede er skrevet
	for inputScanner.Scan() {
		//skjekker om linje nummer overstiger 16755 og bryter ut av loopen hvis det gjør det
		if lineNo > 16755 {
			break
		}
		// Split the line into fields
		fields := strings.Split(inputScanner.Text(), ";")

		// Check that the fields slice has at least 4 elements
		if len(fields) != 4 {
			fmt.Printf("Error on line %d: Invalid input format.\n", lineNo)
			continue
		}

		// Extract the last digit from the fourth column
		temperatureField := fields[3]
		if temperatureField == "" {
			fmt.Printf("wrong line %d: Tempraturfielt er tomt.\n", lineNo)
			continue
		}
		temperature, err := strconv.ParseFloat(temperatureField, 64)
		if err != nil {
			fmt.Printf("wrong line %d: %v\n", lineNo, err)
			continue
		}
		if math.IsNaN(temperature) {
			fmt.Printf("Feil på linje %d: Temperatur er ikke gyldig float64 verdi.\n", lineNo)
			continue
		}
		lastDigit := temperature - float64(int(temperature/10))*10

		// Convert Celsius to Fahrenheit
		/*fahrenheit := conv.CelsiusToFarenheit(lastDigit)*/
		fahrenheit := celsiusToFahrenheit(lastDigit)

		// Replace the temperature in the fourth column with the converted value
		temperatureString := strconv.FormatFloat(fahrenheit, 'f', 2, 64)
		temperatureParts := strings.Split(temperatureString, ".")
		fields[3] = temperatureParts[0] + "." + string(temperatureParts[1][0])

		// Write the updated line to the output CSV file
		fields[3] = fmt.Sprintf("%.1f", fahrenheit)
		if err = outputWriter.Write(fields); err != nil {
			fmt.Println("Error writing line to output file:", err)
			return err
		}

		lineNo++
	}

	dataText := "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av Sebastian Garcia Jessen"
	err = outputWriter.Write([]string{dataText})
	if err != nil {
		fmt.Println("Error writing data text to output file:", err)

	}

	return nil
}

func Average(unit string) (float64, error) {
	var filename string
	var tempColumn int
	var delimeter rune

	if unit == "c" {
		filename = "kjevik-temp-celsius-20220318-20230318.csv"
		tempColumn = 3
		delimeter = ';'
	} else if unit == "f" {
		filename = "kjevik-temp-fahr-20220318-20230318.csv"
		tempColumn = 3
		delimeter = ','
	} else {
		return 0, fmt.Errorf("invalid temperature unit: %s", unit)
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimeter

	var total float64
	var count int

	// Loop through the lines in the CSV file
	for i := 1; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		if i < 2 || i > 16755 {
			//hopper over linjer utenfor rangen
			continue
		}

		if len(record) <= tempColumn {
			return 0, fmt.Errorf("invalid data in file %s", filename)
		}

		temp, err := strconv.ParseFloat(record[tempColumn], 64)
		if err != nil {
			return 0, err
		}

		total += temp
		count++
	}

	if count == 0 {
		return 0, fmt.Errorf("no temperature data found in file %s", filename)
	}

	return total / float64(count), nil
}
