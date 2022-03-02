package quickchartgo

import (
	"net/url"
	"testing"
)

func TestQuickChart_GetUrl(t *testing.T) {
	type fields struct {
		Width             int64
		Height            int64
		DevicePixelRation float64
		Format            string
		BackgroundColor   string
		Key               string
		Version           string
		Config            string
		Scheme            string
		Host              string
		Port              int64
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test #1",
			fields: fields{
				Width:             500,
				Height:            300,
				DevicePixelRation: 1.0,
				Format:            "png",
				BackgroundColor:   "transparent",
				Key:               "",
				Version:           "",
				Config:            `{type:'bar',data:{labels:['Q1','Q2','Q3','Q4'],datasets:[{label:'Users',data:[50,60,70,180]}]}}`,
				Scheme:            "https",
				Host:              "quickchart.io",
				Port:              443,
			},
			want:    `https://quickchart.io:443/chart?w=500&h=300&devicePixelRatio=1.000000&f=png&bkg=transparent&c={type:'bar',data:{labels:['Q1','Q2','Q3','Q4'],datasets:[{label:'Users',data:[50,60,70,180]}]}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qc := &Chart{
				Width:             tt.fields.Width,
				Height:            tt.fields.Height,
				DevicePixelRation: tt.fields.DevicePixelRation,
				Format:            tt.fields.Format,
				BackgroundColor:   tt.fields.BackgroundColor,
				Key:               tt.fields.Key,
				Version:           tt.fields.Version,
				Config:            tt.fields.Config,
				Scheme:            tt.fields.Scheme,
				Host:              tt.fields.Host,
				Port:              tt.fields.Port,
			}
			got, err := qc.GetUrl()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuickChart.GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotUnescaped, err := url.PathUnescape(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuickChart.GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUnescaped != tt.want {
				t.Errorf("QuickChart.GetUrl() = %v, want %v", gotUnescaped, tt.want)
			}
		})
	}
}
