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
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	AnsibleUser      string `json:"ansible_user,omitempty" yaml:"ansible_user,omitempty"`
	AnsibleBecome    string `json:"ansible_become,omitempty" yaml:"ansible_become,omitempty"`
	ClusterLbAddress string `json:"cluster_lb_address,omitempty" yaml:"cluster_lb_address,omitempty"`
	ClusterVip       string `json:"cluster_vip,omitempty" yaml:"cluster_vip,omitempty"`
	ChartRepo        struct {
		AddOns struct {
			Header string `json:"header" yaml:"header"`
			URL    string `json:"url" yaml:"url"`
		} `json:"addons,omitempty" yaml:"addons,omitempty"`
	} `json:"chart_repo,omitempty" yaml:"chart_repo,omitempty"`
	DefaultAdminUser     string   `json:"default_admin_user,omitempty" yaml:"default_admin_user,omitempty"`
	DefaultAdminPassword string   `json:"default_admin_password,omitempty" yaml:"default_admin_password,omitempty"`
	DockerUsername       string   `json:"docker_username,omitempty" yaml:"docker_username,omitempty"`
	DockerPassword       string   `json:"docker_password,omitempty" yaml:"docker_password,omitempty"`
	EtcdExtraArgs        []string `json:"etcd_extra_args" yaml:"etcd_extra_args"`
	ImageRepo            string   `json:"image_repo,omitempty" yaml:"image_repo,omitempty"`
	Managementservices   struct {
		Istio                string `json:"istio" yaml:"istio"`
		VulnerabilityAdvisor string `json:"vulnerability-advisor" yaml:"vulnerability-advisor"`
		StorageGlusterfs     string `json:"storage-glusterfs" yaml:"storage-glusterfs"`
		StorageMinio         string `json:"storage-minio" yaml:"storage-minio"`
	} `json:"management_services" yaml:"management_services"`
	NetworkType           string `json:"network_type" yaml:"network_type"`
	PrivateRegistryEnable string `json:"private_registry_enable,omitempty" yaml:"private_registry_enable,omitempty"`
	PrivateRegistryServer string `json:"private_registry_server,omitempty" yaml:"private_registry_server,omitempty"`
	ProxyLbAddress        string `json:"proxy_lb_address,omitempty" yaml:"proxy_lb_address,omitempty"`
	ProxyVipIface         string `json:"proxy_vip_iface,omitempty" yaml:"proxy_vip_iface,omitempty"`
	ProxyVip              string `json:"proxy_vip,omitempty" yaml:"proxy_vip,omitempty"`
	VipManager            string `json:"vip_manager" yaml:"vip_manager"`
	VipIface              string `json:"vip_iface,omitempty" yaml:"vip_iface,omitempty"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("check for cluster config.yaml exists")
		fmt.Printf("\tdocker_username from config: %s\n", viper.Get("docker_username"))

		// image := imageName("stable", "amd64")
		// err := createConfigYaml(image)
		// if err != nil {
		// 	fmt.Println("Error !!!")
		// }
		var c conf

		c.getConf()

		// fmt.Println(configCmd)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// marshal in the default INCEPTION config.yaml
// and then add our custom settings based on our
// .grandctl/config.yaml
//
func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	//
	// TODO: we're hardcoding this here now just for development work,
	// but this really needs to be read from the .grandctl/config.yaml
	// and updated to the config.yaml
	//
	c.ChartRepo.AddOns.Header = viper.GetString("HEADER")
	c.ChartRepo.AddOns.URL = viper.GetString("URL")
	c.ImageRepo = viper.GetString("image_repo")
	c.DockerUsername = viper.GetString("docker_username")
	c.DockerPassword = viper.GetString("docker_password")
	c.PrivateRegistryEnable = "true"
	c.AnsibleBecome = viper.GetString("ansible_become")
	c.AnsibleUser = viper.GetString("ansible_user")
	c.PrivateRegistryServer = viper.GetString("private_registry_server")

	fmt.Println(c)

	// d, err := yaml.Marshal(&c)
	// m := make(map[interface{}]interface{})
	// err = yaml.Unmarshal([]byte(yamlFile), &m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- m:\n%v\n\n", m)

	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
	return c
}

// construct the appropriate artifactory url and return
//
func imageName(gate string, arch string) string {
	repo := "hyc-cloud-private-" + gate + "-docker-local.artifactory.swg-devops.com"
	imageName := repo + "/ibmcom/icp-inception-" + arch + ":latest"

	if gate != "integration" {
		imageName = repo + "/ibmcom-" + arch + "/icp-inception:latest"
	}
	return imageName
}

// dump the default INCEPTION config.yaml from the docker image
//
// sudo docker run -e LICENSE=accept -v "$(pwd)":/data ibmcom/icp-inception:2.1.0.3 cp -r cluster /data
func createConfigYaml(image string) error {
	cmdRunner := exec.Command("docker", "run", "-e",
		"LICENSE=accept",
		"--net=host", "-t",
		"-v", "/opt/ibm/cluster:/data", image, "cp cluster/config.yaml /data")
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
