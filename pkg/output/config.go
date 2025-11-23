package output

// ChartType represents the type of chart to render
type ChartType string

const (
	ChartTypeLine      ChartType = "line"
	ChartTypeBar       ChartType = "bar"
	ChartTypePie       ChartType = "pie"
	ChartTypeDoughnut  ChartType = "doughnut"
	ChartTypeScatter   ChartType = "scatter"
	ChartTypeBubble    ChartType = "bubble"
	ChartTypeArea      ChartType = "area"
	ChartTypeHistogram ChartType = "histogram"
	ChartTypeRadar     ChartType = "radar"
	ChartTypePolarArea ChartType = "polarArea"
)

// ScaleType represents the scale type for axes
type ScaleType string

const (
	ScaleTypeLinear      ScaleType = "linear"
	ScaleTypeLogarithmic ScaleType = "logarithmic"
	ScaleTypeTime        ScaleType = "time"
	ScaleTypeCategory    ScaleType = "category"
)

// LegendPosition represents where the legend should be positioned
type LegendPosition string

const (
	LegendPositionTop    LegendPosition = "top"
	LegendPositionBottom LegendPosition = "bottom"
	LegendPositionLeft   LegendPosition = "left"
	LegendPositionRight  LegendPosition = "right"
	LegendPositionChart  LegendPosition = "chartArea"
)

// ColorScheme represents a color scheme for the chart
type ColorScheme string

const (
	ColorSchemeDefault  ColorScheme = "default"
	ColorSchemeLegacy   ColorScheme = "legacy"
	ColorSchemeGradient ColorScheme = "gradient"
	ColorSchemeCustom   ColorScheme = "custom"
)

// SettingDefinition represents a single setting with all its metadata
type SettingDefinition struct {
	Name         string   `json:"name"`                   // Setting name (e.g., "animation", "title")
	Type         string   `json:"type"`                   // Type: "bool", "string", "number", "enum"
	Description  string   `json:"description"`            // Human-readable description
	DefaultValue string   `json:"defaultValue,omitempty"` // Default value as string
	ValidValues  []string `json:"validValues,omitempty"`  // Valid values for enum types
	Examples     []string `json:"examples,omitempty"`     // Example values
	Required     bool     `json:"required,omitempty"`     // Whether this setting is required
}

// ChartConfig represents a generic chart configuration
// that can be used across different frontend implementations
type ChartConfig struct {
	// Basic chart properties
	Title     string    `json:"title,omitempty" flag:"title" desc:"Sets the title for the chart" default:"" valid:"" examples:"Sales Report,Monthly Revenue,User Growth"`
	ChartType ChartType `json:"chartType" flag:"chart-type" desc:"Chart type" default:"line" valid:"line,bar,pie,doughnut,scatter,bubble,area,histogram,radar,polarArea" examples:""`

	// Axis labels
	XLabel string `json:"xLabel,omitempty" flag:"x-label" desc:"Sets the label for the x axis" default:"" valid:"" examples:"Date,Time,Category"`
	YLabel string `json:"yLabel,omitempty" flag:"y-label" desc:"Sets the label for the y axis" default:"" valid:"" examples:"Value,Count,Revenue"`

	// Scale configuration
	XScaleType ScaleType `json:"xScaleType,omitempty" flag:"x-scale-type" desc:"X-axis scale type" default:"category" valid:"category,linear,logarithmic,time"`
	YScaleType ScaleType `json:"yScaleType,omitempty" flag:"y-scale-type" desc:"Y-axis scale type" default:"linear" valid:"linear,logarithmic"`

	// Y-axis configuration
	ZeroBased bool `json:"zeroBased,omitempty" flag:"zero-based" desc:"Makes y-axis begin at zero" default:"false" valid:""`

	// Legend configuration
	LegendDisplay  bool           `json:"legendDisplay,omitempty" flag:"legend-display" desc:"Show or hide the legend" default:"true" valid:""`
	LegendPosition LegendPosition `json:"legendPosition,omitempty" flag:"legend-position" desc:"Legend position" default:"top" valid:"top,bottom,left,right,chartArea"`

	// Color configuration
	ColorScheme ColorScheme `json:"colorScheme,omitempty" flag:"color-scheme" desc:"Color scheme" default:"default" valid:"default,legacy,gradient,custom"`
	Colors      []string    `json:"colors,omitempty" flag:"colors" desc:"Custom colors when ColorScheme is 'custom' (comma-separated)" default:"" valid:""`

	// Frontend-specific settings
	// This allows frontends to accept settings that are not part of the generic config
	FrontendSettings map[string]any `json:"frontendSettings,omitempty" flag:"-" desc:"Frontend-specific settings (use with output name prefix)" default:"" valid:""`

	// OutputPath specifies where to write the output
	// - Empty string: use default behavior (open browser for HTML, stdout for JSON)
	// - "-": write to stdout
	// - Any other string: write to that file path
	OutputPath string `json:"outputPath,omitempty" flag:"output-path" desc:"Output path: empty for browser/stdout, '-' for stdout, or file path" default:"" valid:"" examples:"chart.html,output/chart.html,-"`
}
