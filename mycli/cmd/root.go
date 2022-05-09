package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "MyDB Server simple cli",
	Long:  `MyDB Server cli tools`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.bucli.yaml)")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.PersistentFlags().StringP("filepath", "", "./", "file path to data write")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	viper.BindPFlag("filePath", rootCmd.PersistentFlags().Lookup("filepath"))
	viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Config file not found!")
		os.Exit(1)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("config file not found")
		} else {
			// Config file was found but another error was produced
			fmt.Println("config file found but something wrong")
		}

		fmt.Println("config file reader error:", err)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func CreateDir(path string) error {
	// Check if folder exists
	_, err := os.Stat(path)

	// Create directory if not exists
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModeDir|0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateFile(filename string) error {
	// Check if folder exists
	_, err := os.Stat(filename)

	// Create directory if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFileOrDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(filepath string) string {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ""
	}

	return string(data)
}

func WriteFileLineByLine(filepath string, data []string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, line := range data {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteFileAll(filepath string, data string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	file.Sync()

	return nil
}
