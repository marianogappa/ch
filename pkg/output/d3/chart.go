package d3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

type ChartType string

const (
	ChartTypeBar       ChartType = "bar"
	ChartTypePie       ChartType = "pie"
	ChartTypeScatter   ChartType = "scatter"
	ChartTypeHistogram ChartType = "histogram"
)

type Config struct {
	Title     string      `json:"title"`
	ChartType ChartType   `json:"chartType"`
	XLabel    string      `json:"xLabel"`
	YLabel    string      `json:"yLabel"`
	Color     string      `json:"color"`
	Other     interface{} `json:"other,omitempty"`
}

type Chart struct {
	Config Config
	Data   interface{}
}

func NewChart(config Config, data interface{}) *Chart {
	return &Chart{
		Config: config,
		Data:   data,
	}
}

func (c *Chart) Render(w io.Writer) error {
	// 1. Marshal data and config to JSON
	dataJSON, err := json.Marshal(c.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	configJSON, err := json.Marshal(c.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 2. Get the specific chart script template
	scriptTmpl, ok := Templates[string(c.Config.ChartType)]
	if !ok {
		return fmt.Errorf("unsupported chart type: %s", c.Config.ChartType)
	}

	// 3. Render the script part
	var scriptBuf bytes.Buffer
	if err := scriptTmpl.Execute(&scriptBuf, nil); err != nil {
		return fmt.Errorf("failed to execute script template: %w", err)
	}

	// 4. Render the base template with the script inserted
	baseTmpl := Templates["base"]
	return baseTmpl.Execute(w, struct {
		Title       string
		Data        template.JS
		Config      template.JS
		ChartScript template.JS
	}{
		Title:       c.Config.Title,
		Data:        template.JS(dataJSON),
		Config:      template.JS(configJSON),
		ChartScript: template.JS(scriptBuf.String()),
	})
}
