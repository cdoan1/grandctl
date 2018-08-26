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

	"github.com/outten45/aini"
	"github.com/spf13/cobra"
)

// iniCmd represents the ini command
var iniCmd = &cobra.Command{
	Use:   "ini",
	Short: "parse the HOSTS file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ini called")
		v, _ := aini.NewFile("hosts.ini")

		// fmt.Println(v.Groups["worker"])
		// fmt.Println(v.Groups["master"])

		for k := range v.Groups {
			if k == "ungrouped" {
				continue
			}
			fmt.Println(k)
			hosts := v.Groups[k]
			for host := range hosts {
				fmt.Println(hosts[host].Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(iniCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iniCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iniCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
