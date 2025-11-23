package d3

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	chdataset "github.com/marianogappa/ch/dataset"
	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/output"
	"github.com/marianogappa/ch/pkg/output/frontends/d3"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	ch.RegisterOutput(NewD3Driver())
}

type D3Driver struct{}

func NewD3Driver() *D3Driver {
	return &D3Driver{}
}

func (d *D3Driver) Name() string {
	return "d3"
}

// Ensure D3Driver implements ChartOutput interface
var _ output.ChartOutput = (*D3Driver)(nil)

type D3Config struct {
	output.ChartConfig
}

func (d *D3Driver) RegisterFlags(fs *flag.FlagSet) any {
	c := &D3Config{
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
	fs.StringVar((*string)(&c.ChartType), "chart-type", string(output.ChartTypeLine), "Chart type: line, bar, pie, scatter, histogram.")
	fs.StringVar((*string)(&c.XScaleType), "x-scale", string(output.ScaleTypeCategory), "X-axis scale type: category, linear, logarithmic, time.")
	fs.StringVar((*string)(&c.YScaleType), "y-scale", string(output.ScaleTypeLinear), "Y-axis scale type: linear, logarithmic.")
	fs.BoolVar(&c.LegendDisplay, "legend", true, "Show legend.")
	fs.StringVar((*string)(&c.LegendPosition), "legend-position", string(output.LegendPositionTop), "Legend position: top, bottom, left, right, chartArea.")
	fs.StringVar((*string)(&c.ColorScheme), "color", string(output.ColorSchemeDefault), "Color scheme: default, legacy, gradient, custom.")
	fs.StringVar(&c.OutputPath, "output-path", "", "Output path: empty for browser, '-' for stdout, or file path.")

	return c
}

func (d *D3Driver) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   false,
		Interactive: true,
	}
}

// RenderChart renders the chart using the generic ChartConfig
// Implements output.ChartOutput interface
func (d *D3Driver) RenderChart(rows <-chan ch.Row, config *output.ChartConfig) error {
	// Convert generic ChartConfig to D3Config
	cfg := &D3Config{
		ChartConfig: *config,
	}
	return d.render(rows, cfg)
}

