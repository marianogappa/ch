package output

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/marianogappa/ch/pkg/ch"
)

// SettingsCollector collects all settings from generic ChartConfig and output-specific settings
type SettingsCollector struct {
	GenericSettings []SettingDefinition            `json:"genericSettings"`
	OutputSettings  map[string][]SettingDefinition `json:"outputSettings"` // key is output name
}

// CollectAllSettings collects all settings from all registered outputs
func CollectAllSettings() (*SettingsCollector, error) {
	collector := &SettingsCollector{
		GenericSettings: []SettingDefinition{},
		OutputSettings:  make(map[string][]SettingDefinition),
	}

	// Collect generic settings from ChartConfig
	genericSettings, err := collectGenericSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to collect generic settings: %w", err)
	}
	collector.GenericSettings = genericSettings

	// Collect output-specific settings
	outputNames := ch.Outputs()
	for _, outputName := range outputNames {
		outDriver, err := ch.GetOutput(outputName)
		if err != nil {
			continue
		}

		if chartOutput, ok := outDriver.(ChartOutput); ok {
			frontendSettings := chartOutput.GetFrontendSettings()
			if len(frontendSettings) > 0 {
				settings := make([]SettingDefinition, 0, len(frontendSettings))
				for _, setting := range frontendSettings {
					settings = append(settings, setting)
				}
				collector.OutputSettings[outputName] = settings
			}
		}
	}

	return collector, nil
}

// collectGenericSettings collects settings from ChartConfig struct tags
func collectGenericSettings() ([]SettingDefinition, error) {
	var settings []SettingDefinition

	t := reflect.TypeOf(ChartConfig{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		// Get flag tag - skip if "-"
		flagTag := field.Tag.Get("flag")
		if flagTag == "" || flagTag == "-" {
			continue
		}

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			jsonName = field.Name
		}

		// Get description
		description := field.Tag.Get("desc")
		if description == "" {
			description = fmt.Sprintf("Sets the %s", jsonName)
		}

		// Get default value
		defaultValue := field.Tag.Get("default")

		// Get valid values
		validValuesStr := field.Tag.Get("valid")
		var validValues []string
		if validValuesStr != "" {
			validValues = strings.Split(validValuesStr, ",")
			for i := range validValues {
				validValues[i] = strings.TrimSpace(validValues[i])
			}
		}

		// Get examples
		examplesStr := field.Tag.Get("examples")
		var examples []string
		if examplesStr != "" {
			examples = strings.Split(examplesStr, ",")
			for i := range examples {
				examples[i] = strings.TrimSpace(examples[i])
			}
		}

		// Determine type
		settingType := "string"
		switch field.Type.Kind() {
		case reflect.Bool:
			settingType = "bool"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			settingType = "number"
		case reflect.Slice:
			settingType = "array"
		}

		// Check if it's an enum type (has valid values)
		if len(validValues) > 0 {
			settingType = "enum"
		}

		// Check if required (no omitempty in JSON tag)
		required := !strings.Contains(jsonTag, "omitempty")

		// Handle special cases for custom types
		if field.Type == reflect.TypeOf(ChartType("")) {
			settingType = "enum"
			if len(validValues) == 0 {
				validValues = []string{
					string(ChartTypeLine),
					string(ChartTypeBar),
					string(ChartTypePie),
					string(ChartTypeDoughnut),
					string(ChartTypeScatter),
					string(ChartTypeBubble),
					string(ChartTypeArea),
					string(ChartTypeHistogram),
					string(ChartTypeRadar),
					string(ChartTypePolarArea),
				}
			}
		} else if field.Type == reflect.TypeOf(ScaleType("")) {
			settingType = "enum"
			if len(validValues) == 0 {
				validValues = []string{
					string(ScaleTypeLinear),
					string(ScaleTypeLogarithmic),
					string(ScaleTypeTime),
					string(ScaleTypeCategory),
				}
			}
		} else if field.Type == reflect.TypeOf(LegendPosition("")) {
			settingType = "enum"
			if len(validValues) == 0 {
				validValues = []string{
					string(LegendPositionTop),
					string(LegendPositionBottom),
					string(LegendPositionLeft),
					string(LegendPositionRight),
					string(LegendPositionChart),
				}
			}
		} else if field.Type == reflect.TypeOf(ColorScheme("")) {
			settingType = "enum"
			if len(validValues) == 0 {
				validValues = []string{
					string(ColorSchemeDefault),
					string(ColorSchemeLegacy),
					string(ColorSchemeGradient),
					string(ColorSchemeCustom),
				}
			}
		}

		settings = append(settings, SettingDefinition{
			Name:         flagTag,
			Type:         settingType,
			Description:  description,
			DefaultValue: defaultValue,
			ValidValues:  validValues,
			Examples:     examples,
			Required:     required,
		})
	}

	return settings, nil
}
