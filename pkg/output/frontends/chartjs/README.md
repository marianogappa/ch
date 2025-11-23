# Chart.js Go Package

A comprehensive Go package that provides a complete API wrapper for Chart.js, allowing you to build sophisticated charts programmatically and generate HTML/JavaScript output.

## Features

- **Complete Chart.js API Coverage**: All chart types, scales, plugins, and options are supported
- **Type-Safe Configuration**: Full Go struct definitions for all Chart.js options
- **Flexible Data Input**: Support for various data formats including time series, scatter plots, and more
- **Helper Functions**: Convenient helper functions for common configurations
- **Template-Based Output**: Generates clean HTML with embedded Chart.js code

## Chart Types Supported

- Line charts
- Bar charts
- Pie charts
- Doughnut charts
- Polar area charts
- Radar charts
- Bubble charts
- Scatter charts

## Quick Start

### Basic Line Chart

```go
package main

import (
    "os"
    "github.com/marianogappa/ch/pkg/output/frontends/chartjs"
)

func main() {
    chart := chartjs.NewLineChart(
        []string{"Jan", "Feb", "Mar", "Apr", "May"},
        []chartjs.LineDataset{
            {
                Label:           "Sales",
                Data:            []float64{10, 19, 15, 25, 22},
                BorderColor:     "rgb(75, 192, 192)",
                BackgroundColor: "rgba(75, 192, 192, 0.2)",
                Fill:            true,
            },
        },
    )

    chart.Options = &chartjs.Options{}
    chartjs.ApplyOptions(chart.Options,
        chartjs.WithTitle("Monthly Sales"),
        chartjs.WithXAxisLabel("Month"),
        chartjs.WithYAxisLabel("Sales ($)"),
        chartjs.WithYAxisBeginAtZero(true),
    )

    cjs := chartjs.New(chart)
    cjs.MustBuild(chartjs.OutputAll, os.Stdout)
}
```

### Advanced Configuration

```go
chart := &chartjs.Chart{
    Type: chartjs.ChartTypeLine,
    Data: &chartjs.Data{
        Labels: []any{"Jan", "Feb", "Mar"},
        Datasets: []chartjs.Dataset{
            {
                Label:           "Dataset 1",
                Data:            []any{10.0, 19.0, 15.0},
                BorderColor:     "rgb(75, 192, 192)",
                BackgroundColor: "rgba(75, 192, 192, 0.2)",
                Fill:            boolPtr(true),
                Tension:         float64Ptr(0.4),
            },
        },
    },
    Options: &chartjs.Options{
        Responsive: boolPtr(true),
        Plugins: &chartjs.Plugins{
            Title: &chartjs.Title{
                Display:  boolPtr(true),
                Text:     stringPtr("My Chart"),
                Position: stringPtr("top"),
            },
            Legend: &chartjs.Legend{
                Display:  boolPtr(true),
                Position: stringPtr("bottom"),
            },
        },
        Scales: &chartjs.Scales{
            Y: &chartjs.Scale{
                BeginAtZero: boolPtr(true),
                Title: &chartjs.ScaleTitle{
                    Display: boolPtr(true),
                    Text:    stringPtr("Value"),
                },
            },
        },
    },
}

cjs := chartjs.New(chart)
cjs.MustBuild(chartjs.OutputAll, os.Stdout)
```

## Helper Functions

The package provides convenient helper functions for common configurations:

- `WithTitle(text string)` - Add a chart title
- `WithSubtitle(text string)` - Add a chart subtitle
- `WithLegend(display bool, position string)` - Configure legend
- `WithXAxisLabel(label string)` - Set X-axis label
- `WithYAxisLabel(label string)` - Set Y-axis label
- `WithYAxisBeginAtZero(beginAtZero bool)` - Set Y-axis to begin at zero
- `WithLogarithmicScale()` - Use logarithmic Y-axis
- `WithTimeScale(unit string)` - Configure time-based X-axis
- `WithResponsive(responsive bool)` - Set responsive mode
- `WithAspectRatio(ratio float64)` - Set aspect ratio
- `WithAnimation(duration int, easing string)` - Configure animation
- `WithStacked(stacked bool)` - Enable stacked mode

## Output Modes

- `OutputAll` - Complete HTML file with chart
- `OutputHTMLHeader` - HTML header with dependencies
- `OutputDependencies` - Only Chart.js script tags
- `OutputChart` - Only the chart JavaScript configuration
- `OutputHTMLFooter` - HTML footer

## Data Structures

The package supports various data formats:

- **Simple arrays**: `[]float64` for basic charts
- **Data points**: `DataPoint{X, Y, R}` for scatter/bubble charts
- **Time series**: `TimeDataPoint{X: time.Time, Y: float64}` for time-based charts
- **Mixed types**: `[]any` for flexible data input

## Chart.js Features

All Chart.js features are supported:

- **Scales**: Linear, logarithmic, time, category, radial linear
- **Plugins**: Legend, title, tooltip, subtitle, zoom, decimation, filler
- **Elements**: Points, lines, bars, arcs with full customization
- **Animations**: Complete animation configuration
- **Interactions**: Hover, click, and interaction modes
- **Colors**: Support for color strings, arrays, and functions

## Documentation

For detailed Chart.js documentation, see:
- [Chart.js Getting Started](https://www.chartjs.org/docs/latest/getting-started/usage.html)
- [Chart.js Options](https://www.chartjs.org/docs/latest/general/options.html)
- [Chart.js Data Structures](https://www.chartjs.org/docs/latest/general/data-structures.html)
- [Chart.js Colors](https://www.chartjs.org/docs/latest/general/colors.html)

## Examples

See `example_test.go` for comprehensive examples of:
- Basic line charts
- Bar charts with multiple datasets
- Time series charts
- Scatter plots
- Pie charts
- Advanced configurations

