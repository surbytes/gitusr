package models

import (
	"log"
	"os/exec"
	"strings"
)

type User struct {
	Name  string
	Email string
}

// add new user into git config file.
func AddUsr(name, email string) (User, error) {
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

func SetUsr(name string) {
	sectionName := "users." + name + ".name"
	sectionEmail := "users." + name + ".email"

	targetUsr, err := GetUsr(name)
	if err != nil {
		log.Fatal(err)
	}

	currentUsr := GetCurrentUsr()

	// ADD the current user into users.

	_, err = AddUsr(currentUsr.Name, currentUsr.Email)
	if err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("git", "config", "--global", "user.name", targetUsr.Name).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("git", "config", "--global", "user.email", targetUsr.Email).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("git", "config", "--global", "--unset", sectionName).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("git", "config", "--global", "--unset", sectionEmail).Run(); err != nil {
		log.Fatal(err)
	}
}

// get the current user on git config file.
func GetCurrentUsr() User {
	// var email strings.Builder
	// var name strings.Builder
	name, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		log.Fatal(err)
	}

	email, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		log.Fatal(err)
	}

	nname := strings.Trim(string(name), "\f\t\r\n ")
	nemail := strings.Trim(string(email), "\f\t\r\n ")

	return User{
		Name:  nname,
		Email: nemail,
	}
}

// delete user from git config file.
func (usr *User) DelUsr(name, email string) {
	if err := exec.Command("git", "config", "--global", "--unset-all", "users.email", usr.Email).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("git", "config", "--global", "--unset-all", "users.name", usr.Name).Run(); err != nil {
		log.Fatal(err)
	}
}
