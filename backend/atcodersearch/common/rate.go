package common

func RateToColor(rate int) string {
	if rate < 0 {
		return "black"
	} else if rate < 400 {
		return "gray"
	} else if rate < 800 {
		return "brown"
	} else if rate < 1200 {
		return "green"
	} else if rate < 1600 {
		return "cyan"
	} else if rate < 2000 {
		return "blue"
	} else if rate < 2400 {
		return "yellow"
	} else if rate < 2800 {
		return "orange"
	} else if rate < 3200 {
		return "red"
	} else if rate < 3600 {
		return "silver"
	} else {
		return "gold"
	}
}
