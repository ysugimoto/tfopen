package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/hcl/v2/hclparse"
)

func main() {
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
			openBrowser(ws)
			break
		}
	}

	return nil
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
		err = fmt.Errorf("Unsupport runtime OS: %s", runtime.GOOS)
	}

	return err
}
