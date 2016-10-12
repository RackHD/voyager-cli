// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/RackHD/voyager-houston/model"
	"github.com/mitchellh/go-homedir"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "IP endpoint of a Voyager system to target",
	Long: `IP endpoint of a Voyager system to target.
This command will attempt a call to the IP specified, and if successful,
will store the IP for future use in the .voyager file

Usage: mcc target http://<ip address>:<port>
ex.: mcc target http://192.168.1.1:80`,

	Run: func(c *cobra.Command, args []string) {

		// Argument required
		if len(args) > 1 {
			fmt.Printf("Too Many Arguments.\n")
			return
		} else if len(args) < 1 {
			// Check and see if the endpoint has been set
			dir, err := homedir.Dir()

			if err != nil {
				log.Fatal(err)
			}
			fileLocation := fmt.Sprintf("%s/.voyager", dir)
			if _, err := os.Stat(fileLocation); err == nil {
				fileContent, err := ioutil.ReadFile(fileLocation)
				if err != nil {
					log.Fatal(err)
				}
				// Unmarshal data and print
				urlObject := url.URL{}
				err = json.Unmarshal(fileContent, &urlObject)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Current target is ")
				color.Green("%s://%s\n", urlObject.Scheme, urlObject.Host)
			} else {
				color.Red("No target set.\n")
			}
		} else {
			// Parse and validate argument
			target, err := url.Parse(args[0])
			if err != nil {
				color.Red("Could not convert arg to IP Address: %s\n", err)
				return
			}

			if target.Host == "" || target.Scheme == "" {
				fmt.Printf("Please enter a valid target url. ex: http://192.168.1.1:80\n")
				return
			}

			// Convert argument to REST call
			targetURL := fmt.Sprintf("%s://%s/info", target.Scheme, target.Host)

			// Send API call to validate that argument points to running Voyager instance
			res, err := http.Get(targetURL)

			// Error check the REST call
			if err != nil {
				fmt.Printf("Error sending '/info' API call to Voyager: %s\n", err)
				return
			}
			if res.StatusCode != 200 {
				fmt.Printf("Non-success status code returned:")
				color.Red("%d\n", res.StatusCode)
				return
			}

			respBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %s\n", err)
			}

			var houston model.Houston
			err = json.Unmarshal(respBytes, &houston)
			if err != nil {
				fmt.Printf("Invalid response body: %s\nexited with error: %s", respBytes, err)
				return
			}

			color.Green("Target set to %s://%s\n", target.Scheme, target.Host)

			// Store target URL of valid Voyager endpoint
			targetb, err := json.Marshal(target)
			if err != nil {
				fmt.Printf("Could not marshal IP Address to JSON: %s\n", err)
				return
			}

			// Determine where to store state file
			dir, err := homedir.Dir()
			if err != nil {
				log.Fatal(err)
			}
			fileLocation := fmt.Sprintf("%s/.voyager", dir)
			err = ioutil.WriteFile(fileLocation, targetb, 0666)
			if err != nil {
				fmt.Printf("Error storing IP address to file: %s\n", err)
				return
			}
		}
		// Success
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)

}
