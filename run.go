package yeelight

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
)

const cmdName = "yeelight"

// Run the yeelight
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")
	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("no option specified")
	}
	on := fs.Arg(0) == "on"

	yee, err := Discover()
	if err != nil {
		return err
	}
	yee.SetPower(on)

	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
