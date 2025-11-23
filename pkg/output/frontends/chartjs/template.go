package chartjs

import (
	"io"
	"log"
	"text/template"
)

type tpl int

const (
	tplChartDivScript tpl = iota
	tplChartObject
	tplHTMLHeader
	tplHTMLFooter
	tplChartJS
)

func tplToWriter(t tpl, s interface{}, w io.Writer) error {
	return templates[t].Execute(w, s)
}

func init() {
	var err error
	for t, ts := range templateStrings {
		if templates[t], err = template.New("").Parse(ts); err != nil {
			log.Fatalf("chartjs.templates: error parsing text template: %v", err)
		}
	}
}

var templates = map[tpl]*template.Template{}

var templateStrings = map[tpl]string{
	tplChartDivScript: `<div style="height:90vh; width:90vw;">
<canvas id="{{.CanvasID}}"></canvas>
</div>
<script>
const ctx{{.CanvasID}} = document.getElementById("{{.CanvasID}}");
const chart{{.CanvasID}} = new Chart(ctx{{.CanvasID}}, {{.ChartConfig}});
</script>`,

	tplChartObject: `{{.}}`,

	tplHTMLHeader: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chart.js 4.0 Chart</title>
    <style>
        /* CSS Reset */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            padding: 20px;
            background-color: #f5f5f5;
        }
        #chart {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
        }
    </style>
{{.}}
</head>
<body>`,

	tplHTMLFooter: `
</body>
</html>`,

	tplChartJS: `<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>`,
}
