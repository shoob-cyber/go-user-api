package service

import (
	"testing"
	"time"
)

// TestCalculateAge tests the age calculation function
func TestCalculateAge(t *testing.T) {
	service := &UserService{}

	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Person born in 1990, now in 2024, before birthday",
			dob:      time.Date(1990, time.May, 15, 0, 0, 0, 0, time.UTC),
			expected: 33,
		},
		{
			name:     "Person born in 1990, now in 2024, after birthday",
			dob:      time.Date(1990, time.January, 15, 0, 0, 0, 0, time.UTC),
			expected: 34,
		},
		{
			name:     "Person born today (age 0)",
			dob:      time.Now(),
			expected: 0,
		},
		{
			name:     "Person born yesterday (age 0)",
			dob:      time.Now().AddDate(0, 0, -1),
			expected: 0,
		},
		{
			name:     "Person turning age this year",
			dob:      time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateAge(tt.dob)
			if result != tt.expected {
				t.Errorf("CalculateAge(%v) = %d, want %d", tt.dob, result, tt.expected)
			}
		})
	}
}

// TestCalculateAgeEdgeCases tests edge cases
func TestCalculateAgeEdgeCases(t *testing.T) {
	service := &UserService{}

	// Test case: Person born on leap day (Feb 29)
	leapDay := time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC)
	age := service.CalculateAge(leapDay)
	if age < 20 || age > 25 {
		t.Errorf("Leap day age calculation failed: got %d", age)
	}

	// Test case: Very old person
	oldPerson := time.Date(1924, time.June, 15, 0, 0, 0, 0, time.UTC)
	age = service.CalculateAge(oldPerson)
	if age < 95 || age > 105 {
		t.Errorf("Old person age calculation failed: got %d", age)
	}

	// Test case: Person born this year
	thisYear := time.Date(time.Now().Year(), time.February, 15, 0, 0, 0, 0, time.UTC)
	age = service.CalculateAge(thisYear)
	if age != 0 {
		t.Errorf("This year born person should be 0, got %d", age)
	}
}

// BenchmarkCalculateAge benchmarks the age calculation
func BenchmarkCalculateAge(b *testing.B) {
	service := &UserService{}
	dob := time.Date(1990, time.May, 15, 0, 0, 0, 0, time.UTC)

	for i := 0; i < b.N; i++ {
		service.CalculateAge(dob)
	}
}
