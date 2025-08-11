package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/droxer/RxGo/pkg/rxgo"
)

// User represents a user entity
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// UserSubscriber processes user data with transformations
type UserSubscriber struct {
	name string
}

func (s *UserSubscriber) Start() {
	fmt.Printf("[%s] Starting user processing\n", s.name)
}

func (s *UserSubscriber) OnNext(user User) {
	formatted := fmt.Sprintf("User: %s (ID: %d, Email: %s, Age: %d)",
		user.Name, user.ID, user.Email, user.Age)
	var category string
	if user.Age < 18 {
		category = "Minor"
	} else if user.Age < 65 {
		category = "Adult"
	} else {
		category = "Senior"
	}

	fmt.Printf("[%s] Processing %s - Category: %s\n", s.name, formatted, category)
}

func (s *UserSubscriber) OnError(err error) {
	fmt.Printf("[%s] Processing error: %v\n", s.name, err)
}

func (s *UserSubscriber) OnCompleted() {
	fmt.Printf("[%s] User processing completed!\n", s.name)
}

// StringProcessor handles string transformations
type StringProcessor struct {
	name string
}

func (s *StringProcessor) Start() {
	fmt.Printf("[%s] Starting string processing\n", s.name)
}

func (s *StringProcessor) OnNext(text string) {
	upper := strings.ToUpper(text)
	charCount := len(text)
	cleaned := strings.ToLower(strings.ReplaceAll(text, " ", ""))
	isPalindrome := cleaned == reverseString(cleaned)

	fmt.Printf("[%s] Input: '%s' -> Upper: '%s', Chars: %d, Palindrome: %t\n",
		s.name, text, upper, charCount, isPalindrome)
}

func (s *StringProcessor) OnError(err error) {
	fmt.Printf("[%s] Processing error: %v\n", s.name, err)
}

func (s *StringProcessor) OnCompleted() {
	fmt.Printf("[%s] String processing completed!\n", s.name)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NumberCruncher processes numbers with mathematical transformations
type NumberCruncher struct {
	name string
}

func (s *NumberCruncher) Start() {
	fmt.Printf("[%s] Starting number crunching\n", s.name)
}

func (s *NumberCruncher) OnNext(num int) {
	square := num * num
	isPrime := num > 1
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			isPrime = false
			break
		}
	}
	factorial := 1
	if num <= 10 {
		for i := 1; i <= num; i++ {
			factorial *= i
		}
	} else {
		factorial = -1
	}

	fmt.Printf("[%s] Input: %d -> Square: %d, Prime: %t, Factorial: %d\n",
		s.name, num, square, isPrime, factorial)
}

func (s *NumberCruncher) OnError(err error) {
	fmt.Printf("[%s] Processing error: %v\n", s.name, err)
}

func (s *NumberCruncher) OnCompleted() {
	fmt.Printf("[%s] Number crunching completed!\n", s.name)
}

func main() {
	fmt.Println("=== Real-World Data Transformation Examples ===")

	fmt.Println("\n1. User Data Processing:")
	users := rxgo.Just(
		User{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 28},
		User{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 35},
		User{ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 42},
		User{ID: 4, Name: "Diana", Email: "diana@example.com", Age: 19},
		User{ID: 5, Name: "Eve", Email: "eve@example.com", Age: 67},
	)
	users.Subscribe(context.Background(), &UserSubscriber{name: "UserProcessor"})

	fmt.Println("\n2. String Processing:")
	strings := rxgo.Just("hello", "world", "level", "racecar", "Go")
	strings.Subscribe(context.Background(), &StringProcessor{name: "TextProcessor"})

	fmt.Println("\n3. Number Analysis:")
	numbers := rxgo.Range(1, 8)
	numbers.Subscribe(context.Background(), &NumberCruncher{name: "MathProcessor"})

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n=== All transformation examples completed ===")
}
