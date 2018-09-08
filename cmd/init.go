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

	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/withmandala/go-log"
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

// logger := log.New(os.Stdout)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the cluster configuration by creating the config.yaml.",
	Long: `Initialize the cluster configuration by creating the config.yaml
with default values.`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := log.New(os.Stdout).WithColor()
		logger.Infof("dumping default config.yaml from INCEPTION image")

		image := imageName("stable", "amd64")
		err := createConfigYaml(image)
		if err != nil {
			logger.Errorf("failed to create the config.yaml")
		}

		logger.Infof("update config with custom values from ~/.grandctl/config.yaml")

		var c conf
		c.getConf()
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

	logger := log.New(os.Stdout).WithColor()

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		logger.Error("yamlFile.Get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		logger.Fatalf("Unmarshal: %v", err)
	}

	c.ChartRepo.AddOns.Header = viper.GetString("HEADER")
	c.ChartRepo.AddOns.URL = viper.GetString("URL")
	c.ImageRepo = viper.GetString("image_repo")
	c.DockerUsername = viper.GetString("docker_username")
	c.DockerPassword = viper.GetString("docker_password")
	c.PrivateRegistryEnable = "true"
	c.AnsibleBecome = viper.GetString("ansible_become")
	c.AnsibleUser = viper.GetString("ansible_user")
	c.PrivateRegistryServer = viper.GetString("private_registry_server")

	// fmt.Println(c)

	// d, err := yaml.Marshal(&c)
	// m := make(map[interface{}]interface{})
	// err = yaml.Unmarshal([]byte(yamlFile), &m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- m:\n%v\n\n", m)

	d, err := yaml.Marshal(&c)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	logger.Infof("---\n%s\n\n", string(d))
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
	logger := log.New(os.Stdout).WithColor()
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
		logger.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		logger.Fatalf("failed to capture stdout or stderr")
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return err
}
