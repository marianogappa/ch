package chartjs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// marshalChartToJS converts a Chart struct to JavaScript object notation
// This handles special cases like time.Time, JavaScript functions, etc.
func marshalChartToJS(chart *Chart) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	// Type
	buf.WriteString(fmt.Sprintf("  type: '%s',\n", chart.Type))

	// Data
	if chart.Data != nil {
		buf.WriteString("  data: ")
		dataStr, err := marshalDataToJS(chart.Data)
		if err != nil {
			return "", err
		}
		buf.WriteString(dataStr)
		buf.WriteString(",\n")
	}

	// Options
	if chart.Options != nil {
		buf.WriteString("  options: ")
		optionsStr, err := marshalOptionsToJS(chart.Options)
		if err != nil {
			return "", err
		}
		buf.WriteString(optionsStr)
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String(), nil
}

func marshalDataToJS(data *Data) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	// Labels
	if len(data.Labels) > 0 {
		buf.WriteString("    labels: ")
		labelsStr, err := marshalLabelsToJS(data.Labels)
		if err != nil {
			return "", err
		}
		buf.WriteString(labelsStr)
		buf.WriteString(",\n")
	}

	// Datasets
	if len(data.Datasets) > 0 {
		buf.WriteString("    datasets: [\n")
		for i, ds := range data.Datasets {
			if i > 0 {
				buf.WriteString(",\n")
			}
			dsStr, err := marshalDatasetToJS(ds)
			if err != nil {
				return "", err
			}
			buf.WriteString("      " + strings.ReplaceAll(dsStr, "\n", "\n      "))
		}
		buf.WriteString("\n    ]\n")
	}

	buf.WriteString("  }")
	return buf.String(), nil
}

func marshalLabelsToJS(labels []any) (string, error) {
	if len(labels) == 0 {
		return "[]", nil
	}

	var parts []string
	for _, label := range labels {
		switch v := label.(type) {
		case string:
			parts = append(parts, fmt.Sprintf("'%s'", escapeJSString(v)))
		case time.Time:
			parts = append(parts, fmt.Sprintf("'%s'", v.Format(time.RFC3339)))
		case int, int8, int16, int32, int64:
			parts = append(parts, fmt.Sprintf("%v", v))
		case uint, uint8, uint16, uint32, uint64:
			parts = append(parts, fmt.Sprintf("%v", v))
		case float32, float64:
			parts = append(parts, fmt.Sprintf("%v", v))
		default:
			// Fallback to JSON marshaling
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			parts = append(parts, string(jsonBytes))
		}
	}
	return "[" + strings.Join(parts, ", ") + "]", nil
}

func marshalDatasetToJS(ds Dataset) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	parts := []string{}

	if ds.Label != "" {
		parts = append(parts, fmt.Sprintf("      label: '%s'", escapeJSString(ds.Label)))
	}

	// Data
	if len(ds.Data) > 0 {
		dataStr, err := marshalDataArrayToJS(ds.Data)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("      data: %s", dataStr))
	}

	// BackgroundColor
	if ds.BackgroundColor != nil {
		bgStr, err := marshalColorToJS(ds.BackgroundColor)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("      backgroundColor: %s", bgStr))
	}

	// BorderColor
	if ds.BorderColor != nil {
		borderStr, err := marshalColorToJS(ds.BorderColor)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("      borderColor: %s", borderStr))
	}

	// Other optional fields
	if ds.BorderWidth != nil {
		parts = append(parts, fmt.Sprintf("      borderWidth: %d", *ds.BorderWidth))
	}
	if ds.Fill != nil {
		parts = append(parts, fmt.Sprintf("      fill: %v", *ds.Fill))
	}
	if ds.Tension != nil {
		parts = append(parts, fmt.Sprintf("      tension: %g", *ds.Tension))
	}
	if ds.PointRadius != nil {
		parts = append(parts, fmt.Sprintf("      pointRadius: %g", *ds.PointRadius))
	}
	if ds.Stepped != nil {
		parts = append(parts, fmt.Sprintf("      stepped: %v", *ds.Stepped))
	}
	if ds.BarPercentage != nil {
		parts = append(parts, fmt.Sprintf("      barPercentage: %g", *ds.BarPercentage))
	}
	if ds.CategoryPercentage != nil {
		parts = append(parts, fmt.Sprintf("      categoryPercentage: %g", *ds.CategoryPercentage))
	}
	if ds.Offset != nil {
		parts = append(parts, fmt.Sprintf("      offset: %g", *ds.Offset))
	}
	if ds.Radius != nil {
		parts = append(parts, fmt.Sprintf("      radius: %g", *ds.Radius))
	}
	if ds.Stack != nil {
		parts = append(parts, fmt.Sprintf("      stack: '%s'", escapeJSString(*ds.Stack)))
	}
	if ds.IndexAxis != nil {
		parts = append(parts, fmt.Sprintf("      indexAxis: '%s'", *ds.IndexAxis))
	}

	buf.WriteString(strings.Join(parts, ",\n"))
	buf.WriteString("\n    }")
	return buf.String(), nil
}

