package dataStructs

import (
	"errors"
)

type Category struct {
    Name    string
    Actions []string
}

var (
    Mobile = Category{
        Name: "Mobile",
        Actions: []string{"Query1", "Query2"}}
    Energy = Category{
        Name: "Energy",
        Actions: []string{"Query1", "Query2"}}
)

var categories = map[string]Category {
    "Mobile": Mobile ,
    "Energy": Energy,
}

func CheckExistence(input string) (*Category, error) {
    if value, found := categories[input]; found == false {
        return nil, errors.New("no category available")
    } else {
        return &value, nil
    }
}
