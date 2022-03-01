package quickchartgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type QuickChart struct {
	Width             int64   `json:"width"`
	Height            int64   `json:"height"`
	DevicePixelRation float64 `json:"devicePixelRatio"`
	Format            string  `json:"format"`
	BackgroundColor   string  `json:"backgroundColor"`
	Key               string  `json:"key"`
	Version           string  `json:"version"`
	Config            string  `json:"chart"`

	Scheme string `json:"-"`
	Host   string `json:"-"`
	Port   int64  `json:"-"`
}

type getShortURLResponse struct {
	Success bool   `json:"-"`
	URL     string `json:"url`
}

func New() *QuickChart {
	return &QuickChart{
		Width:             500,
		Height:            300,
		DevicePixelRation: 1.0,
		Format:            "png",
		BackgroundColor:   "transparent",

		Scheme: "https",
		Host:   "quickchart.io",
		Port:   443,
	}
}

func (qc *QuickChart) GetUrl() (string, error) {

	if len(qc.Config) == 0 {
		return "", fmt.Errorf("invalid config")
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("w=%d", qc.Width))
	sb.WriteString(fmt.Sprintf("&h=%d", qc.Height))
	sb.WriteString(fmt.Sprintf("&devicePixelRatio=%f", qc.DevicePixelRation))
	sb.WriteString(fmt.Sprintf("&f=%s", qc.Format))
	sb.WriteString(fmt.Sprintf("&bkg=%s", url.QueryEscape(qc.BackgroundColor)))
	sb.WriteString(fmt.Sprintf("&c=%s", url.QueryEscape(qc.Config)))

	if len(qc.Key) > 0 {
		sb.WriteString(fmt.Sprintf("&key=%s", url.QueryEscape(qc.Key)))
	}

	if len(qc.Version) > 0 {
		sb.WriteString(fmt.Sprintf("&v=%s", url.QueryEscape(qc.Key)))
	}

	return fmt.Sprintf("%s://%s:%d/chart?%s", qc.Scheme, qc.Host, qc.Port, sb.String()), nil

}

func (qc *QuickChart) GetShortUrl() (string, error) {

	if len(qc.Config) == 0 {
		return "", fmt.Errorf("invalid config")
	}

	quickChartURL := fmt.Sprintf("%s://%s:%d/chart/create", qc.Scheme, qc.Host, qc.Port)

	jsonEncodedPayload, err := json.Marshal(qc)
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Post(
		quickChartURL,
		"application/json",
		bytes.NewBuffer(jsonEncodedPayload),
	)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response error: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	unescapedResponse, err := url.PathUnescape(string(body))
	if err != nil {
		return "", err
	}

	decodedResponse := &getShortURLResponse{}
	err = json.Unmarshal([]byte(unescapedResponse), decodedResponse)
	if err != nil {
		return "", err
	}

	return decodedResponse.URL, nil

}

func (qc *QuickChart) ToByteArray() ([]byte, error) {
	if len(qc.Config) == 0 {
		return nil, fmt.Errorf("invalid config")
	}

	quickChartURL := fmt.Sprintf("%s://%s:%d/chart", qc.Scheme, qc.Host, qc.Port)

	jsonEncodedPayload, err := json.Marshal(qc)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Post(
		quickChartURL,
		"application/json",
		bytes.NewBuffer(jsonEncodedPayload),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (qc *QuickChart) ToFile(filePath string) error {

	rawFile, err := qc.ToByteArray()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, rawFile, 0660)

}
