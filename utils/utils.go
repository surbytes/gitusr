package utils

// TODO: Add package documentation comment explaining the purpose of this package

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/surbytes/gitusr/models"
	"gopkg.in/ini.v1"
)

// PrintInfo prints information messages with blue color
// TODO: Use color package consistently instead of ANSI codes directly
// TODO: Add exported function documentation
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// global .gitconfigfile path
func globalConfigFile() string {
	home, err := os.UserHomeDir()
	CheckErr(err)
	return filepath.Join(home, ".gitconfig")
}

// load config file
func loadGlobalConfigFile() *ini.File {
	cfg, err := ini.Load(globalConfigFile())
	CheckErr(err)
	return cfg
}

// get users section keys
func getGlobalUsersKeys() []string {
	file, err := os.Open(globalConfigFile())
	CheckErr(err)
	defer file.Close()

	re := regexp.MustCompile(`\[users\s+"(.*?)"\]`)
	scanner := bufio.NewScanner(file)
	var usersKeys []string

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			usersKeys = append(usersKeys, match[1])
		}
	}

	err = scanner.Err()
	CheckErr(err)

	return usersKeys
}

// prepare users section
func prepareUsers(usersKeys []string) []models.User {
	var users []models.User
	for _, v := range usersKeys {
		u := "users \"" + v + "\""
		section := loadGlobalConfigFile().Section(u)
		name, err := section.GetKey("name")
		CheckErr(err)
		email, err := section.GetKey("email")
		CheckErr(err)
		user := models.User{
			Name:  name.String(),
			Email: email.String(),
		}
		users = append(users, user)
	}

	return append(users, models.GetCurrentUsr())
}

// render users
// TODO: Add exported function documentation
// TODO: Return errors instead of handling them internally
func RenderUsers() {

	var users []string
	var currUser string
	for _, usr := range prepareUsers(getGlobalUsersKeys()) {
		if usr.Name == models.GetCurrentUsr().Name && usr.Email == models.GetCurrentUsr().Email {

			currUser = color.YellowString("%s <%s> *", usr.Name, usr.Email)
			users = append(users, currUser)
		} else {
			users = append(users, fmt.Sprintf("%s <%s>", usr.Name, usr.Email))
		}
	}
	slices.Reverse(users)
	prompt := promptui.Select{
		Label: "Select User",
		Items: users,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if len(result) > 0 {
		rs := strings.Split(result, " ")
		if rs[0] == strings.Split(currUser, " ")[0] && rs[1] == ("<"+models.GetCurrentUsr().Email+">") {
			color.Red("Selected user is already the active Git user. No changes made")
		} else {
			models.SetUsr(rs[0])
		}
	}
}

// handle error
// TODO: Use color package consistently instead of ANSI codes directly
// TODO: Consider returning errors instead of immediate exit for better error handling
// TODO: Add exported function documentation
func CheckErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
