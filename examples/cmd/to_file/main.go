package main

import (
	"os"

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

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	qc.Write(file)

}
