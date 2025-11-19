package llm

import (
	"context"
	"fmt"
	"strings"
)

type Controller struct {
	Client Client
}

func NewController(client Client) *Controller {
	return &Controller{Client: client}
}

// SuggestConfiguration takes a sample of data and asks the LLM for configuration.
// Returns (lineFormat, chartType, error)
func (c *Controller) SuggestConfiguration(ctx context.Context, dataSample []string) (string, string, error) {
	prompt := fmt.Sprintf(`
I have a dataset. Here are the first few lines:
%s

Please analyze this data and suggest:
1. A format string for parsing (using 's' for string, 'f' for float, 'd' for date).
2. A suitable chart type (e.g., 'line', 'bar', 'scatter', 'pie').

Reply in the format:
FORMAT: <format_string>
CHART: <chart_type>
`, strings.Join(dataSample, "\n"))

	resp, err := c.Client.Complete(ctx, prompt)
	if err != nil {
		return "", "", err
	}

	// Parse response (naive implementation)
	var format, chart string
	lines := strings.Split(resp, "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "FORMAT: ") {
			format = strings.TrimPrefix(l, "FORMAT: ")
		}
		if strings.HasPrefix(l, "CHART: ") {
			chart = strings.TrimPrefix(l, "CHART: ")
		}
	}

	return strings.TrimSpace(format), strings.TrimSpace(chart), nil
}
