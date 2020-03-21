package main

type Category uint

const (
    Mobile Category = iota + 1
    Energy
)

func (category Category) String() string {
    names := []string{
        "Mobile",
        "Energy"
    }

    if category < Mobile || day > Energy {
      return "Unknown"
    }

    return names[day]
}
