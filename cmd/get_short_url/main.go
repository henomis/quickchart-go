package main

import (
	"fmt"

	quickchartgo "github.com/henomis/quickchart-go"
)

func main() {

	chartConfig := `{
		type: 'bar',
		data: {
			labels: ['Q1', 'Q2', 'Q3', 'Q4'],
			datasets: [{
			label: 'Users',
			data: [50, 60, 70, 180]
			}]
		}
	}`

	qc := quickchartgo.New()
	qc.Config = chartConfig

	quickchartURL, err := qc.GetShortUrl()
	if err != nil {
		panic(err)
	}

	fmt.Println(quickchartURL)

}
