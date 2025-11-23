package chartjs_test

import (
	"os"
	"testing"
	"time"

	"github.com/marianogappa/ch/pkg/output/frontends/chartjs"
)

func TestExampleBasicLineChart(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a simple line chart
	chart := chartjs.NewLineChart(
		[]string{"Jan", "Feb", "Mar", "Apr", "May"},
		[]chartjs.LineDataset{
			{
				Label:           "Sales",
				Data:            []float64{10, 19, 15, 25, 22},
				BorderColor:     "rgb(75, 192, 192)",
				BackgroundColor: "rgba(75, 192, 192, 0.2)",
				Fill:            true,
			},
		},
	)

	// Apply options using helper functions
	chart.Options = &chartjs.Options{}
	chartjs.ApplyOptions(chart.Options,
		chartjs.WithTitle("Monthly Sales"),
		chartjs.WithXAxisLabel("Month"),
		chartjs.WithYAxisLabel("Sales ($)"),
		chartjs.WithYAxisBeginAtZero(true),
		chartjs.WithLegend(true, "top"),
	)

	// Build and output
	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

func TestExampleBarChart(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a bar chart with multiple datasets
	chart := chartjs.NewBarChart(
		[]string{"Q1", "Q2", "Q3", "Q4"},
		[]chartjs.BarDataset{
			{
				Label:           "2023",
				Data:            []float64{65, 59, 80, 81},
				BackgroundColor: []string{"#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0"},
			},
			{
				Label:           "2024",
				Data:            []float64{28, 48, 40, 19},
				BackgroundColor: []string{"#FF9F40", "#FF6384", "#C9CBCF", "#4BC0C0"},
			},
		},
	)

	chart.Options = &chartjs.Options{}
	chartjs.ApplyOptions(chart.Options,
		chartjs.WithTitle("Quarterly Comparison"),
		chartjs.WithStacked(true),
	)

	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

func TestExampleTimeSeriesChart(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a time series chart
	now := time.Now()
	chart := chartjs.NewTimeSeriesChart(
		[]chartjs.TimeSeriesDataset{
			{
				Label: "Temperature",
				Points: []chartjs.TimeSeriesPoint{
					{X: now.Add(-24 * time.Hour), Y: 20.5},
					{X: now.Add(-18 * time.Hour), Y: 22.1},
					{X: now.Add(-12 * time.Hour), Y: 19.8},
					{X: now.Add(-6 * time.Hour), Y: 21.3},
					{X: now, Y: 20.9},
				},
				BorderColor:     "rgb(255, 99, 132)",
				BackgroundColor: "rgba(255, 99, 132, 0.2)",
				Fill:            true,
			},
		},
	)

	chart.Options = &chartjs.Options{}
	chartjs.ApplyOptions(chart.Options,
		chartjs.WithTitle("Temperature Over Time"),
		chartjs.WithTimeScale("hour"),
	)

	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

func TestExampleScatterChart(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a scatter chart
	chart := chartjs.NewScatterChart(
		[]chartjs.ScatterDataset{
			{
				Label: "Dataset 1",
				Points: []chartjs.ScatterPoint{
					{X: 10, Y: 20},
					{X: 15, Y: 10},
					{X: 20, Y: 30},
					{X: 25, Y: 15},
				},
				BackgroundColor: "rgba(255, 99, 132, 0.5)",
				BorderColor:     "rgb(255, 99, 132)",
			},
		},
	)

	chart.Options = &chartjs.Options{}
	chartjs.ApplyOptions(chart.Options,
		chartjs.WithTitle("Scatter Plot"),
		chartjs.WithXAxisLabel("X Axis"),
		chartjs.WithYAxisLabel("Y Axis"),
	)

	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

func TestExamplePieChart(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a pie chart
	chart := chartjs.NewPieChart(
		[]string{"Red", "Blue", "Yellow", "Green", "Purple"},
		[]float64{12, 19, 3, 5, 2},
		[]string{
			"rgba(255, 99, 132, 0.8)",
			"rgba(54, 162, 235, 0.8)",
			"rgba(255, 206, 86, 0.8)",
			"rgba(75, 192, 192, 0.8)",
			"rgba(153, 102, 255, 0.8)",
		},
	)

	chart.Options = &chartjs.Options{}
	chartjs.ApplyOptions(chart.Options,
		chartjs.WithTitle("Color Distribution"),
		chartjs.WithLegend(true, "right"),
	)

	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

func TestExampleAdvancedConfiguration(t *testing.T) {
	t.Skip("Example usage - not a real test")
	// Create a chart with full manual configuration
	chart := &chartjs.Chart{
		Type: chartjs.ChartTypeLine,
		Data: &chartjs.Data{
			Labels: []any{"Jan", "Feb", "Mar", "Apr", "May"},
			Datasets: []chartjs.Dataset{
				{
					Label:           "Dataset 1",
					Data:            []any{10.0, 19.0, 15.0, 25.0, 22.0},
					BorderColor:     "rgb(75, 192, 192)",
					BackgroundColor: "rgba(75, 192, 192, 0.2)",
					Fill:            boolPtr(true),
					Tension:         float64Ptr(0.4),
				},
			},
		},
		Options: &chartjs.Options{
			Responsive: boolPtr(true),
			Plugins: &chartjs.Plugins{
				Title: &chartjs.Title{
					Display:  boolPtr(true),
					Text:     stringPtr("Advanced Chart"),
					Position: stringPtr("top"),
				},
				Legend: &chartjs.Legend{
					Display:  boolPtr(true),
					Position: stringPtr("bottom"),
				},
				Tooltip: &chartjs.Tooltip{
					Enabled: boolPtr(true),
					Mode:    stringPtr("index"),
				},
			},
			Scales: &chartjs.Scales{
				Y: &chartjs.Scale{
					BeginAtZero: boolPtr(true),
					Title: &chartjs.ScaleTitle{
						Display: boolPtr(true),
						Text:    stringPtr("Value"),
					},
				},
			},
			Elements: &chartjs.Elements{
				Line: &chartjs.LineElement{
					Tension: float64Ptr(0.4),
				},
			},
		},
	}

	cjs := chartjs.New(chart)
	cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}

// Helper functions for examples
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

