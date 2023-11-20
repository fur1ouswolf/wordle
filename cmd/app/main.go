package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"
)

func loadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	words := make([]string, 0, 100)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func selectRandomWord(words *[]string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return (*words)[r.Intn(len(*words))]
}

func validateInput(userInput string, n int, words []string) (string, error) {
	userInput = strings.TrimSuffix(userInput, "\n")
	userInput = strings.ToLower(userInput)
	if len(userInput) != n {
		return "", fmt.Errorf("invalid input")
	}
	for _, letter := range userInput {
		if letter < 'a' || letter > 'z' {
			return "", fmt.Errorf("invalid input")
		}
	}
	if slices.Contains(words, userInput) == false {
		return "", fmt.Errorf("unknown word")
	}
	return userInput, nil
}

func compareWords(word string, userInput string) []int {
	wordRunes := []rune(word)
	userInputRunes := []rune(userInput)
	letters := make([]int, len(wordRunes))
	for i := 0; i < len(wordRunes); i++ {
		if wordRunes[i] == userInputRunes[i] {
			letters[i] = 2
			continue
		}
		if slices.Contains(wordRunes, userInputRunes[i]) {
			letters[i] = 1
		} else {
			letters[i] = 0
		}
	}
	return letters
}

func gameLoop(word string, words []string) {
	n := len(word)
	fmt.Println(strings.Repeat("⬜ ", n))
	reader := bufio.NewReader(os.Stdin)
	for try := 0; try < 6; try++ {
		s, err := reader.ReadString('\n')
		userInput, err := validateInput(s, n, words)
		fmt.Print("\033[F\033[K")
		if err != nil {
			fmt.Println(err)
			try--
			continue
		}
		userLetters := compareWords(word, userInput)
		cnt := 0
		for i := 0; i < n; i++ {
			switch userLetters[i] {
			case 0:
				fmt.Printf("\033[31m%c\033[0m ", rune(userInput[i]))
			case 1:
				fmt.Printf("\033[33m%c\033[0m ", rune(userInput[i]))
			case 2:
				fmt.Printf("\033[32m%c\033[0m ", rune(userInput[i]))
				cnt++
			}
		}
		fmt.Println()
		if cnt == n {
			fmt.Print("\033[32mYou have won!\033[0m")
			return
		}

	}
	fmt.Print("\033[31mYou have lost!\033[0m")
}

func main() {
	words, err := loadWords("validwords.txt")
	if err != nil {
		panic(err)
	}

	word := selectRandomWord(&words)
	gameLoop(word, words)
}
