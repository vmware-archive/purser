package main

import (
	"time"
)

func getMonthlyPriceForInstance(instanceType string) (float64, float64, float64) {
	hours := totalHoursTillNow()
	basePrice := getPriceForInstanceType(instanceType)
	totalPrice := hours * basePrice
	cpuMemoryRatio := cpuMemoryRatio(instanceType)
	return totalPrice, totalPrice * cpuMemoryRatio, totalPrice * (1 - cpuMemoryRatio)
}

func totalHoursTillNow() float64 {
	now := time.Now()
	return 24.0*float64(now.Day()-1) + float64(now.Hour())
}

func cpuMemoryRatio(instanceType string) float64 {
	// TODO: enhance this.
	return 0.3
}
