package models

// TODO: Add package documentation comment explaining the purpose of this package

import (
	"log"
	"os/exec"
	"strings"
)

// TODO: Consider using constants for config keys instead of string literals

// User represents a Git user profile with name and email
type User struct {
	Name  string
	Email string
}

// add new user into git config file
// TODO: Follow Go naming convention - rename to AddUser
func AddUsr(name string, email string) (User, error) {
	sectionName := "users." + name + ".name"
	sectionEmail := "users." + name + ".email"
	if err := exec.Command("git", "config", "--global", "--add", sectionEmail, email).Run(); err != nil {
		return User{}, err
	}

	if err := exec.Command("git", "config", "--global", "--add", sectionName, name).Run(); err != nil {
		return User{}, err
	}

	return User{
		Name:  name,
		Email: email,
	}, nil
}

// GetUsr retrieves a user profile from the git config
// TODO: Follow Go naming convention - rename to GetUser
func GetUsr(name string) (User, error) {
	sectionName := "users." + name + ".name"
	sectionEmail := "users." + name + ".email"

	n, err := exec.Command("git", "config", "--global", sectionName).Output()
	if err != nil {
		return User{}, err
	}

	email, err := exec.Command("git", "config", "--global", sectionEmail).Output()

	if err != nil {
		return User{}, err
	}
	nemail := strings.Trim(string(email), "\f\t\r\n ")
	nname := strings.Trim(string(n), "\f\t\r\n ")

	return User{
		Name:  nname,
		Email: nemail,
	}, nil

}

// SetUsr activates a user profile as the current Git user
// TODO: Follow Go naming convention - rename to SetUser
// TODO: Return errors instead of using log.Fatal for better error handling
func SetUsr(name string) {

	sectionName := "users." + name + ".name"
	sectionEmail := "users." + name + ".email"

	targetUsr, err := GetUsr(name)
	if err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	currentUsr := GetCurrentUsr()

	// ADD the current user into users.

	_, err = AddUsr(currentUsr.Name, currentUsr.Email)
	if err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	if err := exec.Command("git", "config", "--global", "user.name", targetUsr.Name).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	if err := exec.Command("git", "config", "--global", "user.email", targetUsr.Email).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	if err := exec.Command("git", "config", "--global", "--unset", sectionName).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	if err := exec.Command("git", "config", "--global", "--unset", sectionEmail).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}
}

// get the current user on git config file
// TODO: Follow Go naming convention - rename to GetCurrentUser
// TODO: Return error instead of using log.Fatal
func GetCurrentUsr() User {

	// TODO: Remove commented code
	//var email strings.Builder
	//var name strings.Builder
	name, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	email, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	nname := strings.Trim(string(name), "\f\t\r\n ")
	nemail := strings.Trim(string(email), "\f\t\r\n ")

	return User{
		Name:  nname,
		Email: nemail,
	}

}

// delete user from git config file
// TODO: Follow Go naming convention - rename to DeleteUser
// TODO: Return error instead of using log.Fatal
// TODO: Clarify function signature - unused parameters (name, email)
func (usr *User) DelUsr(name string, email string) {

	if err := exec.Command("git", "config", "--global", "--unset-all", "users.email", usr.Email).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}

	if err := exec.Command("git", "config", "--global", "--unset-all", "users.name", usr.Name).Run(); err != nil {
		log.Fatal(err) // TODO: Return error instead of fatal exit
	}
}
