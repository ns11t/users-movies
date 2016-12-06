package model

import (
	"errors"
	"regexp"
)

var PhoneBaseRegexp *regexp.Regexp = regexp.MustCompile(`^[+]?([\d]{5,20})$`)

// Represents User entity for API
type Person struct {
	Id          int    `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phoneNumber"`
}

func (p *Person) ValidateLoginPassword() error {
	if p.Login == "" {
		return errors.New("login is incorrect")
	}
	if p.Password == "" {
		return errors.New("password is incorrect")
	}
	return nil
}

func (p *Person) Validate() error {
	err := p.ValidateLoginPassword()
	if err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("name is incorrect")
	}
	if p.Age < 1 {
		return errors.New("age is incorrect")
	}
	if !PhoneBaseRegexp.MatchString(p.PhoneNumber) {
		return errors.New("phone is incorrect")
	}
	return nil
}
