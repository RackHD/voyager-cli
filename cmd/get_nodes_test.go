package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("GetNodes", func() {
	var binLocation string

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/mcc", runtime.GOOS)
	})

	Context("When there is no statefile", func() {
		It("UNIT should print error", func() {
			cmd := exec.Command(binLocation, "nodes")
			out, _ := cmd.StdoutPipe()
			cmd.Start()

			buf := new(bytes.Buffer)
			buf.ReadFrom(out)
			Expect(buf.String()).To(Equal("Cannot find mcc config file (expecting $HOME/.voyager)\n"))
		})
	})

	Context("When there is a valid statefile", func() {
		Context("When the target URL is invalid", func() {
			It("UNIT should print error", func() {
				server := ghttp.NewServer()
				target, _ := url.Parse(server.URL())
				target.Host = "0.0.0.1"
				targetb, err := json.Marshal(target)
				Expect(err).ToNot(HaveOccurred())

				stateFile := os.Getenv("HOME") + "/.voyager"
				err = ioutil.WriteFile(stateFile, targetb, 0666)
				Expect(err).ToNot(HaveOccurred())

				cmd := exec.Command(binLocation, "nodes")
				out, _ := cmd.StdoutPipe()
				cmd.Start()

				// Capture Standard Output to verify
				buf := new(bytes.Buffer)
				buf.ReadFrom(out)
				Expect(buf.String()).To(ContainSubstring("Error sending '/nodes' API call to Voyager at 0.0.0.1: Get http://0.0.0.1/nodes:"))
				os.Remove(stateFile)
			})
		})

		Context("When the target URL is valid", func() {
			var server *ghttp.Server
			var stateFile string

			BeforeEach(func() {
				server = ghttp.NewServer()

				target, _ := url.Parse(server.URL())
				targetb, err := json.Marshal(target)
				Expect(err).ToNot(HaveOccurred())

				stateFile = os.Getenv("HOME") + "/.voyager"
				err = ioutil.WriteFile(stateFile, targetb, 0666)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
				os.Remove(stateFile)
			})

			Context("When mcc sends get nodes to inventory service", func() {
				Context("When response is not empty", func() {
					It("UNIT receive a list of nodes information", func() {
						expectedResponse := `[{"ID":"UNIQUE_ID","Type":"compute","Status":"Added","IP":"1.2.3.4"},{"ID":"UNIQUE_ID1","Type":"compute","Status":"Added","IP":"2.3.4.5"}]`
						expectedString := `+------------+---------+--------+---------+
|     ID     |  TYPE   | STATUS |   IP    |
+------------+---------+--------+---------+
| UNIQUE_ID  | compute | Added  | 1.2.3.4 |
| UNIQUE_ID1 | compute | Added  | 2.3.4.5 |
+------------+---------+--------+---------+
`
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/nodes"),
								ghttp.RespondWith(http.StatusOK, expectedResponse),
							),
						)

						cmd := exec.Command(binLocation, "nodes")
						out, _ := cmd.StdoutPipe()
						cmd.Start()

						buf := new(bytes.Buffer)
						buf.ReadFrom(out)
						Expect(server.ReceivedRequests()).To(HaveLen(1))
						Expect(buf.String()).To(Equal(expectedString))
					})
				})

				Context("When response is empty", func() {
					It("UNIT should show empty table", func() {
						ExpectedResponse := ""
						expectedString := `+----+------+--------+----+
| ID | TYPE | STATUS | IP |
+----+------+--------+----+
+----+------+--------+----+
`
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/nodes"),
								ghttp.RespondWith(http.StatusOK, ExpectedResponse),
							),
						)

						cmd := exec.Command(binLocation, "nodes")
						out, _ := cmd.StdoutPipe()
						cmd.Start()

						buf := new(bytes.Buffer)
						buf.ReadFrom(out)
						Expect(buf.String()).To(ContainSubstring(expectedString))
					})
				})
			})
		})
	})
})
