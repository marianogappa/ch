package output

import "github.com/marianogappa/ch/pkg/ch"

// ChartOutput is an optional interface that chart drivers can implement
// to provide chart-specific functionality
type ChartOutput interface {
	// RenderChart renders the chart using the generic ChartConfig
	// This is an alternative to the standard Render method that uses ChartConfig
	RenderChart(rows <-chan ch.Row, config *ChartConfig) error

	// GetFrontendSettings returns a map of frontend-specific setting definitions
	// The key is the setting name, and the value is its definition
	// Returns nil or empty map if no frontend-specific settings are supported
	GetFrontendSettings() map[string]SettingDefinition
}

// IsChartOutput checks if an output driver implements ChartOutput
func IsChartOutput(driver any) bool {
	_, ok := driver.(ChartOutput)
	return ok
}
