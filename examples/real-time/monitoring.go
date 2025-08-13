package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/droxer/RxGo/pkg/rx"
)

// SensorData represents IoT sensor readings
type SensorData struct {
	DeviceID  string
	Value     float64
	Unit      string
	Timestamp time.Time
	Location  string
}

// Alert represents system alerts
type Alert struct {
	Level   string
	Message string
	Time    time.Time
}

// MonitoringService processes real-time sensor data
type MonitoringService struct {
	name          string
	thresholdHigh float64
	thresholdLow  float64
	alertCount    int
	lastAlert     time.Time
}

func NewMonitoringService(name string, low, high float64) *MonitoringService {
	return &MonitoringService{
		name:          name,
		thresholdLow:  low,
		thresholdHigh: high,
	}
}

func (s *MonitoringService) Start() {
	fmt.Printf("[%s] Starting IoT monitoring service\n", s.name)
	fmt.Printf("[%s] Thresholds: Low=%.2f, High=%.2f\n", s.name, s.thresholdLow, s.thresholdHigh)
}

func (s *MonitoringService) OnNext(data SensorData) {
	// Real-time analysis
	status := "NORMAL"
	if data.Value > s.thresholdHigh {
		status = "HIGH"
		s.generateAlert("HIGH", fmt.Sprintf("%s sensor %.2f%s exceeds high threshold", data.DeviceID, data.Value, data.Unit))
	} else if data.Value < s.thresholdLow {
		status = "LOW"
		s.generateAlert("LOW", fmt.Sprintf("%s sensor %.2f%s below low threshold", data.DeviceID, data.Value, data.Unit))
	}

	trend := "stable"
	if randomFloat64() > 0.7 {
		trend = "increasing"
	} else if randomFloat64() > 0.7 {
		trend = "decreasing"
	}

	fmt.Printf("[%s] %s @ %s: %.2f%s [%s] Trend: %s\n",
		s.name, data.DeviceID, data.Location, data.Value, data.Unit, status, trend)
}

func (s *MonitoringService) generateAlert(level, message string) {
	if time.Since(s.lastAlert) < 2*time.Second {
		return
	}

	alert := Alert{
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}

	fmt.Printf("🚨 [%s] ALERT: %s - %s\n", s.name, alert.Level, alert.Message)
	s.alertCount++
	s.lastAlert = time.Now()
}

func (s *MonitoringService) OnError(err error) {
	fmt.Printf("[%s] Monitoring error: %v\n", s.name, err)
}

func (s *MonitoringService) OnCompleted() {
	fmt.Printf("[%s] Monitoring completed! Total alerts: %d\n", s.name, s.alertCount)
}

// randomFloat64 generates a random float64 in [0.0, 1.0) using crypto/rand
func randomFloat64() float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<53))
	if err != nil {
		return 0.5 // fallback value
	}
	return float64(n.Int64()) / float64(1<<53)
}

// SensorSimulator generates realistic IoT sensor data
func sensorSimulator(ctx context.Context, deviceID, location string) []SensorData {
	var data []SensorData
	baseTemp := 22.0

	for i := 0; i < 15; i++ {
		value := baseTemp + (randomFloat64()*10 - 5)

		data = append(data, SensorData{
			DeviceID:  deviceID,
			Value:     value,
			Unit:      "°C",
			Timestamp: time.Now().Add(time.Duration(i*30) * time.Second),
			Location:  location,
		})
	}
	return data
}

func main() {
	fmt.Println("=== Real-Time IoT Monitoring ===")

	// Create monitoring services
	temperatureMonitor := NewMonitoringService("TempMonitor", 18.0, 26.0)

	// Simulate multiple IoT devices
	devices := []struct {
		id       string
		location string
	}{
		{"sensor-001", "Living Room"},
		{"sensor-002", "Kitchen"},
		{"sensor-003", "Bedroom"},
		{"sensor-004", "Office"},
	}

	// Create temperature monitoring stream
	tempStream := rx.Create(func(ctx context.Context, sub rx.Subscriber[SensorData]) {
		for _, device := range devices {
			data := sensorSimulator(ctx, device.id, device.location)
			for _, reading := range data {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(reading)
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
		sub.OnCompleted()
	})

	tempStream.Subscribe(context.Background(), temperatureMonitor)

	time.Sleep(10 * time.Second)
	fmt.Println("\n=== Real-time monitoring completed ===")
}

// Run with: go run examples/real-time/monitoring.go
