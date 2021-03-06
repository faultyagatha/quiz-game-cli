package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"strings"
	"time"

	"os"
)

//Problem is a shape of our problem
type Problem struct {
	q string
	a string
}

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the format of question-answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("failed to open the file: %s\n", *csvFile))
	}
	//create a reader
	r := csv.NewReader(file)
	//parse the reader
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("failed to parse the file"))
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit * int(time.Second)))
	doQuiz(problems, timer)
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []Problem {
	res := make([]Problem, len(lines))
	for i, v := range lines {
		res[i] = Problem{
			q: v[0],                    //question
			a: strings.TrimSpace(v[1]), //answer with trimmed spaces
		}
	}
	return res
}

func doQuiz(problems []Problem, timer *time.Timer) {
	var answer string
	correct := 0
	for i, v := range problems {
		//start at index 1
		fmt.Printf("Problem #%d: %s = \n", i+1, v.q)
		answerCh := make(chan string)
		//use goroutine to avoid blocking
		go func() {
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer //send the answer to the channel
		}()
		//always wait for an answer from 1 channel or the other
		select {
		case <-timer.C: //if we get a timer back, the result is not valid
			fmt.Printf("\nYou scored %d out of %d.", correct, len(problems))
			return
		case answer = <-answerCh:
			if answer == v.a {
				fmt.Println("Correct")
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.", correct, len(problems))
}
