package main

import (
    "errors"
)

type Category struct {
    name    string
    actions []string
}

var (
    Mobile = Category{
        name: "Mobile",
        actions: []string{"Query1", "Query2"}}
    Energy = Category{
        name: "Energy",
        actions: []string{"Query1", "Query2"}}
)

var categories = map[string]Category {
    "Mobile": Mobile ,
    "Energy": Energy,
}

func checkExistence(input string) (*Category, error) {
    if value, found := categories[input]; found == false {
        return nil, errors.New("no category available")
    } else {
        return &value, nil
    }
}
