package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

var num_correct int

func readArgs() (string, time.Duration) {
	var timeout time.Duration
	filename := flag.String("filename", "problems.csv", "CSV File that conatins quiz questions")
	flag.DurationVar(&timeout, "time", 30*time.Second, "Timeout duration in seconds")

	flag.Parse()
	return *filename, timeout
}

func startTimeout(done chan bool, duration time.Duration) error {
	r := bufio.NewReader(os.Stdin)

	fmt.Println("Press any key to start timer...")
	_, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("error reading input at startTimeout: %v", err)
	}

	go func() {
		time.Sleep(duration)
		done <- true
	}()

	return err
}

func loadCsv(filename string) ([]Problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	defer file.Close()

	var problems []Problem

	for _, row := range records {
		q := strings.TrimSpace(row[0])
		ans := strings.TrimSpace(row[1])
		problems = append(problems, Problem{q, ans})
	}
	return problems, err
}

func runQuiz(problems []Problem) error {
	r := bufio.NewReader(os.Stdin)

	for _, problem := range problems {
		fmt.Printf("%s: ", problem.question)

		input, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}

		input = strings.TrimSpace(input)

		if input != "" && problem.answer == input {
			num_correct++
		}
	}
	return nil
}

func handleExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	filename, timeout := readArgs()
	problems, err := loadCsv(filename)
	handleExit(err)

	done := make(chan bool)
	fake := make(chan bool)

	err = startTimeout(done, timeout)
	handleExit(err)

	go func() {
		err = runQuiz(problems)
		handleExit(err)
	}()

	select {
	case <-done:
		fmt.Println("\n\nTime is up! Exiting...")
		fmt.Printf("\nYou got %d correct, and %d wrong!", num_correct, len(problems)-num_correct)
	case <-fake:
		time.Sleep(1 * time.Millisecond)
	}

	os.Exit(0)
}
