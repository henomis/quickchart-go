package quickchartgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Chart struct {
	Width             int64   `json:"width"`
	Height            int64   `json:"height"`
	DevicePixelRation float64 `json:"devicePixelRatio"`
	Format            string  `json:"format"`
	BackgroundColor   string  `json:"backgroundColor"`
	Key               string  `json:"key"`
	Version           string  `json:"version,omitempty"`
	Config            string  `json:"chart"`

	Scheme  string        `json:"-"`
	Host    string        `json:"-"`
	Port    int64         `json:"-"`
	Timeout time.Duration `json:"-"`
}

type getShortURLResponse struct {
	Success bool   `json:"-"`
	URL     string `json:"url"`
}

func New() *Chart {
	return &Chart{
		Width:             500,
		Height:            300,
		DevicePixelRation: 1.0,
		Format:            "png",
		BackgroundColor:   "#ffffff",

		Scheme:  "https",
		Host:    "quickchart.io",
		Port:    443,
		Timeout: 10 * time.Second,
	}
}

func (qc *Chart) GetUrl() (string, error) {

	if !qc.validateConfig() {
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

func (qc *Chart) GetShortUrl() (string, error) {

	if !qc.validateConfig() {
		return "", fmt.Errorf("invalid config")
	}

	quickChartURL := fmt.Sprintf("%s://%s:%d/chart/create", qc.Scheme, qc.Host, qc.Port)
	bodyStream, err := qc.makePostRequest(quickChartURL)
	if err != nil {
		return "", fmt.Errorf("makePostRequest(%s): %w", quickChartURL, err)
	}

	defer bodyStream.Close()
	body, err := ioutil.ReadAll(bodyStream)
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

func (qc *Chart) Write(output io.Writer) error {

	if !qc.validateConfig() {
		return fmt.Errorf("invalid config")
	}

	quickChartURL := fmt.Sprintf("%s://%s:%d/chart", qc.Scheme, qc.Host, qc.Port)
	bodyStream, err := qc.makePostRequest(quickChartURL)
	if err != nil {
		return fmt.Errorf("makePostRequest(%s): %w", quickChartURL, err)
	}

	defer bodyStream.Close()
	_, err = io.Copy(output, bodyStream)

	return err
}

func (qc *Chart) makePostRequest(endpoint string) (io.ReadCloser, error) {

	jsonEncodedPayload, err := json.Marshal(qc)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: qc.Timeout,
	}

	resp, err := httpClient.Post(
		endpoint,
		"application/json",
		bytes.NewBuffer(jsonEncodedPayload),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error: %d - %s", resp.StatusCode, resp.Header.Get("X-quickchart-error"))
	}

	return resp.Body, nil
}

func (qc *Chart) validateConfig() bool {

	return len(qc.Config) != 0

}
