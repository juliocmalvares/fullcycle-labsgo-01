package utils

import (
	"testing"
)

func TestConvertFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		name       string
		fahrenheit float64
		expected   float64
	}{
		{
			name:       "32°F to 0°C",
			fahrenheit: 32.0,
			expected:   0.0,
		},
		{
			name:       "212°F to 100°C",
			fahrenheit: 212.0,
			expected:   100.0,
		},
		{
			name:       "98.6°F to 37°C",
			fahrenheit: 98.6,
			expected:   37.0,
		},
		{
			name:       "0°F to -17.78°C",
			fahrenheit: 0.0,
			expected:   -17.77777777777778,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertFahrenheitToCelsius(tt.fahrenheit)
			if result != tt.expected {
				t.Errorf("ConvertFahrenheitToCelsius(%f) = %f, expected %f", tt.fahrenheit, result, tt.expected)
			}
		})
	}
}

func TestConvertCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C to 32°F",
			celsius:  0.0,
			expected: 32.0,
		},
		{
			name:     "100°C to 212°F",
			celsius:  100.0,
			expected: 212.0,
		},
		{
			name:     "37°C to 98.6°F",
			celsius:  37.0,
			expected: 98.6,
		},
		{
			name:     "-17.78°C to 0°F",
			celsius:  -17.77777777777778,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCelsiusToFahrenheit(tt.celsius)
			if result != tt.expected {
				t.Errorf("ConvertCelsiusToFahrenheit(%f) = %f, expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestConvertCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C to 273.15K",
			celsius:  0.0,
			expected: 273.15,
		},
		{
			name:     "100°C to 373.15K",
			celsius:  100.0,
			expected: 373.15,
		},
		{
			name:     "-273.15°C to 0K",
			celsius:  -273.15,
			expected: 0.0,
		},
		{
			name:     "25°C to 298.15K",
			celsius:  25.0,
			expected: 298.15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCelsiusToKelvin(tt.celsius)
			if result != tt.expected {
				t.Errorf("ConvertCelsiusToKelvin(%f) = %f, expected %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestConvertKelvinToCelsius(t *testing.T) {
	tests := []struct {
		name     string
		kelvin   float64
		expected float64
	}{
		{
			name:     "273.15K to 0°C",
			kelvin:   273.15,
			expected: 0.0,
		},
		{
			name:     "373.15K to 100°C",
			kelvin:   373.15,
			expected: 100.0,
		},
		{
			name:     "0K to -273.15°C",
			kelvin:   0.0,
			expected: -273.15,
		},
		{
			name:     "298.15K to 25°C",
			kelvin:   298.15,
			expected: 25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertKelvinToCelsius(tt.kelvin)
			if result != tt.expected {
				t.Errorf("ConvertKelvinToCelsius(%f) = %f, expected %f", tt.kelvin, result, tt.expected)
			}
		})
	}
}

func TestFormatTemperatures(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected map[string]float64
	}{
		{
			name:    "25°C",
			celsius: 25.0,
			expected: map[string]float64{
				"temp_C": 25.0,
				"temp_F": 77.0,
				"temp_K": 298.15,
			},
		},
		{
			name:    "0°C",
			celsius: 0.0,
			expected: map[string]float64{
				"temp_C": 0.0,
				"temp_F": 32.0,
				"temp_K": 273.15,
			},
		},
		{
			name:    "100°C",
			celsius: 100.0,
			expected: map[string]float64{
				"temp_C": 100.0,
				"temp_F": 212.0,
				"temp_K": 373.15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTemperatures(tt.celsius)

			if result["temp_C"] != tt.expected["temp_C"] {
				t.Errorf("FormatTemperatures(%f) temp_C = %f, expected %f", tt.celsius, result["temp_C"], tt.expected["temp_C"])
			}
			if result["temp_F"] != tt.expected["temp_F"] {
				t.Errorf("FormatTemperatures(%f) temp_F = %f, expected %f", tt.celsius, result["temp_F"], tt.expected["temp_F"])
			}
			if result["temp_K"] != tt.expected["temp_K"] {
				t.Errorf("FormatTemperatures(%f) temp_K = %f, expected %f", tt.celsius, result["temp_K"], tt.expected["temp_K"])
			}
		})
	}
}
