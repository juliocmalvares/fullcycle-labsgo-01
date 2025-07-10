package utils

func ConvertFahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}

func ConvertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func ConvertKelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func FormatTemperatures(celsius float64) map[string]float64 {
	return map[string]float64{
		"temp_C": celsius,
		"temp_F": ConvertCelsiusToFahrenheit(celsius),
		"temp_K": ConvertCelsiusToKelvin(celsius),
	}
}