package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DbList []struct {
	DbStore
}

type DbStore struct {
	SyncInterval  int      `json:"sync-interval"`
	RedisAddr     string   `json:"redis-addr"`
	RedisPort     int      `json:"redis-port"`
	DbStorage     string   `json:"db-storage"`
	DbName        string   `json:"db-name"`
	DbDisplayName string   `json:"db-display-name"`
	DbPath        string   `json:"db-path"`
	DbReadOnly    bool     `json:"db-read-only"`
	ContainerName string   `json:"container-name"`
	BucketList    []string `json:"bucket-list"`
	IndexList     []string `json:"index-list"`
}

// generatorCmd represents the generator command
var generatorCmd = &cobra.Command{
	Use:   "generator",
	Short: "generator everything",
	Long:  `generator file - folder anything you want`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("generator called")

		// Get Database List
		dbNames := viper.GetStringSlice("databases")

		// Get Service Name for which folder move into the generated files
		serviceName := viper.GetString("service-name")

		if len(dbNames) > 0 {

			// Database List
			dbList := DbList{}

			// Get All Tmpl Files
			tmplFiles, _ := template.ParseGlob("./templates/*")

			for _, dbKey := range dbNames {
				// Store Parameters
				dbParams := viper.GetStringMap(dbKey)
				fmt.Println(dbParams)

				// Sorry about that
				dbStore := &DbStore{}
				dbParamsBytes, _ := json.Marshal(dbParams)
				json.Unmarshal(dbParamsBytes, dbStore)

				// Generate: env file ///////
				generateTemplate2File(tmplFiles, "env.tmpl", fmt.Sprintf("../%s/", serviceName), fmt.Sprintf(".%s.env", dbStore.DbName), dbStore)

				// Generate: service file ///////
				generateTemplate2File(tmplFiles, "service.tmpl", fmt.Sprintf("../%s/", serviceName), fmt.Sprintf("%s.service", dbStore.DbName), dbStore)

				// Generate: store file ///////
				generateTemplate2File(tmplFiles, "store.go.tmpl", fmt.Sprintf("../%s/store/%s/", serviceName, dbStore.DbName), "store.go", dbStore)

				// Collect DbStore data
				dbList = append(dbList, struct{ DbStore }{*dbStore})
			}

			// Generate: storage file ///////
			generateTemplate2File(tmplFiles, "storage.go.tmpl", fmt.Sprintf("../%s/services/", serviceName), "storage.go", dbList)

			// Generate: docker-compose file ///////
			generateTemplate2File(tmplFiles, "docker-compose.yml.tmpl", "../", "docker-compose.yml", dbList)

		}

		fmt.Println("generator process finished")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generatorCmd)
}

func generateTemplate2File(tmplFiles *template.Template, tmplName string, dest, filename string, data interface{}) {
	// Check Directory if need create
	_ = CreateDir(dest)

	// Create a new file
	newFile, err := os.Create(dest + filename)
	checkError(err)

	// Save new file
	err = tmplFiles.ExecuteTemplate(newFile, tmplName, data)
	checkError(err)

	fmt.Println("New file generated", newFile.Name())
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
