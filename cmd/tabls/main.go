package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func timeToString(t time.Time) string {
	return fmt.Sprintf(
		"%d-%02d-%02d %02d:%02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

func sizeToString(sizeBytes int64) string {
	if sizeBytes > 1_000_000_000 {
		return fmt.Sprintf("%.2f GB", float64(sizeBytes) / 1_000_000_000)
	}

	if sizeBytes > 1_000_000 {
		return fmt.Sprintf("%.2f MB", float64(sizeBytes) / 1_000_000)
	}

	if sizeBytes > 1000 {
		return fmt.Sprintf("%.2f KB", float64(sizeBytes) / 1000)
	}

	return fmt.Sprintf("%d B", sizeBytes)
}

func calculateDirSize(dirPath string) (int64, error) {
	var acc int64 = 0
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			acc += info.Size()
		}

		return nil
	})

	return acc, err
}

func main() {
	path := "."
	showHidden := false
	calcDirSize := false

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help":
				color.Blue("Tabls on https://github.com/Moritisimor/tabls")
				color.Green("Usage: tabls <flags...> <directory>")
				color.Green("Flags:")
				color.Magenta("\t-h | --help   \tPrints this")
				color.Magenta("\t-r | --recurse\tRecursively calculates directory sizes")
				color.Magenta("\t-a | --all    \tShows hidden directories (Those starting with '.')")

				return

			case "-r", "--recurse":
				calcDirSize = true

			case "-a", "--all":
				showHidden = true
			
			default:
				color.Red("Unknown flag: %s", arg)
			}

			continue
		}

		path = arg
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
		if strings.HasPrefix(entry.Name(), ".") {
			if !showHidden {
				continue
			}
		}

		info, err := entry.Info()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		entryType := "Symlink"
		if entry.Type().IsRegular() {
			entryType = "File"
		}

		entrySize := info.Size()
		if entry.Type().IsDir() {
			entryType = "Directory"

			if calcDirSize {
				absPath := filepath.Join(path, info.Name())
				entrySize, err = calculateDirSize(absPath)
				if err != nil {
					color.Red(err.Error())
					os.Exit(1)
				}
			}
		}

		tab.AddRow(
			info.Name(),
			entryType,
			sizeToString(entrySize), 
			timeToString(info.ModTime()),
			info.Mode(),
		)
	}

	tab.Print()
}
