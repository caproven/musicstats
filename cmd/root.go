package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caproven/musicstats/internal"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "musicstats",
	Short: "gather stats about your local music collection",
	Long: `musicstats: Gather stats about your local music collection,
such as total duration or breakdown of file type or genre.

All common audio file extensions are supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		folder := viper.GetString("folder")

		musicFilePaths, err := internal.GetAllMusicFiles(folder, []string{})
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.musicstats.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("folder", "f", "", "Folder to scan (default is $HOME/Music)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".musicstats" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".musicstats")
	}

	viper.SetDefault("folder", filepath.Join(home, "Music"))
	viper.BindPFlags(rootCmd.Flags())

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
