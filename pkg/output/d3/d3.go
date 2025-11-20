package d3

import (
	"flag"
	"fmt"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	ch.RegisterOutput(NewD3Output())
}

type D3Output struct{}

func NewD3Output() *D3Output {
	return &D3Output{}
}

func (o *D3Output) Name() string {
	return "d3"
}

type D3Config struct {
	Title     string
	ChartType string
	XLabel    string
	YLabel    string
	Color     string
}

func (o *D3Output) RegisterFlags(fs *flag.FlagSet) any {
	c := &D3Config{}
	fs.StringVar(&c.Title, "title", "D3 Chart", "Title of the chart")
	fs.StringVar(&c.ChartType, "chart-type", "bar", "Chart type: bar, pie, scatter, histogram")
	fs.StringVar(&c.XLabel, "x-label", "", "Label for X axis")
	fs.StringVar(&c.YLabel, "y-label", "", "Label for Y axis")
	fs.StringVar(&c.Color, "color", "", "Color of the chart elements (e.g. 'red', '#ff0000')")
	return c
}

func (o *D3Output) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   false,
		Interactive: true,
	}
}

func (o *D3Output) Render(rows <-chan ch.Row, config any) error {
	cfg, ok := config.(*D3Config)
	if !ok {
		return fmt.Errorf("invalid config type for D3Output")
	}

	var data []interface{}
	for row := range rows {
		// Basic mapping based on chart type
		// This is a simplified implementation. A real one would be more robust.
		switch cfg.ChartType {
		case "bar", "pie":
			if len(row.Strings) > 0 && len(row.Floats) > 0 {
				data = append(data, map[string]interface{}{
					"label": row.Strings[0],
					"value": row.Floats[0],
				})
			}
		case "scatter":
			if len(row.Floats) >= 2 {
				data = append(data, map[string]interface{}{
					"x": row.Floats[0],
					"y": row.Floats[1],
				})
			}
		case "histogram":
			if len(row.Floats) > 0 {
				data = append(data, map[string]interface{}{
					"value": row.Floats[0],
				})
			}
		default:
			// Default to bar-like structure if possible
			if len(row.Strings) > 0 && len(row.Floats) > 0 {
				data = append(data, map[string]interface{}{
					"label": row.Strings[0],
					"value": row.Floats[0],
				})
			}
		}
	}

	chartConfig := Config{
		Title:     cfg.Title,
		ChartType: ChartType(cfg.ChartType),
		XLabel:    cfg.XLabel,
		YLabel:    cfg.YLabel,
		Color:     cfg.Color,
	}

	c := NewChart(chartConfig, data)

	f, err := os.CreateTemp("", "ch-d3-*.html")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.Render(f); err != nil {
		return err
	}

	htmlPath := f.Name() + ".html"
	if err := os.Rename(f.Name(), htmlPath); err != nil {
		return err
	}

	fmt.Printf("Opening chart at %s\n", htmlPath)
	return openBrowser(htmlPath)
}

var openBrowser = open.Run
