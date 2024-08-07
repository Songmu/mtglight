package mtglight

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) (err error) {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")
	var (
		device, event, stateF string
		pid, activeCount      int
	)
	// for OverSight
	// ref. https://objective-see.org/products/oversight.html
	fs.StringVar(&device, "device", "", "device camera/microphone")
	fs.StringVar(&event, "event", "", "event on/off")
	fs.IntVar(&pid, "process", -1, "process ID (note: when off, the process number is empty)")
	fs.IntVar(&activeCount, "activeCount", 0, "active count (total count of cameras and microphones combined)")
	fs.StringVar(&stateF, "state", "", "state file")

	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}
	if device != "camera" {
		return nil
	}
	if event == "" {
		return fmt.Errorf("no events specified")
	}

	var sf *stateFile
	if stateF != "" {
		sf = &stateFile{path: stateF}
		if err := sf.load(); err != nil {
			return err
		}
		if pid > 0 {
			sf.addProcess(pid)
		}
		defer func() {
			if e := sf.save(); e != nil {
				if err == nil {
					err = e
				} else {
					log.Printf("failed to save state file: %s", e)
				}
			}
		}()
	}

	on := event == "on"
	var yee *Yeelight
	if err := retry.Retry(3, time.Second, func() error {
		var err error
		yee, err = Discover()
		return err
	}); err != nil {
		return err
	}

	r := &retryer{}
	if on || activeCount == 0 || (sf != nil && len(sf.state.Processes) == 0) {
		r.run(func() error { return yee.Power(on) })
	}
	if on {
		r.run(func() error { return yee.RGB(0xffff00) })
		r.run(func() error { return yee.Brightness(99) })
	}
	return r.err
}

type stateFile struct {
	path  string
	state state
}

func (sf *stateFile) load() error {
	f, err := os.Open(sf.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&sf.state)
}

func (sf *stateFile) save() error {
	f, err := os.Create(sf.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(sf.state)
}

func (sf *stateFile) updateProcesses() {
	var processes []int
	for _, pid := range sf.state.Processes {
		if isProcessActive(pid) {
			processes = append(processes, pid)
		}
	}
	sf.state.Processes = processes
}

func (sf *stateFile) addProcess(pid int) {
	if pid != -1 {
		sf.state.Processes = append(sf.state.Processes, pid)
	}
	sf.updateProcesses()
}

type state struct {
	Processes []int `json:"processes"`
}

func isProcessActive(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
