package main

func main() {
	// n := average(1, 2, 3)
	// fmt.Println(n)

	// data := []float64{3, 2, 34}
	// n := average(data...)
	// fmt.Println(n)
}

func average(sf ...float64) float64 {
	var total float64
	for _, item := range sf {
		total += item
	}

	return total / float64(len(sf))
}
