package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/surbytes/gitusr/models"
	"github.com/surbytes/gitusr/utils"
	"gopkg.in/ini.v1"
)

func main() {
	home, err := os.UserHomeDir()
	utils.CheckErr(err)

	gitconfig := filepath.Join(home, ".gitconfig")

	file, err := os.Open(gitconfig)
	utils.CheckErr(err)
	defer file.Close()

	if len(os.Args) > 1 {
		models.SetUsr(os.Args[1])
	}
	//re := regexp.MustCompile(`\[users\.(.*?)\]`)
	re := regexp.MustCompile(`\[users\s+"(.*?)"\]`)

	var usersKeys []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			usersKeys = append(usersKeys, match[1])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file")
		return
	}

	//fmt.Println(usersKeys)

	initdata, err := ini.Load(gitconfig)
	utils.CheckErr(err)
	var users []models.User
	for _, v := range usersKeys {
		u := "users \"" + v + "\""
		//fmt.Println(u)
		section := initdata.Section(u)
		fmt.Println("before errr")
		name, err := section.GetKey("name")
		utils.CheckErr(err)
		email, err := section.GetKey("email")
		utils.CheckErr(err)
		user := models.User{
			Name:  name.String(),
			Email: email.String(),
		}
		users = append(users, user)
	}

	fmt.Println(users)
	users = append(users, models.GetCurrentUsr())
	fmt.Println(users)

	for _, usr := range users {
		//utils.PrintInfo("%s <%s>", usr.Name, usr.Email)
		if usr.Name == models.GetCurrentUsr().Name && usr.Email == models.GetCurrentUsr().Email {
			color.Yellow("%s <%s> *", usr.Name, usr.Email)
		} else {
			fmt.Printf("%s <%s>\n", usr.Name, usr.Email)
		}
	}

}
