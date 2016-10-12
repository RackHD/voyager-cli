// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-cleanhttp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var fileLocation string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create configuration manifest",
	Long:  `Create manifest containing the deisred system configuration. Usage: mcc manifest create -f [filename.json]`,

	Run: func(cmd *cobra.Command, args []string) {

		url, err := CheckMccTarget()
		if err != nil {
			log.Fatal(err)
		}

		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		fileLocation := fmt.Sprintf("%s/voyager-manifest.json", dir)

		// Read IP target from file
		_, err = os.Stat(fileLocation)
		if err != nil {
			color.Red("Cannot find voyager manifest file (expecting $HOME/voyager-manifest.json)\n")
			return
		}

		// Read in manifest file to json data object
		manifestFile, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			color.Red(fmt.Sprintf("Cannot find manifest file: %s\n", fileLocation))
			return

		}
		manifestBody := bytes.NewReader(manifestFile)

		// Send http GET request to IP target
		client := cleanhttp.DefaultClient()
		targetURL := fmt.Sprintf("%s://%s/manifest", url.Scheme, url.Host)
		req, err := http.NewRequest("POST", targetURL, manifestBody)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(req)

		if err != nil {
			color.Red("Error sending '/manifest' API call to Voyager at %s: %s\n", url.Host, err)
			return
		}
		defer resp.Body.Close()

		// Parse response
		// responseBytes, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	color.Red("Error reading response body: %s\n", err)
		// 	return
		// }

		if resp.StatusCode != 200 {
			color.Red("Server returned (%d)\n", resp.StatusCode)
			return
		}
		color.Green("Manifest file successfully uploaded\n")
	},
}

func init() {
	manifestCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&fileLocation, "file", "f", "$HOME/voyager-manifest.json", "file location eg. $HOME/voyager-manifest.json")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
