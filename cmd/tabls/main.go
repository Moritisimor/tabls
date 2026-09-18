package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func timeToString(t time.Time) string {
	return fmt.Sprintf(
		"%d-%d-%d %d:%d:%d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	tab := table.New("Name", "Type", "Size", "Modified at", "Mode")
	headerFmt := color.New(color.FgHiGreen, color.Underline).SprintfFunc()
  	columnFmt := color.New(color.FgHiMagenta).SprintfFunc()
	tab.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		entryType := "Symlink"
		if entry.Type().IsRegular() {
			entryType = "File"
		}

		if entry.Type().IsDir() {
			entryType = "Directory"
		}

		tab.AddRow(
			info.Name(),
			entryType,
			info.Size(), 
			timeToString(info.ModTime()),
			info.Mode(),
		)
	}

	tab.Print()
}
