package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caproven/musicstats/internal"
	"github.com/spf13/cobra"
)

var musicDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "musicstats",
	Short: "gather stats about your local music collection",
	Long: `musicstats: Gather stats about your local music collection,
such as total duration or breakdown of file type or genre.

All common audio file extensions are supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		if musicDir == "" {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			musicDir = filepath.Join(home, "Music")
		}

		musicFilePaths, err := internal.GetAllMusicFiles(musicDir, []string{})
		if err != nil {
			cobra.CheckErr(err)
		}

		totalDuration, err := internal.TotalDuration(musicFilePaths)
		if err != nil {
			cobra.CheckErr(err)
		}

		PrintResults(len(musicFilePaths), totalDuration)
	},
}

func PrintResults(numFiles int, duration time.Duration) {
	if numFiles == 0 {
		fmt.Println("No music files found")
	} else {
		fmt.Println("Found a total of", numFiles, "music file(s)")
		fmt.Println("Total duration is", duration)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&musicDir, "folder", "f", "", "directory containing audio/music files (default is $HOME/Music")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
