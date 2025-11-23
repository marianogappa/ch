package d3

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

// D3 is the main struct for building D3.js charts
type D3 struct {
	chart    *Chart
	canvasID string
}

// New creates a new D3 instance
func New(chart *Chart) *D3 {
	return &D3{
		chart:    chart,
		canvasID: "chart",
	}
}

// SetCanvasID sets the canvas element ID (default: "chart")
func (d *D3) SetCanvasID(id string) {
	d.canvasID = id
}

// OutputMode parameterizes the output of Build()
type OutputMode int

const (
	OutputAll          OutputMode = iota // Outputs a complete HTML file
	OutputHTMLHeader                     // Outputs an HTML header, dependencies and opening <body> tag
	OutputDependencies                   // Outputs only <script> tags with D3 dependencies
	OutputChart                          // Outputs the JS code that renders the chart
	OutputHTMLFooter                     // Outputs </body></html>
)

// MustBuild prepares the chart and executes the text template with it. Fatals if there's a problem
func (d *D3) MustBuild(om OutputMode, w io.Writer) {
	if err := d.Build(om, w); err != nil {
		log.Fatal(err)
	}
}

// Build prepares the chart and executes the text template with it. Returns an error if there's a problem
func (d *D3) Build(om OutputMode, w io.Writer) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in D3.Build: %v", r)
		}
	}()

	switch om {
	case OutputDependencies:
		return tplToWriter(tplD3, "", w)
	case OutputHTMLHeader:
		bb := bytes.Buffer{}
		if err := tplToWriter(tplD3, "", &bb); err != nil {
			return err
		}
		return tplToWriter(tplHTMLHeader, bb.String(), w)
	case OutputChart:
		chartScript, err := d.prepareChartScript()
		if err != nil {
			return err
		}
		return tplToWriter(tplChartScript, chartScript, w)
	case OutputHTMLFooter:
		return tplToWriter(tplHTMLFooter, "", w)
	case OutputAll:
		deps := bytes.Buffer{}
		if err := tplToWriter(tplD3, "", &deps); err != nil {
			return err
		}
		if err := tplToWriter(tplHTMLHeader, deps.String(), w); err != nil {
			return err
		}

		chartScript, err := d.prepareChartScript()
		if err != nil {
			return err
		}

		dataJSON, err := json.Marshal(d.chart.Data)
		if err != nil {
			return err
		}

		configJSON, err := json.Marshal(d.chart.Config)
		if err != nil {
			return err
		}

		return tplToWriter(tplChartDivScript, map[string]string{
			"CanvasID":    d.canvasID,
			"ChartScript": chartScript,
			"Data":        string(dataJSON),
			"Config":      string(configJSON),
			"Title":       d.chart.Config.Title,
		}, w)
	}
	return nil
}

func (d *D3) prepareChartScript() (string, error) {
	// Get the specific chart script template based on chart type
	scriptTmpl, ok := chartTemplates[string(d.chart.Type)]
	if !ok {
		return "", nil // Return empty if no specific template
	}

	var scriptBuf bytes.Buffer
	if err := scriptTmpl.Execute(&scriptBuf, nil); err != nil {
		return "", err
	}

	return scriptBuf.String(), nil
}
