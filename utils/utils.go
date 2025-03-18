package utils

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
func RenderUsers(username ...string) {

	if len(username) > 0 && username[0] != "" {
		err := SwitchUserByName(username[0])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		return
	}

	var users []string
	currentUser := models.GetCurrentUsr()
	for _, usr := range prepareUsers(getGlobalUsersKeys()) {
		if usr.Name == currentUser.Name && usr.Email == currentUser.Email {
			users = append(users, color.YellowString("%s <%s> *", usr.Name, usr.Email))
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
		selectedUser := sanitizeResult(result)
		if currentUser.Name == selectedUser.Name && currentUser.Email == selectedUser.Email {
			color.Red("Selected user is already the active Git user. No changes made")
		} else {
			models.SetUsr(currentUser, selectedUser)
		}
	}
}

// Sanitize result with regex
func sanitizeResult(result string) models.User {
	ansiReg := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	reg := regexp.MustCompile(`\<(.*?)\>`)
	email := reg.FindStringSubmatch(result)[1]

	name := reg.ReplaceAllString(result, "")
	if ansiReg.MatchString(name) {
		name = strings.ReplaceAll(name, " *", "")
		name = ansiReg.ReplaceAllString(name, "")
	}
	name = strings.TrimSpace(name)

	return models.User{
		Email: email,
		Name:  name,
	}
}
// SwitchUserByName provides a quick way to switch Git users by name
// without having to go through the interactive selection process
func SwitchUserByName(username string) error {
	usersKeys := getGlobalUsersKeys()
	
	var foundUser models.User
	var found bool
	
	for _, key := range usersKeys {
		if strings.EqualFold(key, username) {
			u := "users \"" + key + "\""
			section := loadGlobalConfigFile().Section(u)
			
			name, err := section.GetKey("name")
			if err != nil {
				return fmt.Errorf("user found but missing name property")
			}
			
			email, err := section.GetKey("email")
			if err != nil {
				return fmt.Errorf("user found but missing email property")
			}
			
			foundUser = models.User{
				Name:  name.String(),
				Email: email.String(),
			}
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("user '%s' not found", username)
	}

	currentUser := models.GetCurrentUsr()
	if currentUser.Name == foundUser.Name && currentUser.Email == foundUser.Email {
		color.Yellow("User '%s' is already the active Git user", username)
		return nil
	}
	
	models.SetUsr(currentUser, foundUser)
	PrintInfo("Switched to Git user: %s <%s>", foundUser.Name, foundUser.Email)
	return nil
}

// handle error
func CheckErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
