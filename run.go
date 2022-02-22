package mtglight

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Songmu/retry"
)

type retryer struct {
	err error
}

func (rt *retryer) run(f func() error) {
	if rt.err != nil {
		return
	}
	rt.err = retry.Retry(3, time.Second, f)
}

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

	var yee *Yeelight
	if err := retry.Retry(3, time.Second, func() error {
		var err error
		yee, err = Discover()
		return err
	}); err != nil {
		return err
	}

	r := &retryer{}
	r.run(func() error { return yee.Power(on) })
	if on {
		r.run(func() error { return yee.RGB(0x993333) })
		r.run(func() error { return yee.Brightness(1) })
	}
	return r.err
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
