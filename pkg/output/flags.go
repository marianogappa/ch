package output

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

// FieldDoc represents documentation for a config field
type FieldDoc struct {
	FieldName    string   `json:"fieldName"`
	JSONName     string   `json:"jsonName"`
	FlagName     string   `json:"flagName"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	EnumValues   []string `json:"enumValues,omitempty"`
	Required     bool     `json:"required,omitempty"`
	OutputName   string   `json:"outputName,omitempty"` // For frontend-specific flags
}

// ConfigDocumentation represents full documentation for a config struct
type ConfigDocumentation struct {
	ConfigType   string     `json:"configType"`
	Description  string     `json:"description"`
	Fields       []FieldDoc `json:"fields"`
	FrontendDocs string     `json:"frontendDocs,omitempty"` // For frontend-specific settings
}

// AllFlagsDocumentation represents documentation for all flags from all outputs
type AllFlagsDocumentation struct {
	GenericFlags []FieldDoc            `json:"genericFlags"` // ChartConfig flags
	OutputFlags  map[string][]FieldDoc `json:"outputFlags"`  // Frontend-specific flags by output name
}

// flagParser holds temporary flag values that need to be parsed after flag parsing
type flagParser struct {
	chartTypeStr   string
	xScaleTypeStr  string
	yScaleTypeStr  string
	legendPosStr   string
	colorSchemeStr string
	colorsStr      string
	chartTypePtr   *ChartType
	xScaleTypePtr  *ScaleType
	yScaleTypePtr  *ScaleType
	legendPosPtr   *LegendPosition
	colorSchemePtr *ColorScheme
	colorsPtr      *[]string
}

// RegisterChartConfigFlags registers flags for ChartConfig on the given flag set
// Returns the config struct that will be populated when flags are parsed
// and a function to call after flag parsing to handle special cases
func RegisterChartConfigFlags(fs *pflag.FlagSet, outputName string) (*ChartConfig, func(outputName string) error) {
	cfg := &ChartConfig{
		ChartType:        ChartTypeLine,
		XScaleType:       ScaleTypeCategory,
		YScaleType:       ScaleTypeLinear,
		LegendDisplay:    true,
		LegendPosition:   LegendPositionTop,
		ColorScheme:      ColorSchemeDefault,
		FrontendSettings: make(map[string]any),
	}

	parser := &flagParser{
		chartTypeStr:   string(ChartTypeLine),
		xScaleTypeStr:  string(ScaleTypeCategory),
		yScaleTypeStr:  string(ScaleTypeLinear),
		legendPosStr:   string(LegendPositionTop),
		colorSchemeStr: string(ColorSchemeDefault),
		chartTypePtr:   &cfg.ChartType,
		xScaleTypePtr:  &cfg.XScaleType,
		yScaleTypePtr:  &cfg.YScaleType,
		legendPosPtr:   &cfg.LegendPosition,
		colorSchemePtr: &cfg.ColorScheme,
		colorsPtr:      &cfg.Colors,
	}

	// Get settings from collector to include examples in descriptions
	collector, _ := CollectAllSettings()
	settingMap := make(map[string]SettingDefinition)
	if collector != nil {
		for _, setting := range collector.GenericSettings {
			settingMap[setting.Name] = setting
		}
	}

	// Helper to build description with examples
	buildDesc := func(name string, defaultDesc string) string {
		if setting, ok := settingMap[name]; ok {
			desc := setting.Description
			if len(setting.Examples) > 0 {
				desc += fmt.Sprintf(" (examples: %s)", strings.Join(setting.Examples, ", "))
			}
			return desc
		}
		return defaultDesc
	}

	// Register flags manually for each field to handle types correctly
	fs.StringVar(&cfg.Title, "title", "", buildDesc("title", "Sets the title for the chart"))
	fs.StringVar(&parser.chartTypeStr, "chart-type", string(ChartTypeLine), buildDesc("chart-type", "Chart type: line, bar, pie, doughnut, scatter, bubble, area, histogram, radar, polarArea"))
	fs.StringVar(&cfg.XLabel, "x-label", "", buildDesc("x-label", "Sets the label for the x axis"))
	fs.StringVar(&cfg.YLabel, "y-label", "", buildDesc("y-label", "Sets the label for the y axis"))
	fs.StringVar(&parser.xScaleTypeStr, "x-scale-type", string(ScaleTypeCategory), buildDesc("x-scale-type", "X-axis scale type: category, linear, logarithmic, time"))
	fs.StringVar(&parser.yScaleTypeStr, "y-scale-type", string(ScaleTypeLinear), buildDesc("y-scale-type", "Y-axis scale type: linear, logarithmic"))
	fs.BoolVar(&cfg.ZeroBased, "zero-based", false, buildDesc("zero-based", "Makes y-axis begin at zero"))
	fs.BoolVar(&cfg.LegendDisplay, "legend-display", true, buildDesc("legend-display", "Show or hide the legend"))
	fs.StringVar(&parser.legendPosStr, "legend-position", string(LegendPositionTop), buildDesc("legend-position", "Legend position: top, bottom, left, right, chartArea"))
	fs.StringVar(&parser.colorSchemeStr, "color-scheme", string(ColorSchemeDefault), buildDesc("color-scheme", "Color scheme: default, legacy, gradient, custom"))

	// Handle Colors slice specially
	fs.StringVar(&parser.colorsStr, "colors", "", buildDesc("colors", "Custom colors when ColorScheme is 'custom' (comma-separated)"))

	fs.StringVar(&cfg.OutputPath, "output-path", "", buildDesc("output-path", "Output path: empty for browser/stdout, '-' for stdout, or file path"))

	// Return a function to parse special cases after flag parsing
	parseFunc := func(selectedOutputName string) error {
		*parser.chartTypePtr = ChartType(parser.chartTypeStr)
		*parser.xScaleTypePtr = ScaleType(parser.xScaleTypeStr)
		*parser.yScaleTypePtr = ScaleType(parser.yScaleTypeStr)
		*parser.legendPosPtr = LegendPosition(parser.legendPosStr)
		*parser.colorSchemePtr = ColorScheme(parser.colorSchemeStr)

		if parser.colorsStr != "" {
			*parser.colorsPtr = strings.Split(parser.colorsStr, ",")
			// Trim spaces
			for i := range *parser.colorsPtr {
				(*parser.colorsPtr)[i] = strings.TrimSpace((*parser.colorsPtr)[i])
			}
		}
		return nil
	}

	return cfg, parseFunc
}

// RegisterFrontendSettingsFlags registers flags for frontend-specific settings with prefix
func RegisterFrontendSettingsFlags(fs *pflag.FlagSet, outputName string, settings map[string]any, docs map[string]string) {
	// This is a placeholder - frontend-specific flags will be registered dynamically
	// based on what the driver supports. For now, we'll use a generic approach.
	// Drivers can register their own flags in RegisterFlags if needed.
}

// camelToKebab converts camelCase to kebab-case
func camelToKebab(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// GenerateConfigDocumentation generates LLM-readable documentation for ChartConfig
// This is a legacy function that uses the old FieldDoc format for backward compatibility
// New code should use GenerateAllFlagsDocumentation() instead
func GenerateConfigDocumentation(outputName string, chartOutput ChartOutput) (string, error) {
	// Use SettingsCollector to get generic settings
	collector, err := CollectAllSettings()
	if err != nil {
		return "", err
	}

	// Convert SettingDefinition to FieldDoc for backward compatibility
	doc := ConfigDocumentation{
		ConfigType:  "ChartConfig",
		Description: "Configuration for chart rendering",
		Fields:      []FieldDoc{},
	}

	for _, setting := range collector.GenericSettings {
		// Convert SettingDefinition to FieldDoc
		fieldDoc := FieldDoc{
			FlagName:    setting.Name,
			Type:        setting.Type,
			Description: setting.Description,
		}
		if setting.DefaultValue != "" {
			fieldDoc.DefaultValue = setting.DefaultValue
		}
		if len(setting.ValidValues) > 0 {
			fieldDoc.EnumValues = setting.ValidValues
		}
		fieldDoc.Required = setting.Required

		// Try to infer JSONName and FieldName from flag name
		// This is approximate but should work for most cases
		camel := kebabToCamel(setting.Name)
		fieldDoc.JSONName = camel
		// Capitalize first letter for FieldName
		if len(camel) > 0 {
			fieldDoc.FieldName = strings.ToUpper(camel[:1]) + camel[1:]
		} else {
			fieldDoc.FieldName = camel
		}

		doc.Fields = append(doc.Fields, fieldDoc)
	}

	// Add frontend-specific documentation if available
	if chartOutput != nil {
		frontendSettings := chartOutput.GetFrontendSettings()
		if len(frontendSettings) > 0 {
			frontendDoc := fmt.Sprintf("Frontend-specific settings (use with prefix %s-):\n", outputName)
			for name, setting := range frontendSettings {
				frontendDoc += fmt.Sprintf("  - %s: %s (type: %s", name, setting.Description, setting.Type)
				if setting.DefaultValue != "" {
					frontendDoc += fmt.Sprintf(", default: %s", setting.DefaultValue)
				}
				if len(setting.ValidValues) > 0 {
					frontendDoc += fmt.Sprintf(", valid: %v", setting.ValidValues)
				}
				frontendDoc += ")\n"
			}
			doc.FrontendDocs = frontendDoc
		}
	}

	// Marshal to JSON for LLM consumption
	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// kebabToCamel converts kebab-case to camelCase
func kebabToCamel(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) == 0 {
		return s
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			// Capitalize first letter
			first := strings.ToUpper(parts[i][:1])
			rest := parts[i][1:]
			result += first + rest
		}
	}
	return result
}

// ConfigToFlags converts a ChartConfig to command-line flags
// This is used by LLM to output the flags that correspond to a chosen configuration
func ConfigToFlags(cfg *ChartConfig, outputName string) []string {
	var flags []string

	// Use reflection to convert config to flags
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(*cfg)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanInterface() {
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

		flagName := camelToKebab(jsonName)

		// Skip if zero value and has omitempty
		if strings.Contains(jsonTag, "omitempty") {
			zeroVal := reflect.Zero(fieldValue.Type())
			if reflect.DeepEqual(fieldValue.Interface(), zeroVal.Interface()) {
				continue
			}
		}

		// Convert value to flag
		switch fieldValue.Kind() {
		case reflect.String:
			val := fieldValue.String()
			if val != "" {
				flags = append(flags, fmt.Sprintf("--%s=%s", flagName, val))
			}
		case reflect.Bool:
			val := fieldValue.Bool()
			if val {
				flags = append(flags, fmt.Sprintf("--%s", flagName))
			}
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				val := fieldValue.Interface().([]string)
				if len(val) > 0 {
					flags = append(flags, fmt.Sprintf("--%s=%s", flagName, strings.Join(val, ",")))
				}
			}
		default:
			// For custom types, convert to string
			val := fmt.Sprintf("%v", fieldValue.Interface())
			if val != "" && val != "<nil>" {
				flags = append(flags, fmt.Sprintf("--%s=%s", flagName, val))
			}
		}
	}

	// Handle frontend-specific settings with prefix
	for key, value := range cfg.FrontendSettings {
		prefixedKey := outputName + "-" + camelToKebab(key)
		flags = append(flags, fmt.Sprintf("--%s=%v", prefixedKey, value))
	}

	return flags
}

// frontendFlagInfo holds information about a frontend-specific flag
type frontendFlagInfo struct {
	outputName string
	flagName   string
	flagType   string
}

var globalFrontendFlagInfos = make(map[string]frontendFlagInfo) // key is prefixed flag name

// RegisterAllFlags registers all flags (generic + frontend-specific) to the given flag set
// Returns the config and a parse function that populates frontend settings for the selected output
func RegisterAllFlags(fs *pflag.FlagSet) (*ChartConfig, func(outputName string) error) {
	// Register generic ChartConfig flags
	cfg, parseFunc := RegisterChartConfigFlags(fs, "")

	// Get all settings from SettingsCollector
	collector, err := CollectAllSettings()
	if err != nil {
		// If collection fails, just register generic flags
		return cfg, parseFunc
	}

	// Register frontend-specific flags for all outputs
	for outputName, settings := range collector.OutputSettings {
		for _, setting := range settings {
			prefixedFlagName := outputName + "-" + setting.Name

			// Store flag info for later lookup
			globalFrontendFlagInfos[prefixedFlagName] = frontendFlagInfo{
				outputName: outputName,
				flagName:   setting.Name,
				flagType:   setting.Type,
			}

			// Build description with examples if available
			desc := setting.Description
			if len(setting.Examples) > 0 {
				desc += fmt.Sprintf(" (examples: %s)", strings.Join(setting.Examples, ", "))
			}

			// Register flag based on type
			switch setting.Type {
			case "bool":
				defaultVal := setting.DefaultValue == "true"
				var val bool
				fs.BoolVar(&val, prefixedFlagName, defaultVal, fmt.Sprintf("[%s] %s", outputName, desc))
			case "number":
				var val string
				fs.StringVar(&val, prefixedFlagName, setting.DefaultValue, fmt.Sprintf("[%s] %s", outputName, desc))
			default: // string, enum
				var val string
				fs.StringVar(&val, prefixedFlagName, setting.DefaultValue, fmt.Sprintf("[%s] %s", outputName, desc))
			}
		}
	}

	// Create combined parse function that queries the flag set
	combinedParseFunc := func(selectedOutputName string) error {
		if parseFunc != nil {
			if err := parseFunc(selectedOutputName); err != nil {
				return err
			}
		}

		// After parsing, copy frontend settings for the selected output to the config
		if cfg != nil {
			if cfg.FrontendSettings == nil {
				cfg.FrontendSettings = make(map[string]any)
			}

			// Query flag values from the flag set for the selected output
			for prefixedFlagName, flagInfo := range globalFrontendFlagInfos {
				if flagInfo.outputName == selectedOutputName {
					flag := fs.Lookup(prefixedFlagName)
					if flag == nil {
						continue
					}

					switch flagInfo.flagType {
					case "bool":
						if val, err := fs.GetBool(prefixedFlagName); err == nil && val {
							cfg.FrontendSettings[flagInfo.flagName] = val
						}
					case "number":
						if val, err := fs.GetString(prefixedFlagName); err == nil && val != "" {
							cfg.FrontendSettings[flagInfo.flagName] = val
						}
					default: // string, enum
						if val, err := fs.GetString(prefixedFlagName); err == nil && val != "" {
							cfg.FrontendSettings[flagInfo.flagName] = val
						}
					}
				}
			}
		}

		return nil
	}

	return cfg, combinedParseFunc
}

// GenerateAllFlagsDocumentation generates LLM-readable documentation for all flags
// This is a legacy function that generates the old format
// New code should use GenerateChartConfigDocumentation() instead
func GenerateAllFlagsDocumentation() (string, error) {
	collector, err := CollectAllSettings()
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(collector, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// ChartConfigDocumentation represents documentation for creating ChartConfig JSON
type ChartConfigDocumentation struct {
	Description      string                          `json:"description"`
	Example          *ChartConfig                    `json:"example"`
	Fields           []ChartConfigFieldDocumentation `json:"fields"`
	FrontendSettings map[string]FrontendSettingsDoc  `json:"frontendSettings"`
}

// ChartConfigFieldDocumentation documents a single ChartConfig field
type ChartConfigFieldDocumentation struct {
	JSONName     string   `json:"jsonName"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	ValidValues  []string `json:"validValues,omitempty"`
	Examples     []string `json:"examples,omitempty"`
	Required     bool     `json:"required,omitempty"`
}

