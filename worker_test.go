package main

import (
	"testing"
)

func TestFloatToString(t *testing.T) {
	output := FloatToString(2.23)
	if output != "2.23" {
		t.Error("Expected 2.23 but was ", output)
	}
}

func TestFloatToStringWithBiggerPrecision(t *testing.T) {
	output := FloatToString(2.2323)
	if output != "2.23" {
		t.Error("Expected 2.23 but was ", output)
	}
}

func TestIntToString(t *testing.T) {
	output := IntToString(12)
	if output != "12" {
		t.Error("Expected 12 but was ", output)
	}
}

func TestCalculateFareAmountWhenSpeedBiggerThanTenAndNightShift(t *testing.T) {
	output := calculateFareAmount(20, 0.0, 10, true)
	if output != 13 {
		t.Error("Expected 13 but was ", output)
	}
}

func TestCalculateFareAmountWhenSpeedBiggerThanTenAndDayShift(t *testing.T) {
	output := calculateFareAmount(20, 0.0, 10, false)
	if output != 7.4 {
		t.Error("Expected 7.4 but was ", output)
	}
}

func TestCalculateFareAmountWhenIdle(t *testing.T) {
	output := calculateFareAmount(5, 0.5, 10, true)
	if output != 5.95 {
		t.Error("Expected 5.95 but was ", output)
	}
}
