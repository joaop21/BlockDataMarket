package main

type Category uint

const (
    Mobile Category = iota
    Energy
)

func (category Category) String() string {
    names := []string{
        "Mobile",
        "Energy",
    }

    if category < Mobile || category > Energy {
        return "Unknown"
    }

    return names[category]
}
