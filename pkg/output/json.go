package output

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
)

func init() {
	ch.RegisterOutput(NewJSONOutput())
}

type JSONOutput struct{}

func NewJSONOutput() *JSONOutput {
	return &JSONOutput{}
}

func (o *JSONOutput) Name() string {
	return "json"
}

type JSONConfig struct {
	Pretty bool
}

func (o *JSONOutput) RegisterFlags(fs *flag.FlagSet) any {
	c := &JSONConfig{}
	fs.BoolVar(&c.Pretty, "pretty", false, "Pretty print JSON output.")
	return c
}

func (o *JSONOutput) Capabilities() ch.Capabilities {
	return ch.Capabilities{
		Streaming:   true,
		Interactive: false,
	}
}

func (o *JSONOutput) Render(rows <-chan ch.Row, config any) error {
	cfg, _ := config.(*JSONConfig)

	enc := json.NewEncoder(os.Stdout)
	if cfg != nil && cfg.Pretty {
		enc.SetIndent("", "  ")
	}

	for row := range rows {
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}
