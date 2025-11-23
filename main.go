package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/input"
	"github.com/marianogappa/ch/pkg/llm"
	"github.com/marianogappa/ch/pkg/output"
	_ "github.com/marianogappa/ch/pkg/output/drivers/chartjs"
	_ "github.com/marianogappa/ch/pkg/output/drivers/d3"
	_ "github.com/marianogappa/ch/pkg/output/drivers/debug"
	"github.com/marianogappa/ch/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	outputName  string
	separator   string
	dateFormat  string
	lineFormat  string
	interactive bool
	apiKey      string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ch",
	Short: "Chart data from stdin",
	Long:  "Chart data from stdin with various output formats",
	RunE:  run,
}

var chartConfig *output.ChartConfig
var parseChartConfigFunc func(outputName string) error

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputName, "output", "o", "chartjs", "Output driver (chartjs, d3, debug)")
	rootCmd.PersistentFlags().StringVar(&separator, "separator", "\t", "Column separator")
	rootCmd.PersistentFlags().StringVar(&dateFormat, "date-format", "", "Date format")
	rootCmd.PersistentFlags().StringVar(&lineFormat, "format", "", "Line format (e.g. 'sfd')")
	rootCmd.PersistentFlags().BoolVar(&interactive, "interactive", false, "Interactive mode (LLM)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "LLM API Key")

	// Register all chart config flags (generic + all frontend-specific with prefixes)
	chartConfig, parseChartConfigFunc = output.RegisterAllFlags(rootCmd.PersistentFlags())
}

func run(cmd *cobra.Command, args []string) error {
	// Get output driver
	outDriver, err := ch.GetOutput(outputName)
	if err != nil {
		return fmt.Errorf("error: %v. Available outputs: %v", err, ch.Outputs())
	}

	// All current drivers implement ChartOutput
	var chartOutput output.ChartOutput

	if output.IsChartOutput(outDriver) {
		chartOutput = outDriver.(output.ChartOutput)
	} else {
		// Fallback for drivers that don't implement ChartOutput
		// Convert pflag to flag.FlagSet for compatibility
		fs := flag.NewFlagSet("ch", flag.ContinueOnError)
		// Note: This won't work perfectly, but it's a fallback
		// Drivers should implement ChartOutput
		_ = outDriver.RegisterFlags(fs)
		return fmt.Errorf("driver %s does not implement ChartOutput", outputName)
	}

	// Parse flags
	if err := cmd.ParseFlags(args); err != nil {
		return err
	}

	// Call parse function for special cases (like Colors) and populate frontend settings
	if parseChartConfigFunc != nil {
		if err := parseChartConfigFunc(outputName); err != nil {
			return err
		}
	}

	// Setup Input
	in := input.NewReaderInput(os.Stdin)

	// Setup Parser
	sepRune := []rune(separator)[0] // simplistic
	if separator == "\\t" {
		sepRune = '\t'
	} // handle escaped tab from shell

	p := parser.NewCSVParser(sepRune, dateFormat)
	if lineFormat != "" {
		p.LineFormat = lineFormat
	}

	// Interactive Mode
	if interactive {
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			return fmt.Errorf("interactive mode requires OPENAI_API_KEY")
		}
		client := llm.NewOpenAIClient(apiKey)
		controller := llm.NewController(client)
		_ = controller
		fmt.Println("Interactive mode enabled (Mocked)")
		// TODO: Implement full interactive flow
	}

	// Run
	stream, err := in.Stream()
	if err != nil {
		return fmt.Errorf("error creating input stream: %v", err)
	}

	rows, err := p.Parse(stream)
	if err != nil {
		return fmt.Errorf("error creating parser: %v", err)
	}

	// Render using ChartConfig
	if chartConfig != nil && chartOutput != nil {
		if err := chartOutput.RenderChart(rows, chartConfig); err != nil {
			return fmt.Errorf("error rendering output: %v", err)
		}
	} else {
		return fmt.Errorf("chart config not available")
	}

	return nil
}
