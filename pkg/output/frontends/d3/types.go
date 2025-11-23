package d3

// ChartType represents the type of chart to render
type ChartType string

const (
	ChartTypeBar       ChartType = "bar"
	ChartTypePie       ChartType = "pie"
	ChartTypeScatter   ChartType = "scatter"
	ChartTypeHistogram ChartType = "histogram"
	ChartTypeLine      ChartType = "line"
)

// Chart represents the complete D3 chart configuration
type Chart struct {
	Type   ChartType `json:"type"`
	Data   []any     `json:"data"`
	Config *Config   `json:"config,omitempty"`
}

// Config holds chart configuration options
type Config struct {
	Title     string    `json:"title"`
	ChartType ChartType `json:"chartType"`
	XLabel    string    `json:"xLabel,omitempty"`
	YLabel    string    `json:"yLabel,omitempty"`
	Color     string    `json:"color,omitempty"`
	Colors    []string  `json:"colors,omitempty"`
}
