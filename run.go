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

const cmdName = "mtglight"

// Run the yeelight
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")

	// Command line flags passed from overSight
	// ref. https://objective-see.org/products/oversight.html
	var (
		device, event    string
		pid, activeCount int
	)
	fs.StringVar(&device, "device", "", "device: camera/microphone")
	fs.StringVar(&event, "event", "", "event: on/off")
	// NOTE: This process ID is not very useful. It is not an online meeting process but a device
	// management process above it. So, it is impossible to determine whether a meeting is in session
	// by checking for the existence of this process.
	fs.IntVar(&pid, "process", -1, "process ID: (when off, the process number is empty)")
	fs.IntVar(&activeCount, "activeCount", -1, "active count: (total count of cameras and microphones combined)")

	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}

	if event == "" {
		return fmt.Errorf("no events specified")
	}
	on := event == "on"

	// Manipulate the light when the camera is turned on or when the activeCount reaches 0.
	// The reason for not operating the camera when it is off is that even if one camera is turned off,
	// it may still be being used on other screens, so we use activeCount to determine if the camera is
	// off. Also, not that activeCount is the sum of both cameras and microphones, and the device flag
	// might be specified as a microphone when it reaches 0. Apropos, when activeCount is 0, the event
	// is always off.
	// Either way, when activeCount reaches 0, it shows the online meeting is inevitably over,
	// so we'll turn off the light.
	// Fortunately, it prevents the lights from going out even when we temporarily turn off the
	// camera during a meeting.
	if !((on && device == "camera") || activeCount == 0) {
		return nil
	}

	var yee *Yeelight
	r := &retryer{}
	r.run(func() error {
		var err error
		yee, err = Discover()
		return err
	})
	if yee != nil {
		r.run(func() error { return yee.Power(on) })
		if on {
			r.run(func() error { return yee.RGB(0xffff00) })
			r.run(func() error { return yee.Brightness(99) })
		}
	}
	return r.err
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
