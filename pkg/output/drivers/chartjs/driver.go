package chartjs

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	chdataset "github.com/marianogappa/ch/dataset"
	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/output"
	"github.com/marianogappa/ch/pkg/output/frontends/chartjs"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	ch.RegisterOutput(NewChartJSDriver())
}

type ChartJSDriver struct{}

func NewChartJSDriver() *ChartJSDriver {
	return &ChartJSDriver{}
}

func (d *ChartJSDriver) Name() string {
	return "chartjs"
}

// Ensure ChartJSDriver implements ChartOutput interface
var _ output.ChartOutput = (*ChartJSDriver)(nil)

type ChartJSConfig struct {
	output.ChartConfig
}

func (d *ChartJSDriver) RegisterFlags(fs *flag.FlagSet) any {
	c := &ChartJSConfig{
		ChartConfig: output.ChartConfig{
			ChartType:        output.ChartTypeLine,
			XScaleType:       output.ScaleTypeCategory,
			YScaleType:       output.ScaleTypeLinear,
			LegendDisplay:    true,
			LegendPosition:   output.LegendPositionTop,
			ColorScheme:      output.ColorSchemeDefault,
			FrontendSettings: make(map[string]any),
		},
	}

	fs.StringVar(&c.Title, "title", "", "Sets the title for the chart.")
	fs.StringVar(&c.XLabel, "x", "", "Sets the label for the x axis.")
	fs.StringVar(&c.YLabel, "y", "", "Sets the label for the y axis.")
	fs.BoolVar(&c.ZeroBased, "zero-based", false, "Makes y-axis begin at zero.")
	fs.StringVar((*string)(&c.ChartType), "chart-type", string(output.ChartTypeLine), "Chart type: line, bar, pie, doughnut, scatter, bubble, radar, polarArea.")
	fs.StringVar((*string)(&c.XScaleType), "x-scale", string(output.ScaleTypeCategory), "X-axis scale type: category, linear, logarithmic, time.")
	fs.StringVar((*string)(&c.YScaleType), "y-scale", string(output.ScaleTypeLinear), "Y-axis scale type: linear, logarithmic.")
	fs.BoolVar(&c.LegendDisplay, "legend", true, "Show legend.")
	fs.StringVar((*string)(&c.LegendPosition), "legend-position", string(output.LegendPositionTop), "Legend position: top, bottom, left, right, chartArea.")
	fs.StringVar((*string)(&c.ColorScheme), "color", string(output.ColorSchemeDefault), "Color scheme: default, legacy, gradient, custom.")
	fs.StringVar(&c.OutputPath, "output-path", "", "Output path: empty for browser, '-' for stdout, or file path.")

	return c
}

func (d *ChartJSDriver) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   false,
		Interactive: true,
	}
}

// RenderChart renders the chart using the generic ChartConfig
// Implements output.ChartOutput interface
func (d *ChartJSDriver) RenderChart(rows <-chan ch.Row, config *output.ChartConfig) error {
	// Convert generic ChartConfig to ChartJSConfig
	cfg := &ChartJSConfig{
		ChartConfig: *config,
	}
	return d.render(rows, cfg)
}

