/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"agenda/entity"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called")
		_username,_ := cmd.Flags().GetString("user")
		_password,_ := cmd.Flags().GetString("password")
		_email,_ := cmd.Flags().GetString("email")
		_phone,_ := cmd.Flags().GetString("phone")
		fmt.Println("username:"+_username+" password:"+_password+" email:"+_email+" phone:"+_phone)
		if _username == "default user" {
			fmt.Println("username can not be empty")
			return 
		}
		users := entity.ReadUsers()
		for i := 0; i < len(users); i++ {
			if users[i].Username == _username {
				fmt.Println("this username has been registered")
				return
			}
		}
		newUser := entity.User{_username,_password,_email,_phone}
		users = append(users,newUser)
		entity.WriteUsers(users)
	},
}

func init() {
	registerCmd.Flags().StringP("user","u","default user","Help message for username")
	registerCmd.Flags().StringP("password","p","123456","password of user")
	registerCmd.Flags().StringP("email","e","123456@163.com","email of user")
	registerCmd.Flags().StringP("phone","t","12345647897","phone number of user")//telephone
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
