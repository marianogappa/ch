package chartjs

import (
	"bytes"
	"io"
	"log"
	"time"
)

// ChartJS is the main struct for building Chart.js charts
type ChartJS struct {
	chart    *Chart
	canvasID string
}

// New creates a new ChartJS instance
func New(chart *Chart) *ChartJS {
	return &ChartJS{
		chart:    chart,
		canvasID: "chart",
	}
}

// SetCanvasID sets the canvas element ID (default: "chart")
func (c *ChartJS) SetCanvasID(id string) {
	c.canvasID = id
}

// OutputMode parameterizes the output of Build()
type OutputMode int

const (
	OutputAll          OutputMode = iota // Outputs a complete HTML file
	OutputHTMLHeader                     // Outputs an HTML header, dependencies and opening <body> tag
	OutputDependencies                   // Outputs only <script> tags with Chart.js dependencies
	OutputChart                          // Outputs the JS object that Chart.js expects: new Chart(ctx, {{.}})
	OutputHTMLFooter                     // Outputs </body></html>
)

// MustBuild prepares the chart and executes the text template with it. Fatals if there's a problem
func (c *ChartJS) MustBuild(om OutputMode, w io.Writer) {
	if err := c.Build(om, w); err != nil {
		log.Fatal(err)
	}
}

// Build prepares the chart and executes the text template with it. Returns an error if there's a problem
func (c *ChartJS) Build(om OutputMode, w io.Writer) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in ChartJS.Build: %v", r)
		}
	}()

	switch om {
	case OutputDependencies:
		return tplToWriter(tplChartJS, "", w)
	case OutputHTMLHeader:
		bb := bytes.Buffer{}
		if err := tplToWriter(tplChartJS, "", &bb); err != nil {
			return err
		}
		return tplToWriter(tplHTMLHeader, bb.String(), w)
	case OutputChart:
		chartObj, err := c.prepareChartObject()
		if err != nil {
			return err
		}
		return tplToWriter(tplChartObject, chartObj, w)
	case OutputHTMLFooter:
		return tplToWriter(tplHTMLFooter, "", w)
	case OutputAll:
		deps := bytes.Buffer{}
		if err := tplToWriter(tplChartJS, "", &deps); err != nil {
			return err
		}
		if err := tplToWriter(tplHTMLHeader, deps.String(), w); err != nil {
			return err
		}

		chartObj, err := c.prepareChartObject()
		if err != nil {
			return err
		}

		chartScript := bytes.Buffer{}
		if err := tplToWriter(tplChartObject, chartObj, &chartScript); err != nil {
			return err
		}

		return tplToWriter(tplChartDivScript, map[string]string{
			"CanvasID":    c.canvasID,
			"ChartConfig": chartScript.String(),
		}, w)
	}
	return nil
}

// templateData holds data for template execution
type templateData struct {
	CanvasID    string
	ChartConfig string
}

func (c *ChartJS) prepareChartObject() (string, error) {
	// Use custom marshaler that handles time.Time, functions, etc.
	return marshalChartToJS(c.chart)
}

// Helper functions for building charts from common data structures

// NewLineChart creates a line chart from simple data
func NewLineChart(labels []string, datasets []LineDataset) *Chart {
	data := &Data{
		Labels:   stringSliceToAny(labels),
		Datasets: make([]Dataset, len(datasets)),
	}

	for i, ds := range datasets {
		data.Datasets[i] = Dataset{
			Label:           ds.Label,
			Data:            floatSliceToAny(ds.Data),
			BorderColor:     ds.BorderColor,
			BackgroundColor: ds.BackgroundColor,
			Fill:            boolPtr(ds.Fill),
			Tension:         ds.Tension,
		}
	}

	return &Chart{
		Type: ChartTypeLine,
		Data: data,
	}
}

