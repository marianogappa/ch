package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/input"
	"github.com/marianogappa/ch/pkg/llm"
	_ "github.com/marianogappa/ch/pkg/output/chartjs"
	_ "github.com/marianogappa/ch/pkg/output/d3"
	_ "github.com/marianogappa/ch/pkg/output/json"
	"github.com/marianogappa/ch/pkg/parser"
)

func main() {
	if err := Run(os.Args, os.Stdin); err != nil {
		log.Fatal(err)
	}
}

func Run(args []string, stdin io.Reader) error {
	// 1. Parse global flags to determine output driver and common options
	outputName := "chartjs" // default
	for i, arg := range args {
		if (arg == "--output" || arg == "-o") && i+1 < len(args) {
			outputName = args[i+1]
			break
		}
	}

	outDriver, err := ch.GetOutput(outputName)
	if err != nil {
		return fmt.Errorf("error: %v. Available outputs: %v", err, ch.Outputs())
	}

	// Setup FlagSet
	fs := flag.NewFlagSet("ch", flag.ContinueOnError)

	// Global flags
	var (
		separator     string
		dateFormat    string
		rawLineFormat string
		interactive   bool
		apiKey        string
	)

	fs.StringVar(&separator, "separator", "\t", "Column separator")
	fs.StringVar(&dateFormat, "date-format", "", "Date format")
	fs.StringVar(&rawLineFormat, "format", "", "Line format (e.g. 'sfd')")
	fs.BoolVar(&interactive, "interactive", false, "Interactive mode (LLM)")
	fs.StringVar(&apiKey, "api-key", "", "LLM API Key")

	// Register output flags
	outConfig := outDriver.RegisterFlags(fs)

	// Parse
	var dummyOutput string
	fs.StringVar(&dummyOutput, "output", "chartjs", "Output driver")
	fs.StringVar(&dummyOutput, "o", "chartjs", "Output driver")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	// 2. Setup Input
	in := input.NewReaderInput(stdin)

	// 3. Setup Parser
	sepRune := []rune(separator)[0] // simplistic
	if separator == "\\t" {
		sepRune = '\t'
	} // handle escaped tab from shell

	p := parser.NewCSVParser(sepRune, dateFormat)
	if rawLineFormat != "" {
		p.LineFormat = rawLineFormat
	}

	// 4. Interactive Mode
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

	// 5. Run
	stream, err := in.Stream()
	if err != nil {
		return fmt.Errorf("error creating input stream: %v", err)
	}

	rows, err := p.Parse(stream)
	if err != nil {
		return fmt.Errorf("error creating parser: %v", err)
	}

	if err := outDriver.Render(rows, outConfig); err != nil {
		return fmt.Errorf("error rendering output: %v", err)
	}

	return nil
}
