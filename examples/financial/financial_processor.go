package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/rxgo"
)

// Trade represents a financial transaction
type Trade struct {
	Symbol    string
	Price     float64
	Volume    int
	Timestamp time.Time
}

// Portfolio represents an investment portfolio
type Portfolio struct {
	Holdings map[string]int
	Cash     float64
}

// FinancialProcessor processes trading data
type FinancialProcessor struct {
	portfolio *Portfolio
	name      string
}

func NewFinancialProcessor(name string, initialCash float64) *FinancialProcessor {
	return &FinancialProcessor{
		name: name,
		portfolio: &Portfolio{
			Holdings: make(map[string]int),
			Cash:     initialCash,
		},
	}
}

func (s *FinancialProcessor) Start() {
	fmt.Printf("[%s] Starting financial data processing\n", s.name)
	fmt.Printf("[%s] Initial portfolio: $%.2f cash\n", s.name, s.portfolio.Cash)
}

func (s *FinancialProcessor) OnNext(trade Trade) {
	tradeValue := trade.Price * float64(trade.Volume)
	if _, exists := s.portfolio.Holdings[trade.Symbol]; !exists {
		s.portfolio.Holdings[trade.Symbol] = 0
	}
	s.portfolio.Holdings[trade.Symbol] += trade.Volume
	s.portfolio.Cash -= tradeValue
	
	totalValue := s.portfolio.Cash
	for symbol, shares := range s.portfolio.Holdings {
		if symbol == trade.Symbol {
			totalValue += trade.Price * float64(shares)
		} else {
			totalValue += 100.0 * float64(shares)
		}
	}
	
	var riskLevel string
	if trade.Volume > 1000 {
		riskLevel = "HIGH"
	} else if trade.Volume > 100 {
		riskLevel = "MEDIUM"
	} else {
		riskLevel = "LOW"
	}

	fmt.Printf("[%s] %s: %d shares @ $%.2f = $%.2f [Risk: %s] Portfolio: $%.2f\n",
		s.name, trade.Symbol, trade.Volume, trade.Price, tradeValue, riskLevel, totalValue)
}

func (s *FinancialProcessor) OnError(err error) {
	fmt.Printf("[%s] Financial processing error: %v\n", s.name, err)
}

func (s *FinancialProcessor) OnCompleted() {
	fmt.Printf("[%s] Financial processing completed!\n", s.name)
	fmt.Printf("[%s] Final portfolio:\n", s.name)
	for symbol, shares := range s.portfolio.Holdings {
		fmt.Printf("[%s]   %s: %d shares\n", s.name, symbol, shares)
	}
	fmt.Printf("[%s]   Cash: $%.2f\n", s.name, s.portfolio.Cash)
}

// MarketDataGenerator generates mock trading data
func generateMarketData(ctx context.Context) []Trade {
	return []Trade{
		{Symbol: "AAPL", Price: 175.25, Volume: 150, Timestamp: time.Now()},
		{Symbol: "GOOGL", Price: 2850.50, Volume: 25, Timestamp: time.Now().Add(1 * time.Minute)},
		{Symbol: "MSFT", Price: 340.75, Volume: 200, Timestamp: time.Now().Add(2 * time.Minute)},
		{Symbol: "TSLA", Price: 250.30, Volume: 75, Timestamp: time.Now().Add(3 * time.Minute)},
		{Symbol: "AMZN", Price: 145.80, Volume: 1200, Timestamp: time.Now().Add(4 * time.Minute)},
	}
}

func main() {
	fmt.Println("=== Financial Trading Example ===")

	processor := NewFinancialProcessor("TradeEngine", 10000.0)

	// Simulate real-time trading data
	trades := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[Trade]) {
		data := generateMarketData(ctx)

		for _, trade := range data {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(trade)
				time.Sleep(800 * time.Millisecond)
			}
		}
		sub.OnCompleted()
	})

	trades.Subscribe(context.Background(), processor)

	time.Sleep(5 * time.Second)
	fmt.Println("\n=== Financial trading completed ===")
}
