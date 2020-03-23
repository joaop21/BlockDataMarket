package main

import (
    "errors"
)

type Category string

const (
    Mobile Category = "Mobile"
    Energy Category = "Energy"
)

var categories = map[string]Category {
    "Mobile": Mobile,
    "Energy": Energy,
}

func checkExistence(input string) (string, error) {
    if value, found := categories[input]; found == false {
        return "", errors.New("no key found")
    } else {
        return string(value), nil
    }
}
