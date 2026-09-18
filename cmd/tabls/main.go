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
		"%d-%d-%d %d:%d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
	)
}

func main() {
	entries, err := os.ReadDir(".")
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

		entryType := "File"
		if info.IsDir() {
			entryType = "Directory"
		}

		tab.AddRow(info.Name(), entryType, info.Size(), timeToString(info.ModTime()), info.Mode())
	}

	tab.Print()
}
