package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
)

// APIResponse represents a mock API response
type APIResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
	Status string `json:"status"`
}

// APISubscriber processes HTTP streaming data
type APISubscriber struct {
	name string
}

func (s *APISubscriber) Start() {
	fmt.Printf("[%s] Starting HTTP data processing\n", s.name)
}

func (s *APISubscriber) OnNext(response APIResponse) {
	shortTitle := response.Title
	if len(shortTitle) > 30 {
		shortTitle = shortTitle[:27] + "..."
	}
	category := "Unknown"
	switch response.Status {
	case "published":
		category = "Live"
	case "draft":
		category = "In Progress"
	case "archived":
		category = "Historical"
	}

	fmt.Printf("[%s] Post %d: %s [%s] - User: %d\n",
		s.name, response.ID, shortTitle, category, response.UserID)
}

func (s *APISubscriber) OnError(err error) {
	fmt.Printf("[%s] API processing error: %v\n", s.name, err)
}

func (s *APISubscriber) OnCompleted() {
	fmt.Printf("[%s] HTTP streaming completed!\n", s.name)
}

// Mock API data generator
func mockAPIStream(ctx context.Context) []APIResponse {
	return []APIResponse{
		{ID: 1, Title: "Introduction to Reactive Programming", Body: "Learn the basics", UserID: 101, Status: "published"},
		{ID: 2, Title: "Advanced Go Patterns", Body: "Deep dive into Go", UserID: 102, Status: "draft"},
		{ID: 3, Title: "Building Microservices", Body: "Architecture guide", UserID: 103, Status: "published"},
		{ID: 4, Title: "Legacy System Migration", Body: "Modernization tips", UserID: 101, Status: "archived"},
		{ID: 5, Title: "Performance Optimization Techniques for High-Scale Applications", Body: "Optimization strategies", UserID: 104, Status: "published"},
	}
}

func RunHTTPStreamingExamples() {
	httpStreamingExamples()
}

func main() {
	httpStreamingExamples()
}

func httpStreamingExamples() {
	fmt.Println("=== HTTP Streaming Example ===")

	apiStream := observable.Create(func(ctx context.Context, sub observable.Subscriber[APIResponse]) {
		data := mockAPIStream(ctx)

		for _, item := range data {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(item)
				time.Sleep(500 * time.Millisecond)
			}
		}
		sub.OnCompleted()
	})

	apiStream.Subscribe(context.Background(), &APISubscriber{name: "HTTPProcessor"})

	time.Sleep(3 * time.Second)
	fmt.Println("\n=== HTTP streaming completed ===")
}
