package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// NodesResponse is ...
type NodesResponse struct {
	ID     string `json:"ID"`
	Type   string `json:"Type"`
	Status string `json:"Status"`
	IP     string `json:"IP"`
}

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List nodes' information",
	Long: `List nodes' ID, Type, Status
Usage: mcc nodes`,

	Run: func(cmd *cobra.Command, args []string) {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		urlObject := url.URL{}
		fileLocation := fmt.Sprintf("%s/.voyager", dir)

		// Read IP target from file
		_, err = os.Stat(fileLocation)
		if err != nil {
			color.Red("Cannot find mcc config file (expecting $HOME/.voyager)\n")
			return
		}

		fileContent, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshal data and print
		err = json.Unmarshal(fileContent, &urlObject)
		if err != nil {
			log.Fatal(err)
		}

		// Send http GET request to IP target
		targetURL := fmt.Sprintf("%s://%s/nodes", urlObject.Scheme, urlObject.Host)
		res, err := http.Get(targetURL)
		if err != nil {
			color.Red("Error sending '/nodes' API call to Voyager at %s: %s\n", urlObject.Host, err)
			return
		}

		// Parse response
		responseBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			color.Red("Error reading response body: %s\n", err)
			return
		}

		if res.StatusCode != 200 {
			color.Red("Server returned (%d): %s\n", res.StatusCode, responseBytes)
			return
		}

		jsonToTable(responseBytes)
	},
}

func jsonToTable(jsonIn []byte) {
	var nodes []NodesResponse
	err := json.Unmarshal(jsonIn, &nodes)
	if err != nil {
		fmt.Println("error:", err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Type", "Status", "IP"})

	for _, v := range nodes {
		table.Append([]string{v.ID, v.Type, v.Status, v.IP})
	}

	table.Render()
}

func init() {
	RootCmd.AddCommand(nodesCmd)
}