func marshalDataArrayToJS(data []any) (string, error) {
	if len(data) == 0 {
		return "[]", nil
	}

	var parts []string
	for _, item := range data {
		switch v := item.(type) {
		case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			parts = append(parts, fmt.Sprintf("%v", v))
		case DataPoint:
			dpStr := fmt.Sprintf("{x: %g, y: %g", v.X, v.Y)
			if v.R != nil {
				dpStr += fmt.Sprintf(", r: %g", *v.R)
			}
			dpStr += "}"
			parts = append(parts, dpStr)
		case TimeDataPoint:
			dpStr := fmt.Sprintf("{x: '%s', y: %g", v.X.Format(time.RFC3339), v.Y)
			if v.R != nil {
				dpStr += fmt.Sprintf(", r: %g", *v.R)
			}
			dpStr += "}"
			parts = append(parts, dpStr)
		case time.Time:
			parts = append(parts, fmt.Sprintf("'%s'", v.Format(time.RFC3339)))
		case string:
			parts = append(parts, fmt.Sprintf("'%s'", escapeJSString(v)))
		default:
			// Try to marshal as JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("failed to marshal data point: %w", err)
			}
			parts = append(parts, string(jsonBytes))
		}
	}
	return "[" + strings.Join(parts, ", ") + "]", nil
}

func marshalColorToJS(color any) (string, error) {
	if color == nil {
		return "null", nil
	}

	// Check if it's a string (JavaScript function)
	if str, ok := color.(string); ok {
		// Check if it looks like a function
		if strings.HasPrefix(strings.TrimSpace(str), "function") ||
			strings.HasPrefix(strings.TrimSpace(str), "(") {
			return str, nil
		}
		// Otherwise it's a color string
		return fmt.Sprintf("'%s'", escapeJSString(str)), nil
	}

	// Check if it's a slice of strings
	rv := reflect.ValueOf(color)
	if rv.Kind() == reflect.Slice {
		var parts []string
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			if str, ok := item.(string); ok {
				parts = append(parts, fmt.Sprintf("'%s'", escapeJSString(str)))
			} else {
				jsonBytes, err := json.Marshal(item)
				if err != nil {
					return "", err
				}
				parts = append(parts, string(jsonBytes))
			}
		}
		return "[" + strings.Join(parts, ", ") + "]", nil
	}

	// Fallback to JSON
	jsonBytes, err := json.Marshal(color)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func marshalOptionsToJS(options *Options) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	parts := []string{}

	if options.Responsive != nil {
		parts = append(parts, fmt.Sprintf("    responsive: %v", *options.Responsive))
	}
	if options.MaintainAspectRatio != nil {
		parts = append(parts, fmt.Sprintf("    maintainAspectRatio: %v", *options.MaintainAspectRatio))
	}
	if options.AspectRatio != nil {
		parts = append(parts, fmt.Sprintf("    aspectRatio: %g", *options.AspectRatio))
	}
	if options.IndexAxis != nil {
		parts = append(parts, fmt.Sprintf("    indexAxis: '%s'", *options.IndexAxis))
	}

	// Scales
	if options.Scales != nil {
		scalesStr, err := marshalScalesToJS(options.Scales)
		if err != nil {
			return "", err
		}
		parts = append(parts, "    scales: "+strings.ReplaceAll(scalesStr, "\n", "\n    "))
	}

	// Plugins
	if options.Plugins != nil {
		pluginsStr, err := marshalPluginsToJS(options.Plugins)
		if err != nil {
			return "", err
		}
		parts = append(parts, "    plugins: "+strings.ReplaceAll(pluginsStr, "\n", "\n    "))
	}

	// Elements
	if options.Elements != nil {
		elementsStr, err := marshalElementsToJS(options.Elements)
		if err != nil {
			return "", err
		}
		parts = append(parts, "    elements: "+strings.ReplaceAll(elementsStr, "\n", "\n    "))
	}

	// Animation
	if options.Animation != nil {
		animStr, err := marshalAnimationToJS(options.Animation)
		if err != nil {
			return "", err
		}
		parts = append(parts, "    animation: "+strings.ReplaceAll(animStr, "\n", "\n    "))
	}

	// Interaction
	if options.Interaction != nil {
		interactionStr, err := marshalInteractionToJS(options.Interaction)
		if err != nil {
			return "", err
		}
		parts = append(parts, "    interaction: "+strings.ReplaceAll(interactionStr, "\n", "\n    "))
	}

	buf.WriteString(strings.Join(parts, ",\n"))
	buf.WriteString("\n  }")
	return buf.String(), nil
}