// FrontendSettingsDoc documents frontend-specific settings
type FrontendSettingsDoc struct {
	Description string                          `json:"description"`
	Settings    []ChartConfigFieldDocumentation `json:"settings"`
}

// GenerateChartConfigDocumentation generates documentation that teaches the LLM
// how to create a ChartConfig JSON structure
func GenerateChartConfigDocumentation() (string, error) {
	collector, err := CollectAllSettings()
	if err != nil {
		return "", err
	}

	// Build example ChartConfig
	example := &ChartConfig{
		Title:          "Example Chart",
		ChartType:      ChartTypeLine,
		XLabel:         "X Axis",
		YLabel:         "Y Axis",
		XScaleType:     ScaleTypeCategory,
		YScaleType:     ScaleTypeLinear,
		ZeroBased:      false,
		LegendDisplay:  true,
		LegendPosition: LegendPositionTop,
		ColorScheme:    ColorSchemeDefault,
		Colors:         []string{},
		FrontendSettings: map[string]any{
			"animation": true,
		},
		OutputPath: "",
	}

	// Convert generic settings to field documentation
	// Use reflection to get JSON names directly from ChartConfig struct tags
	chartConfigType := reflect.TypeOf(ChartConfig{})
	flagToJSONName := make(map[string]string)
	for i := 0; i < chartConfigType.NumField(); i++ {
		field := chartConfigType.Field(i)
		flagTag := field.Tag.Get("flag")
		jsonTag := field.Tag.Get("json")
		if flagTag != "" && flagTag != "-" && jsonTag != "" && jsonTag != "-" {
			jsonName := strings.Split(jsonTag, ",")[0]
			flagToJSONName[flagTag] = jsonName
		}
	}

	var fields []ChartConfigFieldDocumentation
	for _, setting := range collector.GenericSettings {
		// Get JSON name from the map, fallback to kebabToCamel if not found
		jsonName, ok := flagToJSONName[setting.Name]
		if !ok {
			jsonName = kebabToCamel(setting.Name)
		}

		fieldDoc := ChartConfigFieldDocumentation{
			JSONName:     jsonName,
			Type:         setting.Type,
			Description:  setting.Description,
			DefaultValue: setting.DefaultValue,
			ValidValues:  setting.ValidValues,
			Examples:     setting.Examples,
			Required:     setting.Required,
		}
		fields = append(fields, fieldDoc)
	}

	// Build frontend settings documentation
	frontendDocs := make(map[string]FrontendSettingsDoc)
	for outputName, settings := range collector.OutputSettings {
		var settingDocs []ChartConfigFieldDocumentation
		for _, setting := range settings {
			settingDoc := ChartConfigFieldDocumentation{
				JSONName:     setting.Name, // Frontend settings use their flag name as-is
				Type:         setting.Type,
				Description:  setting.Description,
				DefaultValue: setting.DefaultValue,
				ValidValues:  setting.ValidValues,
				Examples:     setting.Examples,
				Required:     setting.Required,
			}
			settingDocs = append(settingDocs, settingDoc)
		}
		frontendDocs[outputName] = FrontendSettingsDoc{
			Description: fmt.Sprintf("Frontend-specific settings for %s output driver. These should be placed in the 'frontendSettings' object as key-value pairs.", outputName),
			Settings:    settingDocs,
		}
	}

	doc := ChartConfigDocumentation{
		Description:      "This documentation teaches you how to create a ChartConfig JSON object. The ChartConfig is used to configure chart rendering with generic settings and frontend-specific settings.",
		Example:          example,
		Fields:           fields,
		FrontendSettings: frontendDocs,
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
