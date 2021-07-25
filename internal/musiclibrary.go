package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

// TODO switch to channels for processing instead of returning slice.
// exiftool library takes a slice of paths (strings), so see if it can
// be adapted to use a channel
func GetAllMusicFiles(dir string, files []string) ([]string, error) {
	fmt.Println("Scanning", dir, "for music files..")
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// TODO how to differentiate between ReadDir failing and fs.Stat failing?
			return filepath.SkipDir
		}

		if !d.IsDir() {
			for _, ext := range AUDIO_FILE_EXTENSIONS {
				if filepath.Ext(path) == ext {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	return files, err
}

func TotalDuration(files []string) (time.Duration, error) {
	fmt.Println("Calculating total duration..")
	var total time.Duration

	exif, err := exiftool.NewExiftool()
	if err != nil {
		return total, err
	}
	defer exif.Close()

	for _, metadata := range exif.ExtractMetadata(files...) {
		if val, ok := metadata.Fields["Duration"]; ok {
			duration := exifDuration(val)
			total += duration
		}
	}
	return total, nil
}

func exifDuration(data interface{}) time.Duration {
	// Parses a time.Duration from the given exif data's
	// "Duration" field. If the data cannot be parsed or
	// interpreted, the value of 0 is returned.
	switch data := data.(type) {
	case string:
		// if duration < 1m, syntax is "#.## s[ (approx)]"
		// if duration >= 1m, syntax is "##.##.##[ (approx)]"
		data = strings.TrimSuffix(data, "(approx)")
		data = strings.ReplaceAll(data, " ", "")

		var hours, minutes, seconds int
		n, err := fmt.Sscanf(data, "%d:%d:%d", &hours, &minutes, &seconds)
		if err != nil || n != 3 {
			duration, err := time.ParseDuration(data)
			if err != nil {
				return 0
			}
			return duration
		} else {
			return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
		}
	}

	// return 0 duration if we don't know how to parse it
	return 0
}
