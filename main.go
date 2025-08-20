package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/pkg/errors"
)

var version string

func main() {
	if len(os.Args) > 0 && os.Args[0] == "version" {
		fmt.Println(version)
		os.Exit(1)
	}

	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() error {
	parser := hclparse.NewParser()

	matches, err := filepath.Glob("*.tf")
	if err != nil {
		return err
	}

	for i := range matches {
		ws, err := parseHCLFile(parser, matches[i])
		if err != nil {
			return err
		}
		if ws != nil {
			if err := openBrowser(ws); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	}

	return errors.New("no Terraform configuation found")
}

func openBrowser(ws *Workspace) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", ws.toURL()).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", ws.toURL()).Start()
	case "darwin":
		err = exec.Command("open", ws.toURL()).Start()
	default:
		err = fmt.Errorf("unsupport runtime OS: %s", runtime.GOOS)
	}

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
