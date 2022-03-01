package main

import (
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

	err := qc.ToFile("pippo.png")
	if err != nil {
		panic(err)
	}

}
