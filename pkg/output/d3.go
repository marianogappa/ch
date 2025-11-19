package output

import (
	"flag"
	"fmt"
	"html/template"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
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
	Title string
}

func (o *D3Output) RegisterFlags(fs *flag.FlagSet) any {
	c := &D3Config{}
	fs.StringVar(&c.Title, "title", "D3 Chart", "Title of the chart")
	return c
}

func (o *D3Output) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   false,
		Interactive: true,
	}
}

func (o *D3Output) Render(rows <-chan ch.Row, config any) error {
	cfg, _ := config.(*D3Config)

	// Collect data
	var data []ch.Row
	for row := range rows {
		data = append(data, row)
	}

	// Frequency counting for string-only input
	if len(data) > 0 && len(data[0].Floats) == 0 && len(data[0].Strings) > 0 {
		counts := make(map[string]float64)
		for _, row := range data {
			if len(row.Strings) > 0 {
				counts[row.Strings[0]]++
			}
		}

		// Rebuild data with counts
		// We'll put the count in Floats[0]
		// And the label in Strings[0]
		var newData []ch.Row
		for k, v := range counts {
			newData = append(newData, ch.Row{
				Floats:  []float64{v},
				Strings: []string{k},
			})
		}
		data = newData
	}

	// Create HTML with D3
	// This is a very basic scatter plot example for now
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div id="chart"></div>
    <script>
        var data = {{.Data}};
        
        // Extract floats
        var points = data.map(function(d) { return d.Floats; });
        var labels = data.map(function(d) { return d.Strings ? d.Strings[0] : ""; });
        
        var width = 800, height = 600;
        var svg = d3.select("#chart").append("svg")
            .attr("width", width)
            .attr("height", height);
            
        svg.selectAll("circle")
            .data(points)
            .enter()
            .append("circle")
            .attr("cx", function(d, i) { 
                if (d && d.length >= 2) return d[0] * 10 + 50;
                return i * 50 + 50; // Use index for X if 1D
            }) 
            .attr("cy", function(d) { 
                if (d && d.length >= 2) return d[1] * 10 + 50;
                if (d && d.length >= 1) return 600 - (d[0] * 10 + 50); // Invert Y for 1D
                return 0;
            })
            .attr("r", 5)
            .append("title") // Simple tooltip
            .text(function(d, i) { return labels[i] + ": " + (d ? d[0] : ""); });
    </script>
</body>
</html>
`
	t, err := template.New("d3").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "ch-d3-*.html")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, struct {
		Title string
		Data  []ch.Row
	}{
		Title: cfg.Title,
		Data:  data,
	}); err != nil {
		return err
	}

	htmlPath := f.Name() + ".html"
	if err := os.Rename(f.Name(), htmlPath); err != nil {
		return err
	}

	fmt.Printf("Opening D3 chart at %s\n", htmlPath)
	return openBrowser(htmlPath)
}
