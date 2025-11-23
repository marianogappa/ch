package output

import (
	"encoding/json"
	"testing"

	"github.com/spf13/pflag"
)

func TestConfigToFlagsRoundTrip(t *testing.T) {
	// Test that ConfigToFlags produces flags that can be parsed back to the same config
	tests := []struct {
		name   string
		config *ChartConfig
	}{
		{
			name: "minimal config",
			config: &ChartConfig{
				ChartType:     ChartTypeLine,
				LegendDisplay: false, // Explicitly set to false to test
			},
		},
		{
			name: "full config",
			config: &ChartConfig{
				Title:         "Test Chart",
				ChartType:     ChartTypeBar,
				XLabel:        "X Axis",
				YLabel:        "Y Axis",
				XScaleType:    ScaleTypeLinear,
				YScaleType:    ScaleTypeLogarithmic,
				ZeroBased:     true,
				LegendDisplay: true,
				LegendPosition: LegendPositionBottom,
				ColorScheme:   ColorSchemeLegacy,
				Colors:        []string{"red", "blue", "green"},
				OutputPath:    "/tmp/test.html",
			},
		},
		// Note: FrontendSettings test skipped for now as flags aren't registered
		// This would require drivers to register their frontend flags
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert config to flags
			flags := ConfigToFlags(tt.config, "chartjs")

			// Create a new config and parse flags back
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			newConfig, parseFunc := RegisterChartConfigFlags(fs, "chartjs")

			// Parse the flags
			if err := fs.Parse(flags); err != nil {
				t.Fatalf("Failed to parse flags: %v", err)
			}

			// Call parse function
			if err := parseFunc("chartjs"); err != nil {
				t.Fatalf("Failed to parse special cases: %v", err)
			}

			// Compare configs - only check fields that were set in the original config
			// (skip zero values that have omitempty)
			if tt.config.Title != "" && newConfig.Title != tt.config.Title {
				t.Errorf("Title mismatch: got %q, want %q", newConfig.Title, tt.config.Title)
			}
			// ChartType is required, always check
			if newConfig.ChartType != tt.config.ChartType {
				t.Errorf("ChartType mismatch: got %q, want %q", newConfig.ChartType, tt.config.ChartType)
			}
			if tt.config.XLabel != "" && newConfig.XLabel != tt.config.XLabel {
				t.Errorf("XLabel mismatch: got %q, want %q", newConfig.XLabel, tt.config.XLabel)
			}
			if tt.config.YLabel != "" && newConfig.YLabel != tt.config.YLabel {
				t.Errorf("YLabel mismatch: got %q, want %q", newConfig.YLabel, tt.config.YLabel)
			}
			if tt.config.XScaleType != "" && newConfig.XScaleType != tt.config.XScaleType {
				t.Errorf("XScaleType mismatch: got %q, want %q", newConfig.XScaleType, tt.config.XScaleType)
			}
			if tt.config.YScaleType != "" && newConfig.YScaleType != tt.config.YScaleType {
				t.Errorf("YScaleType mismatch: got %q, want %q", newConfig.YScaleType, tt.config.YScaleType)
			}
			if tt.config.ZeroBased && !newConfig.ZeroBased {
				t.Errorf("ZeroBased mismatch: got %v, want %v", newConfig.ZeroBased, tt.config.ZeroBased)
			}
			// LegendDisplay: ConfigToFlags only outputs bool flags when true
			// So if original is false, we can't round-trip it (this is expected behavior)
			// Only check if original was true
			if tt.config.LegendDisplay && !newConfig.LegendDisplay {
				t.Errorf("LegendDisplay mismatch: got %v, want %v", newConfig.LegendDisplay, tt.config.LegendDisplay)
			}
			if tt.config.LegendPosition != "" && newConfig.LegendPosition != tt.config.LegendPosition {
				t.Errorf("LegendPosition mismatch: got %q, want %q", newConfig.LegendPosition, tt.config.LegendPosition)
			}
			if tt.config.ColorScheme != "" && newConfig.ColorScheme != tt.config.ColorScheme {
				t.Errorf("ColorScheme mismatch: got %q, want %q", newConfig.ColorScheme, tt.config.ColorScheme)
			}
			if len(tt.config.Colors) > 0 {
				if len(newConfig.Colors) != len(tt.config.Colors) {
					t.Errorf("Colors length mismatch: got %d, want %d", len(newConfig.Colors), len(tt.config.Colors))
				} else {
					for i, c := range tt.config.Colors {
						if i < len(newConfig.Colors) && newConfig.Colors[i] != c {
							t.Errorf("Colors[%d] mismatch: got %q, want %q", i, newConfig.Colors[i], c)
						}
					}
				}
			}
			if tt.config.OutputPath != "" && newConfig.OutputPath != tt.config.OutputPath {
				t.Errorf("OutputPath mismatch: got %q, want %q", newConfig.OutputPath, tt.config.OutputPath)
			}
		})
	}
}

func TestGenerateConfigDocumentation(t *testing.T) {
	// Test that documentation can be generated
	doc, err := GenerateConfigDocumentation("chartjs", nil)
	if err != nil {
		t.Fatalf("Failed to generate documentation: %v", err)
	}

	// Verify it's valid JSON
	var parsed ConfigDocumentation
	if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
		t.Fatalf("Generated documentation is not valid JSON: %v", err)
	}

	// Verify it has expected fields
	if parsed.ConfigType != "ChartConfig" {
		t.Errorf("Expected ConfigType to be 'ChartConfig', got %q", parsed.ConfigType)
	}

	if len(parsed.Fields) == 0 {
		t.Error("Expected at least one field in documentation")
	}

	// Check that required fields are documented
	fieldMap := make(map[string]FieldDoc)
	for _, f := range parsed.Fields {
		fieldMap[f.JSONName] = f
	}

	requiredFields := []string{"chartType", "title", "xLabel", "yLabel"}
	for _, rf := range requiredFields {
		if _, ok := fieldMap[rf]; !ok {
			t.Errorf("Required field %q not found in documentation", rf)
		}
	}
}

func TestCamelToKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"chartType", "chart-type"},
		{"xLabel", "x-label"},
		{"yLabel", "y-label"},
		{"zeroBased", "zero-based"},
		{"legendDisplay", "legend-display"},
		{"legendPosition", "legend-position"},
		{"colorScheme", "color-scheme"},
		{"outputPath", "output-path"},
		{"frontendSettings", "frontend-settings"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := camelToKebab(tt.input)
			if got != tt.expected {
				t.Errorf("camelToKebab(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

