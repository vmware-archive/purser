package main

import (
	"time"
)

func getMonthToDateCostForInstanceType(instanceType string) (float64, float64, float64) {
	hours := totalHoursTillNow()
	basePrice := getPriceForInstanceType(instanceType)
	totalPrice := hours * basePrice
	cpuMemoryRatio := cpuMemoryRatio(instanceType)
	return totalPrice, totalPrice * cpuMemoryRatio, totalPrice * (1 - cpuMemoryRatio)
}

func getMonthToDateCostForStorageClass(storageClass string) float64 {
	percentageOfHoursElapsed := percentageOfHoursElapsedInCurrentMonth()
	basePrice := getPriceForVolumeType(storageClass)
	return basePrice * percentageOfHoursElapsed
}

func percentageOfHoursElapsedInCurrentMonth() float64 {
	now := time.Now()
	hoursTillNow := 24.0*float64(now.Day()-1) + float64(now.Hour())
	totalDays := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Add(-time.Nanosecond).Day()
	return hoursTillNow / (float64)(totalDays)
}

func totalHoursTillNow() float64 {
	now := time.Now()
	return 24.0*float64(now.Day()-1) + float64(now.Hour())
}

func cpuMemoryRatio(instanceType string) float64 {
	// TODO: enhance this.
	return 0.3
}
