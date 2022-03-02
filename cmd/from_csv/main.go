package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	quickchartgo "github.com/henomis/quickchart-go"
)

func main() {

	// chartConfig := `{
	// 	type: 'bar',
	// 	data: {
	// 		labels: ['Q1', 'Q2', 'Q3', 'Q4'],
	// 		datasets: [{
	// 		label: 'Users',
	// 		data: [50, 60, 70, 180]
	// 		}]
	// 	}
	// }`

	file, _ := os.Open("data.csv")
	output, _ := buildJSONDataFromCSV(file)

	chartConfig := fmt.Sprintf("{type:'line',%s}", output)

	qc := quickchartgo.New()
	qc.Config = chartConfig

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	err = qc.Write(file)
	if err != nil {
		panic(err)
	}

}

func buildJSONDataFromCSV(input io.Reader) (string, error) {

	// data: {
	// 	labels: ['Q1', 'Q2', 'Q3', 'Q4'],
	// 	datasets: [{
	// 	label: 'Users',
	// 	data: [50, 60, 70, 180]
	// 	}]
	// }

	records, err := importCSV(input)
	if err != nil {
		return "", err
	}

	labels := []string{}
	datasetValues := []string{}
	datasetLabel := ""

	for i, record := range records {

		if i == 0 {
			datasetLabel = record[1]
			continue
		}

		labels = append(labels, record[0])
		datasetValues = append(datasetValues, record[1])

	}

	var sb strings.Builder

	sb.WriteString("data:{labels:[")
	for i, v := range labels {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("'%s'", v))

	}
	sb.WriteString(fmt.Sprintf("],datasets:[{label:'%s',data:[", datasetLabel))
	for i, v := range datasetValues {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("'%s'", v))

	}
	sb.WriteString("]}]}")

	return sb.String(), nil

}

func importCSV(input io.Reader) ([][]string, error) {

	csvStream := csv.NewReader(input)

	records, err := csvStream.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
