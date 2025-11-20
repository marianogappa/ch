package chartjs

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	chdataset "github.com/marianogappa/ch/dataset"
	"github.com/marianogappa/ch/pkg/ch"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	ch.RegisterOutput(NewChartJSOutput())
}

type ChartJSOutput struct{}

func NewChartJSOutput() *ChartJSOutput {
	return &ChartJSOutput{}
}

func (o *ChartJSOutput) Name() string {
	return "chartjs"
}

type ChartJSConfig struct {
	Title     string
	XLabel    string
	YLabel    string
	ZeroBased bool
	ChartType string
	ScaleType string
	ColorType string
}

func (o *ChartJSOutput) RegisterFlags(fs *flag.FlagSet) any {
	c := &ChartJSConfig{}
	fs.StringVar(&c.Title, "title", "", "Sets the title for the chart.")
	fs.StringVar(&c.XLabel, "x", "", "Sets the label for the x axis.")
	fs.StringVar(&c.YLabel, "y", "", "Sets the label for the y axis.")
	fs.BoolVar(&c.ZeroBased, "zero-based", false, "Makes y-axis begin at zero.")
	fs.StringVar(&c.ChartType, "chart-type", "line", "Chart type: line, bar, pie, scatter.") // Renamed from implicit arg
	fs.StringVar(&c.ScaleType, "scale", "linear", "Scale type: linear, logarithmic.")
	fs.StringVar(&c.ColorType, "color", "default", "Color type: default, legacy, gradient.")
	return c
}

func (o *ChartJSOutput) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   false, // ChartJS implementation here is not streaming
		Interactive: true,
	}
}

func (o *ChartJSOutput) Render(rows <-chan ch.Row, config any) error {
	cfg, ok := config.(*ChartJSConfig)
	if !ok {
		return fmt.Errorf("invalid config type for ChartJSOutput")
	}

	// Buffer all rows to build a Dataset
	// This is a bridge between streaming architecture and legacy Dataset struct
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
			// Try to parse with a few common formats if possible, or just use a default
			// For now, we rely on the fact that if it was parsed as DateTime, it should be parseable.
			// But we don't have the format here.
			// TODO: Pass format or parsed time in Row.
			t, _ := time.Parse("2006-01-02", dStr) // minimal assumption
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
		if cfg.ChartType == "line" { // "line" is the default in RegisterFlags
			cfg.ChartType = "bar"
		}
	}

	cOpts := Options{
		Title:     cfg.Title,
		ScaleType: NewScaleType(cfg.ScaleType),
		XLabel:    cfg.XLabel,
		YLabel:    cfg.YLabel,
		ZeroBased: cfg.ZeroBased,
		ColorType: NewColorType(cfg.ColorType),
	}

	// Now use the legacy chartjs package
	c := New(
		NewChartType(cfg.ChartType),
		*ds,
		cOpts,
	)

	// We need to handle the temp file creation here or inside chartjs?
	// The original main.go did it.
	// I'll replicate that logic here.
	// I need to copy `tmpfile.go` logic or reimplement it.
	// I'll just use `os.CreateTemp`.

	f, err := os.CreateTemp("", "ch-*.html")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.Build(OutputAll, f); err != nil {
		return err
	}

	// Rename to .html to ensure browser opens it correctly
	htmlPath := f.Name() + ".html"
	if err := os.Rename(f.Name(), htmlPath); err != nil {
		return err
	}

	fmt.Printf("Opening chart at %s\n", htmlPath)
	fmt.Printf("Opening chart at %s\n", htmlPath)
	return openBrowser(htmlPath)
}

var openBrowser = open.Run
