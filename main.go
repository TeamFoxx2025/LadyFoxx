package main

import (
	_ "embed"

	"github.com/TeamFoxx2025/LadyFoxx/command/root"
	"github.com/TeamFoxx2025/LadyFoxx/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