// GetFrontendSettings returns frontend-specific setting definitions
// Implements output.ChartOutput interface
func (d *D3Driver) GetFrontendSettings() map[string]output.SettingDefinition {
	return map[string]output.SettingDefinition{
		"animation": {
			Name:         "animation",
			Type:         "bool",
			Description:  "Enable/disable animations",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
		"tooltip-customization": {
			Name:         "tooltip-customization",
			Type:         "string",
			Description:  "Customize tooltip behavior",
			DefaultValue: "",
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
		"interactive": {
			Name:         "interactive",
			Type:         "bool",
			Description:  "Enable interactive features",
			DefaultValue: "false",
			ValidValues:  nil,
			Required:     false,
		},
	}
}

func (d *D3Driver) Render(rows <-chan ch.Row, config any) error {
	cfg, ok := config.(*D3Config)
	if !ok {
		return fmt.Errorf("invalid config type for D3Driver")
	}
	return d.render(rows, cfg)
}

func (d *D3Driver) render(rows <-chan ch.Row, cfg *D3Config) error {

	// Validate chart type is supported
	supportedChartTypes := []output.ChartType{
		output.ChartTypeLine,
		output.ChartTypeBar,
		output.ChartTypePie,
		output.ChartTypeScatter,
		output.ChartTypeHistogram,
	}
	chartTypeSupported := false
	for _, ct := range supportedChartTypes {
		if ct == cfg.ChartType {
			chartTypeSupported = true
			break
		}
	}
	if !chartTypeSupported {
		return fmt.Errorf("chart type %q is not supported by d3 driver", cfg.ChartType)
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

	// Convert generic config to d3 Chart
	chart := d.convertToD3(cfg, ds)

	// Create frontend instance
	frontend := d3.New(chart)

	// Determine output destination
	var outputWriter *os.File
	var htmlPath string
	var shouldOpenBrowser bool
	var tempFilePath string

	if cfg.OutputPath == "" {
		// Default behavior: create temp file and open browser
		f, err := os.CreateTemp("", "ch-d3-*.html")
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
	if err := frontend.Build(d3.OutputAll, outputWriter); err != nil {
		outputWriter.Close()
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

func (d *D3Driver) convertToD3(cfg *D3Config, ds *chdataset.Dataset) *d3.Chart {
	// Convert chart type
	chartType := d3.ChartType(cfg.ChartType)

	// Build data based on chart type
	var data []any
	switch cfg.ChartType {
	case output.ChartTypeBar, output.ChartTypePie:
		// For bar and pie charts, use label/value format
		if len(ds.SSS) > 0 && len(ds.FSS) > 0 {
			data = make([]any, len(ds.SSS))
			for i := range ds.SSS {
				item := make(map[string]any)
				if len(ds.SSS[i]) > 0 {
					item["label"] = ds.SSS[i][0]
				}
				if len(ds.FSS[i]) > 0 {
					item["value"] = ds.FSS[i][0]
				}
				data[i] = item
			}
		}
	case output.ChartTypeScatter:
		// For scatter plots, use x/y format
		if len(ds.FSS) > 0 {
			data = make([]any, len(ds.FSS))
			for i := range ds.FSS {
				item := make(map[string]any)
				if len(ds.FSS[i]) >= 2 {
					item["x"] = ds.FSS[i][0]
					item["y"] = ds.FSS[i][1]
				} else if len(ds.FSS[i]) == 1 {
					// If only one float, use index as x
					item["x"] = float64(i)
					item["y"] = ds.FSS[i][0]
				}
				data[i] = item
			}
		}
	case output.ChartTypeHistogram:
		// For histograms, use value format
		if len(ds.FSS) > 0 {
			data = make([]any, len(ds.FSS))
			for i := range ds.FSS {
				item := make(map[string]any)
				if len(ds.FSS[i]) > 0 {
					item["value"] = ds.FSS[i][0]
				}
				data[i] = item
			}
		}
	case output.ChartTypeLine:
		// For line charts, use label/value format (similar to bar)
		if len(ds.SSS) > 0 && len(ds.FSS) > 0 {
			data = make([]any, len(ds.SSS))
			for i := range ds.SSS {
				item := make(map[string]any)
				if len(ds.SSS[i]) > 0 {
					item["label"] = ds.SSS[i][0]
				}
				if len(ds.FSS[i]) > 0 {
					item["value"] = ds.FSS[i][0]
				}
				data[i] = item
			}
		} else if len(ds.FSS) > 0 {
			// If no strings, use index as label
			data = make([]any, len(ds.FSS))
			for i := range ds.FSS {
				item := make(map[string]any)
				item["label"] = fmt.Sprintf("%d", i)
				if len(ds.FSS[i]) > 0 {
					item["value"] = ds.FSS[i][0]
				}
				data[i] = item
			}
		}
	}

	// Build config
	config := &d3.Config{
		Title:     cfg.Title,
		ChartType: chartType,
		XLabel:    cfg.XLabel,
		YLabel:    cfg.YLabel,
		Color:     d.getColor(cfg),
	}

	// Apply color scheme
	if len(cfg.Colors) > 0 {
		config.Colors = cfg.Colors
	}

	return &d3.Chart{
		Type:   chartType,
		Data:   data,
		Config: config,
	}
}

func (d *D3Driver) getColor(cfg *D3Config) string {
	switch cfg.ColorScheme {
	case output.ColorSchemeLegacy:
		return "#97BBCD"
	case output.ColorSchemeGradient:
		return "#3692EB"
	case output.ColorSchemeCustom:
		if len(cfg.Colors) > 0 {
			return cfg.Colors[0]
		}
		return ""
	default:
		return "steelblue"
	}
}

var openBrowser = open.Run
