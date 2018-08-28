// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

grandctl install --gate stable`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")

		gate := cmd.Flag("gate").Value.String()
		fmt.Println(" gate:", gate)

		release := cmd.Flag("release").Value.String()
		fmt.Println(" release:", release)

		arch := cmd.Flag("arch").Value.String()
		fmt.Println(" arch:", arch)

		if cmd.Flag("toggle").Value.String() == "true" {
			fmt.Println(" toggle: enabled")
		}

		repo := "hyc-cloud-private-" + gate + "-docker-local.artifactory.swg-devops.com"
		image := repo + "/ibmcom/icp-inception-" + arch + ":" + "latest"

		if gate != "integration" {
			image = repo + "/ibmcom-" + arch + "/icp-inception" + ":" + "latest"
		}

		fmt.Println(" repo:", repo)
		fmt.Println(" image:", image)
		fmt.Println("")
		fmt.Println("")

		fmt.Println("local: install inception")
		uninstall := dockerRunIcpInstall(image)
		if uninstall != nil {
			fmt.Println("Error !!!")
		}

	},
}

// docker run -e LICENSE=accept --net=host -t -v "$(pwd)":/installer/cluster $IMAGE uninstall
func dockerRunIcpInstall(image string) error {
	cmdRunner := exec.Command("docker", "run", "-e", "LICENSE=accept",
		"-e", "http_proxy=${http_proxy}",
		"-e", "https_proxy=${https_proxy}",
		"-e", "no_proxy=${no_proxy}",
		"--net=host",
		"-t",
		"-v", "/opt/ibm/cluster:/installer/cluster", image, "install")
	cmdRunner.Dir = "/opt/ibm/cluster"
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmdRunner.StdoutPipe()
	stderrIn, _ := cmdRunner.StderrPipe()
	cmdRunner.Start()

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	err := cmdRunner.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatalf("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return err
}

// https://github.com/kjk/go-cookbook/blob/master/LICENSE
// func copyAndCapture1(w io.Writer, r io.Reader) ([]byte, error) {
// 	var out []byte
// 	buf := make([]byte, 1024, 1024)
// 	for {
// 		n, err := r.Read(buf[:])
// 		if n > 0 {
// 			d := buf[:n]
// 			out = append(out, d...)
// 			_, err := w.Write(d)
// 			if err != nil {
// 				return out, err
// 			}
// 		}
// 		if err != nil {
// 			// Read returns io.EOF at the end of file, which is not an error for us
// 			if err == io.EOF {
// 				err = nil
// 			}
// 			return out, err
// 		}
// 	}
// }

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")
	uninstallCmd.PersistentFlags().String("gate", "daily", "Gate or Stage to work with, options are: integration, daily, edge, stable")
	uninstallCmd.PersistentFlags().String("release", "3.1.0", "ICP release, options are: latest, 3.1.0, 2.1.0.3-ga, 2.1.0.2-ga")
	uninstallCmd.PersistentFlags().String("arch", "amd64", "architecture, options are: amd64, ppc64le")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
