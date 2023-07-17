package atcoder_search

func RateToColor(rate int) string {
	if 0 <= rate && rate < 400 {
		return "gray"
	} else if 400 <= rate && rate < 800 {
		return "brown"
	} else if 800 <= rate && rate < 1200 {
		return "green"
	} else if 1200 <= rate && rate < 1600 {
		return "cyan"
	} else if 1600 <= rate && rate < 2000 {
		return "blue"
	} else if 2000 <= rate && rate < 2400 {
		return "yellow"
	} else if 2400 <= rate && rate < 2800 {
		return "orange"
	} else if 2800 <= rate && rate < 3200 {
		return "red"
	} else if 3200 <= rate && rate < 3600 {
		return "silver"
	} else {
		return "gold"
	}
}
