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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Manifest for voyager configuration deployment.",
	Long: `A manifest that voyager uses to deploy and configure infrastructure.
	Usage: manifest [create|retrieve|update|delete]`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(manifestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// CheckMccTarget is a helper function that reads mcc target
func CheckMccTarget() (url.URL, error) {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	fileLocation := fmt.Sprintf("%s/.voyager", dir)

	url := url.URL{}
	// Read IP target from file
	_, err = os.Stat(fileLocation)
	if err != nil {
		color.Red("Cannot find mcc config file (expecting $HOME/.voyager)\n")
		return url, err
	}

	fileContent, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal data and print
	err = json.Unmarshal(fileContent, &url)
	if err != nil {
		log.Fatal(err)
		return url, err
	}
	return url, nil
}
