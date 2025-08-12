package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
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

// SchedulerProcessor demonstrates scheduler usage in data processing
func SchedulerProcessor(scheduler observable.Scheduler, name string, data []User) {
	fmt.Printf("\n[%s] Processing with %s scheduler\n", name, scheduler)

	start := time.Now()
	processor := &UserSubscriber{name: name}

	users := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[User]) {
		for _, user := range data {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				scheduler.Schedule(func() {
					sub.OnNext(user)
				})
				time.Sleep(10 * time.Millisecond)
			}
		}
		sub.OnCompleted()
	})

	users.Subscribe(context.Background(), processor)

	if name == "SingleThread" {
		defer scheduler.(*observable.SingleThreadScheduler).Close()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[%s] Completed in %v\n", name, time.Since(start))
}

func main() {
	fmt.Println("=== Real-World Data Transformation Examples ===")

	// Sample data
	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 28},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 35},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 42},
		{ID: 4, Name: "Diana", Email: "diana@example.com", Age: 19},
		{ID: 5, Name: "Eve", Email: "eve@example.com", Age: 67},
	}

	fmt.Println("\n1. User Data Processing:")
	rxgo.Just(users...).Subscribe(context.Background(), &UserSubscriber{name: "UserProcessor"})

	fmt.Println("\n2. String Processing:")
	rxgo.Just("hello", "world", "level", "racecar", "Go").Subscribe(context.Background(), &StringProcessor{name: "TextProcessor"})

	fmt.Println("\n3. Number Analysis:")
	rxgo.Range(1, 8).Subscribe(context.Background(), &NumberCruncher{name: "MathProcessor"})

	fmt.Println("\n4. Scheduler Comparison:")
	schedulers := map[string]observable.Scheduler{
		"Immediate":    observable.NewImmediateScheduler(),
		"NewThread":    observable.NewNewThreadScheduler(),
		"SingleThread": observable.NewSingleThreadScheduler(),
	}

	for name, scheduler := range schedulers {
		SchedulerProcessor(scheduler, name+"-User", users)
	}

	fmt.Println("\n=== All transformation examples completed ===")
}