// GetFrontendSettings returns frontend-specific setting definitions
// Implements output.ChartOutput interface
func (d *ChartJSDriver) GetFrontendSettings() map[string]output.SettingDefinition {
	return map[string]output.SettingDefinition{
		"animation": {
			Name:         "animation",
			Type:         "bool",
			Description:  "Enable/disable animations",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
		"zoom": {
			Name:         "zoom",
			Type:         "bool",
			Description:  "Enable zoom functionality",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
		"pan": {
			Name:         "pan",
			Type:         "bool",
			Description:  "Enable pan functionality",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
		"responsive": {
			Name:         "responsive",
			Type:         "bool",
			Description:  "Make chart responsive",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
		"aspect-ratio": {
			Name:         "aspect-ratio",
			Type:         "number",
			Description:  "Set aspect ratio",
			DefaultValue: "",
			ValidValues:  nil,
			Required:     false,
		},
		"stacked": {
			Name:         "stacked",
			Type:         "bool",
			Description:  "Stack datasets",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
	}
}

func (d *ChartJSDriver) Render(rows <-chan ch.Row, config any) error {
	cfg, ok := config.(*ChartJSConfig)
	if !ok {
		return fmt.Errorf("invalid config type for ChartJSDriver")
	}
	return d.render(rows, cfg)
}

func (d *ChartJSDriver) render(rows <-chan ch.Row, cfg *ChartJSConfig) error {

	// Validate chart type is supported
	supportedChartTypes := []output.ChartType{
		output.ChartTypeLine,
		output.ChartTypeBar,
		output.ChartTypePie,
		output.ChartTypeDoughnut,
		output.ChartTypeScatter,
		output.ChartTypeBubble,
		output.ChartTypeRadar,
		output.ChartTypePolarArea,
	}
	chartTypeSupported := false
	for _, ct := range supportedChartTypes {
		if ct == cfg.ChartType {
			chartTypeSupported = true
			break
		}
	}
	if !chartTypeSupported {
		return fmt.Errorf("chart type %q is not supported by chartjs driver", cfg.ChartType)
	}

	// Buffer all rows to build a Dataset
	ds := &chdataset.Dataset{
		FSS: make([][]float64, 0),
		SSS: make([][]string, 0),
		TSS: make([][]time.Time, 0),
	}

	for row := range rows {
		ds.FSS = append(ds.FSS, row.Floats)
		ds.SSS = append(ds.SSS, row.Strings)

		ts := make([]time.Time, 0)
		for _, dStr := range row.DateTimes {
			t, _ := time.Parse("2006-01-02", dStr)
			ts = append(ts, t)
		}
		ds.TSS = append(ds.TSS, ts)
	}

	// If we have strings but no floats, we probably want to count frequencies
	if len(ds.FSS) > 0 && len(ds.FSS[0]) == 0 && len(ds.SSS) > 0 && len(ds.SSS[0]) > 0 {
		counts := make(map[string]float64)
		for _, ss := range ds.SSS {
			if len(ss) > 0 {
				counts[ss[0]]++
			}
		}

		// Rebuild dataset with counts
		ds.FSS = make([][]float64, 0, len(counts))
		ds.SSS = make([][]string, 0, len(counts))

		// Sort by count descending for better visualization
		type kv struct {
			Key   string
			Value float64
		}
		var ss []kv
		for k, v := range counts {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss {
			ds.FSS = append(ds.FSS, []float64{kv.Value})
			ds.SSS = append(ds.SSS, []string{kv.Key})
		}

		// Default to bar chart for frequency counts if not specified
		if cfg.ChartType == output.ChartTypeLine {
			cfg.ChartType = output.ChartTypeBar
		}
	}

	// Convert generic config to chartjs Chart
	chart := d.convertToChartJS(cfg, ds)

	// Create frontend instance
	frontend := chartjs.New(chart)

	// Determine output destination
	var outputWriter *os.File
	var htmlPath string
	var shouldOpenBrowser bool
	var tempFilePath string

	if cfg.OutputPath == "" {
		// Default behavior: create temp file and open browser
		f, err := os.CreateTemp("", "ch-chartjs-*.html")
		if err != nil {
			return err
		}
		outputWriter = f
		tempFilePath = f.Name()
		htmlPath = tempFilePath + ".html"
		shouldOpenBrowser = true
	} else if cfg.OutputPath == "-" {
		// Write to stdout
		outputWriter = os.Stdout
		htmlPath = ""
		shouldOpenBrowser = false
	} else {
		// Write to specified file path
		f, err := os.Create(cfg.OutputPath)
		if err != nil {
			return err
		}
		outputWriter = f
		htmlPath = cfg.OutputPath
		shouldOpenBrowser = false
	}

	// Build and render
	if err := frontend.Build(chartjs.OutputAll, outputWriter); err != nil {
		if outputWriter != os.Stdout {
			outputWriter.Close()
		}
		return err
	}

	// Close file before renaming if needed
	if outputWriter != os.Stdout {
		outputWriter.Close()
	}

	// If we created a temp file, rename it
	if shouldOpenBrowser && tempFilePath != "" {
		if err := os.Rename(tempFilePath, htmlPath); err != nil {
			return err
		}
		fmt.Printf("Opening chart at %s\n", htmlPath)
		return openBrowser(htmlPath)
	}

	if htmlPath != "" && !shouldOpenBrowser {
		fmt.Printf("Chart written to %s\n", htmlPath)
	}

	return nil
}

func (d *ChartJSDriver) convertToChartJS(cfg *ChartJSConfig, ds *chdataset.Dataset) *chartjs.Chart {
	// Convert chart type
	chartType := chartjs.ChartType(cfg.ChartType)

	// Build labels
	var labels []any
	if len(ds.SSS) > 0 && len(ds.SSS[0]) > 0 {
		labels = make([]any, len(ds.SSS))
		for i, ss := range ds.SSS {
			if len(ss) > 0 {
				labels[i] = ss[0]
			}
		}
	}

	// Build datasets
	datasets := make([]chartjs.Dataset, 0)
	if len(ds.FSS) > 0 {
		// For now, create one dataset per float column
		numCols := len(ds.FSS[0])
		for colIdx := 0; colIdx < numCols; colIdx++ {
			data := make([]any, len(ds.FSS))
			for rowIdx := 0; rowIdx < len(ds.FSS); rowIdx++ {
				if colIdx < len(ds.FSS[rowIdx]) {
					data[rowIdx] = ds.FSS[rowIdx][colIdx]
				}
			}

			dataset := chartjs.Dataset{
				Label: fmt.Sprintf("Series %d", colIdx+1),
				Data:  data,
			}

			// Apply color scheme
			d.applyColorScheme(&dataset, cfg, colIdx)

			datasets = append(datasets, dataset)
		}
	}

	data := &chartjs.Data{
		Labels:   labels,
		Datasets: datasets,
	}

	// Build options
	options := d.buildOptions(cfg)

	return &chartjs.Chart{
		Type:    chartType,
		Data:    data,
		Options: options,
	}
}

func (d *ChartJSDriver) buildOptions(cfg *ChartJSConfig) *chartjs.Options {
	opts := &chartjs.Options{}

	// Title
	if cfg.Title != "" {
		opts.Plugins = &chartjs.Plugins{}
		opts.Plugins.Title = &chartjs.Title{
			Display: boolPtr(true),
			Text:    stringPtr(cfg.Title),
		}
	}

	// Legend
	if cfg.LegendDisplay {
		if opts.Plugins == nil {
			opts.Plugins = &chartjs.Plugins{}
		}
		opts.Plugins.Legend = &chartjs.Legend{
			Display:  boolPtr(true),
			Position: stringPtr(string(cfg.LegendPosition)),
		}
	}

	// Scales
	opts.Scales = &chartjs.Scales{}

	// X-axis
	if cfg.XLabel != "" || cfg.XScaleType != output.ScaleTypeCategory {
		opts.Scales.X = &chartjs.Scale{}
		if cfg.XLabel != "" {
			opts.Scales.X.Title = &chartjs.ScaleTitle{
				Display: boolPtr(true),
				Text:    stringPtr(cfg.XLabel),
			}
		}
		if cfg.XScaleType == output.ScaleTypeTime {
			opts.Scales.X.Type = stringPtr("time")
		} else if cfg.XScaleType == output.ScaleTypeLogarithmic {
			opts.Scales.X.Type = stringPtr("logarithmic")
		} else if cfg.XScaleType == output.ScaleTypeLinear {
			opts.Scales.X.Type = stringPtr("linear")
		} else {
			opts.Scales.X.Type = stringPtr("category")
		}
	}

	// Y-axis
	if cfg.YLabel != "" || cfg.YScaleType != output.ScaleTypeLinear || cfg.ZeroBased {
		opts.Scales.Y = &chartjs.Scale{}
		if cfg.YLabel != "" {
			opts.Scales.Y.Title = &chartjs.ScaleTitle{
				Display: boolPtr(true),
				Text:    stringPtr(cfg.YLabel),
			}
		}
		if cfg.YScaleType == output.ScaleTypeLogarithmic {
			opts.Scales.Y.Type = stringPtr("logarithmic")
		} else {
			opts.Scales.Y.Type = stringPtr("linear")
		}
		if cfg.ZeroBased {
			opts.Scales.Y.BeginAtZero = boolPtr(true)
		}
	}

	// Apply frontend-specific settings
	d.applyFrontendSettings(opts, cfg.FrontendSettings)

	return opts
}

func (d *ChartJSDriver) applyColorScheme(dataset *chartjs.Dataset, cfg *ChartJSConfig, seriesIdx int) {
	switch cfg.ColorScheme {
	case output.ColorSchemeLegacy:
		// Use legacy colors
		legacyColors := []string{
			"rgba(220, 220, 220, 0.2)",
			"rgba(151, 187, 205, 0.2)",
			"rgba(151, 187, 205, 0.2)",
		}
		if seriesIdx < len(legacyColors) {
			dataset.BackgroundColor = legacyColors[seriesIdx]
		}
	case output.ColorSchemeGradient:
		// Use gradient (simplified - could be more sophisticated)
		dataset.BackgroundColor = "rgba(54, 162, 235, 0.2)"
	case output.ColorSchemeCustom:
		if len(cfg.Colors) > seriesIdx {
			dataset.BackgroundColor = cfg.Colors[seriesIdx]
		}
	default:
		// Default Chart.js colors
		dataset.BackgroundColor = nil // Let Chart.js use defaults
	}
}

func (d *ChartJSDriver) applyFrontendSettings(opts *chartjs.Options, settings map[string]any) {
	// Apply frontend-specific settings from the map
	// This allows users to pass Chart.js-specific options that aren't in the generic config
	// For example: animation duration, responsive settings, etc.
	// Implementation can be extended as needed
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

var openBrowser = open.Run
