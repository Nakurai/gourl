package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/nakurai/gourl/db"
	"gorm.io/gorm"
)

var CurrentEnv *Environment

type Environment struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Variables   JSONMap `gorm:"type:json"`
	Current     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e Environment) String() string {
	res := e.Name
	if e.Current{
		res = "* "+res
	}
	if e.Description != "" {
		res += " - " + e.Description
	}
	return res
}

func InitEnvironment() error {

	existDefault, err := GetEnv("default")
	if err != nil {
		return err
	}
	if existDefault == nil {
		defaultEnv := Environment{
			Name:        "default",
			Description: "default environment created on first execution",
			Variables:   map[string]string{},
			Current:     true,
		}
		return CreateEnv(&defaultEnv)
	}

	curEnv, err := GetCurrentEnv()
	if err != nil {
		return err
	}
	CurrentEnv = curEnv

	return nil

}

func CreateEnv(env *Environment) error {
	res := db.Db.Create(env)
	if res.Error != nil {
		return fmt.Errorf("while saving the new environment %s: %v", env.Name, res.Error)
	}
	return nil
}

func GetEnv(name string) (*Environment, error) {

	var existingEnv Environment
	res := db.Db.Where("name = ?", name).First(&existingEnv)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("while fetching existing environment %s: %v", name, res.Error)
		}
	}
	return &existingEnv, nil
}

func GetCurrentEnv() (*Environment, error) {

	var existingEnv Environment
	res := db.Db.Where("current = true").First(&existingEnv)
	if res.Error != nil {
		return nil, fmt.Errorf("while fetching current environment: %v", res.Error)
	}
	return &existingEnv, nil
}

func LoadEnv(name string) error{
	curEnv, err := GetEnv(name)
	if err != nil {
		return err
	}
	if curEnv == nil{
		return fmt.Errorf("Environment named %s does nto exist", name)
	}
	curEnv.Current = true
	db.Db.Save(curEnv)
	CurrentEnv = curEnv

	return nil
}
