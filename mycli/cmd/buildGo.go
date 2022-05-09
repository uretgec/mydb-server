package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildGoCmd represents the buildGo command
var buildGoCmd = &cobra.Command{
	Use:   "go",
	Short: "build all services binary file with necessary folders",
	Long:  `build all services binary file with necessary folders`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("buildGo called")

		// OS Env
		osType := viper.GetString("ostypes") // runtime.GOOS
		if osType == "none" {
			osType = runtime.GOOS
		}

		archType := viper.GetString("archtypes") //runtime.GOARCH
		if archType == "none" {
			archType = runtime.GOARCH
		}

		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Printf("go build error: %v\n", err)
			panic(err)
		}

		// Services
		services := viper.GetStringSlice("services")

		if len(services) > 0 {
			for _, service := range services {

				// Service Parameters
				serviceParams := viper.GetStringMapString(service)

				// Create Dest Folders
				err := CreateDir(serviceParams["dest"])
				if err != nil {
					fmt.Printf("dest folder create error: %v\n", err)
					os.Exit(1)
				}

				// Go Build
				gcmd := exec.Command("/bin/sh", "./sh/gobuild.sh", service, osType, archType, serviceParams["path"], serviceParams["dest"], serviceParams["output"], serviceParams["copyto"], serviceParams["senv"])
				gcmd.Env = append(os.Environ(), "GOOS="+osType, "GOARCH="+archType)
				gcmd.Dir = currentPath

				fmt.Printf("Run Command: %s\n", gcmd.String())

				out, err := gcmd.CombinedOutput()
				if err != nil {
					fmt.Printf("go build out: %v\n", out)
					fmt.Printf("go build error: %v\n", err)
					panic(err)
				}

				fmt.Printf("%s service build command output: %s\n", service, out)

				time.Sleep(1 * time.Second)
			}
		}

		fmt.Println("service build process completed")

		return nil
	},
}

func init() {
	buildCmd.AddCommand(buildGoCmd)
}