// NewBarChart creates a bar chart from simple data
func NewBarChart(labels []string, datasets []BarDataset) *Chart {
	data := &Data{
		Labels:   stringSliceToAny(labels),
		Datasets: make([]Dataset, len(datasets)),
	}

	for i, ds := range datasets {
		data.Datasets[i] = Dataset{
			Label:           ds.Label,
			Data:            floatSliceToAny(ds.Data),
			BackgroundColor: ds.BackgroundColor,
			BorderColor:     ds.BorderColor,
			BorderWidth:     ds.BorderWidth,
		}
	}

	return &Chart{
		Type: ChartTypeBar,
		Data: data,
	}
}

// NewPieChart creates a pie chart from simple data
func NewPieChart(labels []string, data []float64, colors []string) *Chart {
	datasets := []Dataset{{
		Label:           "",
		Data:            floatSliceToAny(data),
		BackgroundColor: colors,
	}}

	return &Chart{
		Type: ChartTypePie,
		Data: &Data{
			Labels:   stringSliceToAny(labels),
			Datasets: datasets,
		},
	}
}

// NewScatterChart creates a scatter chart from x,y points
func NewScatterChart(datasets []ScatterDataset) *Chart {
	chartDatasets := make([]Dataset, len(datasets))

	for i, ds := range datasets {
		points := make([]any, len(ds.Points))
		for j, p := range ds.Points {
			points[j] = DataPoint{X: p.X, Y: p.Y, R: p.R}
		}
		chartDatasets[i] = Dataset{
			Label:           ds.Label,
			Data:            points,
			BackgroundColor: ds.BackgroundColor,
			BorderColor:     ds.BorderColor,
		}
	}

	return &Chart{
		Type: ChartTypeScatter,
		Data: &Data{
			Datasets: chartDatasets,
		},
	}
}

// NewTimeSeriesChart creates a time series line chart
func NewTimeSeriesChart(datasets []TimeSeriesDataset) *Chart {
	chartDatasets := make([]Dataset, len(datasets))

	for i, ds := range datasets {
		points := make([]any, len(ds.Points))
		for j, p := range ds.Points {
			points[j] = TimeDataPoint{X: p.X, Y: p.Y, R: p.R}
		}
		chartDatasets[i] = Dataset{
			Label:           ds.Label,
			Data:            points,
			BorderColor:     ds.BorderColor,
			BackgroundColor: ds.BackgroundColor,
			Fill:            boolPtr(ds.Fill),
			Tension:         ds.Tension,
		}
	}

	// Configure time scale
	options := &Options{
		Scales: &Scales{
			X: &Scale{
				Type: stringPtr("time"),
				Time: &TimeScale{
					Unit: stringPtr("day"),
				},
			},
		},
	}

	return &Chart{
		Type:    ChartTypeLine,
		Data:    &Data{Datasets: chartDatasets},
		Options: options,
	}
}

// Convenience types for common chart creation
type LineDataset struct {
	Label           string
	Data            []float64
	BorderColor     string
	BackgroundColor string
	Fill            bool
	Tension         *float64
}

type BarDataset struct {
	Label           string
	Data            []float64
	BackgroundColor any // string or []string
	BorderColor     string
	BorderWidth     *int
}

type ScatterDataset struct {
	Label           string
	Points          []ScatterPoint
	BackgroundColor string
	BorderColor     string
}

type ScatterPoint struct {
	X float64
	Y float64
	R *float64 // For bubble charts
}

type TimeSeriesDataset struct {
	Label           string
	Points          []TimeSeriesPoint
	BorderColor     string
	BackgroundColor string
	Fill            bool
	Tension         *float64
}

type TimeSeriesPoint struct {
	X time.Time
	Y float64
	R *float64
}

// Helper functions
func stringSliceToAny(ss []string) []any {
	result := make([]any, len(ss))
	for i, s := range ss {
		result[i] = s
	}
	return result
}

func floatSliceToAny(fs []float64) []any {
	result := make([]any, len(fs))
	for i, f := range fs {
		result[i] = f
	}
	return result
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
