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

var _ = Describe("Manifest", func() {
	var binLocation string
	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/mcc", runtime.GOOS)
	})

	Context("When there is a valid statefile", func() {

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

			Context("When there is not a manifest file", func() {
				It("UNIT should print error", func() {
					cmd := exec.Command(binLocation, "manifest", "create")
					out, _ := cmd.StdoutPipe()
					cmd.Start()

					buf := new(bytes.Buffer)
					buf.ReadFrom(out)
					Expect(buf.String()).To(Equal("Cannot find voyager manifest file (expecting $HOME/voyager-manifest.json)\n"))
				})
			})

			Context("When there is a manifest file", func() {

				Context("When Creating a manifest", func() {

					type Manifest struct {
						ID          string `json:"id"`
						Environment string `json:"environemnt"`
						DNS         string `json:"dns"`
						IP          string `json:"ip"`
					}
					var manifestFile string

					var data Manifest
					BeforeEach(func() {
						manifestFile = os.Getenv("HOME") + "/voyager-manifest.json"
						data = Manifest{
							ID:          "123456",
							Environment: "ESXi",
							DNS:         "SOME DNS HERE :)",
							IP:          "0.0.0.0",
						}
						manifestData, err := json.Marshal(data)
						Expect(err).ToNot(HaveOccurred())
						err = ioutil.WriteFile(manifestFile, manifestData, 0666)
						Expect(err).ToNot(HaveOccurred())

					})

					AfterEach(func() {
						os.Remove(manifestFile)
					})

					It("UNIT should create the manifest file and return a success", func() {

						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("POST", "/manifest"),
								ghttp.RespondWith(http.StatusOK, nil),
							),
						)

						cmd := exec.Command(binLocation, "manifest", "create", "-f $HOME/voyager-manifest.json")
						out, _ := cmd.StdoutPipe()
						cmd.Start()

						fileContent, err := ioutil.ReadFile(manifestFile)
						Expect(err).ToNot(HaveOccurred())
						jsonData, err := json.Marshal(data)
						Expect(err).ToNot(HaveOccurred())
						Expect(fileContent).To(Equal(jsonData))

						buf := new(bytes.Buffer)
						buf.ReadFrom(out)
						Expect(buf.String()).To(Equal("Manifest file successfully uploaded\n"))
					})
				})
			})
		})
	})
})
