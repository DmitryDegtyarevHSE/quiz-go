package main

import (
	"bufio"
	"encoding/csv"
	"errors"
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

func parseProblem(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))
	for i := 0; i < len(lines); i++ {
		problems[i] = Problem{
			question: lines[i][0],
			answer:   lines[i][1],
		}
	}
	return problems
}

func problemPuller(fileName string) ([]Problem, error) {
	if fileObj, err := os.Open(fileName); err == nil {
		csvReader := csv.NewReader(fileObj)
		if lines, err := csvReader.ReadAll(); err == nil {
			return parseProblem(lines), nil
		} else {
			return nil, errors.New("error with parsing csv")
		}
	} else {
		return nil, errors.New("error with opening file")
	}
}

func main() {
	const (
		fName   = "quiz.csv"
		seconds = 30
	)
	timer := flag.Int("t", seconds, "timer for the quiz")
	problems, err := problemPuller(fName)
	if err != nil {
		fmt.Printf("something went wrong: %s", err)
		os.Exit(-1)
	}

	ansChan := make(chan string)
	correctAns := 0
	timeObj := time.NewTimer(time.Duration(*timer) * time.Second)
	fmt.Printf("You have %d seconds!\n", seconds)
problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem %d: %s", i+1, problem.question)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			ansChan <- answer
		}()

		select {
		case <-timeObj.C:
			fmt.Println("\nTime is over!")
			close(ansChan)
			break problemLoop
		case iAns := <-ansChan:
			if strings.TrimSpace(iAns) == problem.answer {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansChan)
			}
		}
	}
	<-ansChan
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
}