func marshalScalesToJS(scales *Scales) (string, error) {
	var parts []string

	if scales.X != nil {
		xStr, err := marshalScaleToJS("x", scales.X)
		if err != nil {
			return "", err
		}
		parts = append(parts, xStr)
	}
	if scales.Y != nil {
		yStr, err := marshalScaleToJS("y", scales.Y)
		if err != nil {
			return "", err
		}
		parts = append(parts, yStr)
	}
	if scales.R != nil {
		rStr, err := marshalScaleToJS("r", scales.R)
		if err != nil {
			return "", err
		}
		parts = append(parts, rStr)
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func marshalScaleToJS(name string, scale *Scale) (string, error) {
	var parts []string

	if scale.Type != nil {
		parts = append(parts, fmt.Sprintf("%s: {\n        type: '%s'", name, *scale.Type))
	} else {
		parts = append(parts, fmt.Sprintf("%s: {", name))
	}

	scaleParts := []string{}

	if scale.Display != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        display: %v", *scale.Display))
	}
	if scale.Position != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        position: '%s'", *scale.Position))
	}
	if scale.BeginAtZero != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        beginAtZero: %v", *scale.BeginAtZero))
	}
	if scale.Reverse != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        reverse: %v", *scale.Reverse))
	}
	if scale.Stacked != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        stacked: %v", *scale.Stacked))
	}
	if scale.Min != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        min: %g", *scale.Min))
	}
	if scale.Max != nil {
		scaleParts = append(scaleParts, fmt.Sprintf("        max: %g", *scale.Max))
	}

	// Title
	if scale.Title != nil {
		titleParts := []string{}
		if scale.Title.Display != nil {
			titleParts = append(titleParts, fmt.Sprintf("          display: %v", *scale.Title.Display))
		}
		if scale.Title.Text != nil {
			titleParts = append(titleParts, fmt.Sprintf("          text: '%s'", escapeJSString(*scale.Title.Text)))
		}
		if len(titleParts) > 0 {
			scaleParts = append(scaleParts, "        title: {\n"+strings.Join(titleParts, ",\n")+"\n        }")
		}
	}

	// Ticks
	if scale.Ticks != nil {
		tickParts := []string{}
		if scale.Ticks.Display != nil {
			tickParts = append(tickParts, fmt.Sprintf("          display: %v", *scale.Ticks.Display))
		}
		if scale.Ticks.Color != nil {
			tickParts = append(tickParts, fmt.Sprintf("          color: '%s'", *scale.Ticks.Color))
		}
		if scale.Ticks.Callback != nil {
			// JavaScript function - output as-is
			tickParts = append(tickParts, fmt.Sprintf("          callback: %s", *scale.Ticks.Callback))
		}
		if len(tickParts) > 0 {
			scaleParts = append(scaleParts, "        ticks: {\n"+strings.Join(tickParts, ",\n")+"\n        }")
		}
	}

	// Grid
	if scale.Grid != nil {
		gridParts := []string{}
		if scale.Grid.Display != nil {
			gridParts = append(gridParts, fmt.Sprintf("          display: %v", *scale.Grid.Display))
		}
		if scale.Grid.Color != nil {
			gridParts = append(gridParts, fmt.Sprintf("          color: '%s'", *scale.Grid.Color))
		}
		if len(gridParts) > 0 {
			scaleParts = append(scaleParts, "        grid: {\n"+strings.Join(gridParts, ",\n")+"\n        }")
		}
	}

	// Time scale
	if scale.Time != nil {
		timeParts := []string{}
		if scale.Time.Unit != nil {
			timeParts = append(timeParts, fmt.Sprintf("          unit: '%s'", *scale.Time.Unit))
		}
		if scale.Time.Parser != nil {
			timeParts = append(timeParts, fmt.Sprintf("          parser: '%s'", escapeJSString(*scale.Time.Parser)))
		}
		if len(timeParts) > 0 {
			scaleParts = append(scaleParts, "        time: {\n"+strings.Join(timeParts, ",\n")+"\n        }")
		}
	}

	if len(scaleParts) > 0 {
		parts[0] += ",\n" + strings.Join(scaleParts, ",\n")
	}
	parts[0] += "\n      }"

	return parts[0], nil
}

