package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct{
	Username	string
	Password	string
	Email		string
	Phone 		string
}

func checkerr(err error){
	if err != nil {
		fmt.Println(err)
	}
}

func GetCurrentUserName()(username string){
	dir,err := os.Getwd()
	checkerr(err)
	b,err := ioutil.ReadFile(dir+"/entity/currentUser.txt")
	checkerr(err)
	username = string(b)
	return username
}

func SetCurrentUserName(username string){
	dir,err := os.Getwd()
	checkerr(err)
	b := []byte(username)
	err = ioutil.WriteFile(dir+"/entity/currentUser.txt",b,0777)
	checkerr(err)
}

func ReadUsers()(user []User){
	dir,err := os.Getwd()
	checkerr(err)
	b,err := ioutil.ReadFile(dir+"/entity/Users.txt")
	var users []User
	json.Unmarshal(b,&users)
	return users
}

func WriteUsers(users []User){
	dir,err := os.Getwd()
	checkerr(err)
	data,err := json.Marshal(users)
	checkerr(err)
	b := []byte(data)
	err = ioutil.WriteFile(dir+"/entity/Users.txt",b,0777)
	checkerr(err)
}