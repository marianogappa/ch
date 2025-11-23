package debug

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/output"
	"github.com/marianogappa/ch/pkg/output/frontends/debug"
)

func init() {
	ch.RegisterOutput(NewDebugDriver())
}

type DebugDriver struct{}

func NewDebugDriver() *DebugDriver {
	return &DebugDriver{}
}

func (d *DebugDriver) Name() string {
	return "debug"
}

// Ensure DebugDriver implements ChartOutput interface
var _ output.ChartOutput = (*DebugDriver)(nil)

type DebugConfig struct {
	Pretty     bool
	OutputPath string
}

func (d *DebugDriver) RegisterFlags(fs *flag.FlagSet) any {
	c := &DebugConfig{}
	fs.BoolVar(&c.Pretty, "pretty", false, "Pretty print JSON output.")
	fs.StringVar(&c.OutputPath, "output-path", "", "Output path: empty for stdout, '-' for stdout, or file path.")
	return c
}

func (d *DebugDriver) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   true,
		Interactive: false,
	}
}

// GetFrontendSettings returns frontend-specific setting definitions
// Implements output.ChartOutput interface
func (d *DebugDriver) GetFrontendSettings() map[string]output.SettingDefinition {
	return nil // Debug driver doesn't support frontend-specific settings
}

// RenderChart renders the chart using the generic ChartConfig
// Implements output.ChartOutput interface
// Debug output outputs chart config as first line, then raw data
func (d *DebugDriver) RenderChart(rows <-chan ch.Row, config *output.ChartConfig) error {
	// Use default pretty=false for RenderChart
	frontend := debug.New(false)

	// Determine output destination
	var outputWriter *os.File
	if config.OutputPath == "" || config.OutputPath == "-" {
		outputWriter = os.Stdout
	} else {
		f, err := os.Create(config.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		outputWriter = f
	}

	// Output chart config as first line
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling chart config: %w", err)
	}
	if _, err := fmt.Fprintf(outputWriter, "%s\n", configJSON); err != nil {
		return fmt.Errorf("error writing chart config: %w", err)
	}

	// Stream rows to frontend
	for row := range rows {
		if err := frontend.Encode(row, outputWriter); err != nil {
			return err
		}
	}

	if config.OutputPath != "" && config.OutputPath != "-" {
		fmt.Printf("Debug output written to %s\n", config.OutputPath)
	}

	return nil
}

func (d *DebugDriver) Render(rows <-chan ch.Row, config any) error {
	cfg, ok := config.(*DebugConfig)
	if !ok {
		return fmt.Errorf("invalid config type for DebugDriver")
	}

	// Create frontend instance
	frontend := debug.New(cfg.Pretty)

	// Determine output destination
	var outputWriter *os.File
	if cfg.OutputPath == "" || cfg.OutputPath == "-" {
		outputWriter = os.Stdout
	} else {
		f, err := os.Create(cfg.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		outputWriter = f
	}

	// Stream rows to frontend
	for row := range rows {
		if err := frontend.Encode(row, outputWriter); err != nil {
			return err
		}
	}

	if cfg.OutputPath != "" && cfg.OutputPath != "-" {
		fmt.Printf("Debug output written to %s\n", cfg.OutputPath)
	}

	return nil
}