func marshalPluginsToJS(plugins *Plugins) (string, error) {
	var parts []string

	if plugins.Legend != nil {
		legendStr, err := marshalLegendToJS(plugins.Legend)
		if err != nil {
			return "", err
		}
		parts = append(parts, "legend: "+strings.ReplaceAll(legendStr, "\n", "\n    "))
	}

	if plugins.Title != nil {
		titleStr, err := marshalTitleToJS(plugins.Title)
		if err != nil {
			return "", err
		}
		parts = append(parts, "title: "+strings.ReplaceAll(titleStr, "\n", "\n    "))
	}

	if plugins.Tooltip != nil {
		tooltipStr, err := marshalTooltipToJS(plugins.Tooltip)
		if err != nil {
			return "", err
		}
		parts = append(parts, "tooltip: "+strings.ReplaceAll(tooltipStr, "\n", "\n    "))
	}

	return "{\n    " + strings.Join(parts, ",\n    ") + "\n  }", nil
}

func marshalLegendToJS(legend *Legend) (string, error) {
	var parts []string

	if legend.Display != nil {
		parts = append(parts, fmt.Sprintf("display: %v", *legend.Display))
	}
	if legend.Position != nil {
		parts = append(parts, fmt.Sprintf("position: '%s'", *legend.Position))
	}
	if legend.Align != nil {
		parts = append(parts, fmt.Sprintf("align: '%s'", *legend.Align))
	}
	if legend.OnClick != nil {
		parts = append(parts, fmt.Sprintf("onClick: %s", *legend.OnClick))
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func marshalTitleToJS(title *Title) (string, error) {
	var parts []string

	if title.Display != nil {
		parts = append(parts, fmt.Sprintf("display: %v", *title.Display))
	}
	if title.Text != nil {
		parts = append(parts, fmt.Sprintf("text: '%s'", escapeJSString(*title.Text)))
	}
	if title.Position != nil {
		parts = append(parts, fmt.Sprintf("position: '%s'", *title.Position))
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func marshalTooltipToJS(tooltip *Tooltip) (string, error) {
	var parts []string

	if tooltip.Enabled != nil {
		parts = append(parts, fmt.Sprintf("enabled: %v", *tooltip.Enabled))
	}
	if tooltip.Mode != nil {
		parts = append(parts, fmt.Sprintf("mode: '%s'", *tooltip.Mode))
	}
	if tooltip.Callbacks != nil && tooltip.Callbacks.Label != nil {
		parts = append(parts, fmt.Sprintf("callbacks: {\n        label: %s\n      }", *tooltip.Callbacks.Label))
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func marshalElementsToJS(elements *Elements) (string, error) {
	var parts []string

	if elements.Line != nil {
		lineParts := []string{}
		if elements.Line.Tension != nil {
			lineParts = append(lineParts, fmt.Sprintf("        tension: %g", *elements.Line.Tension))
		}
		if elements.Line.BorderWidth != nil {
			lineParts = append(lineParts, fmt.Sprintf("        borderWidth: %g", *elements.Line.BorderWidth))
		}
		if len(lineParts) > 0 {
			parts = append(parts, "line: {\n"+strings.Join(lineParts, ",\n")+"\n      }")
		}
	}

	if len(parts) == 0 {
		return "{}", nil
	}
	return "{\n    " + strings.Join(parts, ",\n    ") + "\n  }", nil
}

func marshalAnimationToJS(animation *Animation) (string, error) {
	var parts []string

	if animation.Duration != nil {
		parts = append(parts, fmt.Sprintf("duration: %d", *animation.Duration))
	}
	if animation.Easing != nil {
		parts = append(parts, fmt.Sprintf("easing: '%s'", *animation.Easing))
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func marshalInteractionToJS(interaction *Interaction) (string, error) {
	var parts []string

	if interaction.Mode != nil {
		parts = append(parts, fmt.Sprintf("mode: '%s'", *interaction.Mode))
	}
	if interaction.Intersect != nil {
		parts = append(parts, fmt.Sprintf("intersect: %v", *interaction.Intersect))
	}

	return "{\n      " + strings.Join(parts, ",\n      ") + "\n    }", nil
}

func escapeJSString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
