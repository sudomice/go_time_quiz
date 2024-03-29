package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type Question struct {
	question string
	answer   string
}
type quiz struct {
	answered        int
	answeredCorrect int
	questions       []Question
}

var (
	scanner   = bufio.NewScanner(os.Stdin)
	fileName  = flag.String("filename", "problems.csv", "CSV File that conatins quiz questions")
	timeLimit = flag.Int("limit", 30, "Time Limit for each question")
)

func main() {
	f, err := openFile(*fileName)
	if err != nil {
		return
	}

	quiz, _ := readCSV(f)
	quiz.run()
	fmt.Printf("\nAnswered %d correct out of %d answered questions\nTotal Questions %d", quiz.answeredCorrect, quiz.answered, len(quiz.questions))

}
func openFile(fileName string) (io.Reader, error) {
	return os.Open(fileName)
}
func readCSV(f io.Reader) (*quiz, error) {
	allQuestions, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	numQues := len(allQuestions)
	if numQues == 0 {
		return nil, fmt.Errorf("No question in file")
	}
	var quiz quiz
	for _, line := range allQuestions {
		ques := Question{}
		ques.question = line[0]
		ques.answer = line[1]
		quiz.questions = append(quiz.questions, ques)
	}
	return &quiz, nil
}

func (quiz *quiz) run() {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	fmt.Println("asking a q")
quizLoop:
	for _, question := range quiz.questions {
		fmt.Println(question.question)
		ansCh := make(chan string)
		go func() {
			scanner.Scan()
			answer := scanner.Text()
			ansCh <- answer
		}()
		select {
		case <-timer.C:
			break quizLoop
		case answer := <-ansCh:
			if answer == question.answer {
				quiz.answeredCorrect++
			}
			quiz.answered++
		}
	}
	return
}
