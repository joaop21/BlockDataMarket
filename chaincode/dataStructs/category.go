package dataStructs

// Category object that represents a category in the World State
type Category struct {
    Type               string     `json:"type"`
    Name               string     `json:"name"`
    PossibleQueries    []string   `json:"possibleQueries"`
}

// Constructor for Category
func NewCategory(name string, queries []string) *Category {

    return &Category{
        Type:            "Category",
        Name:            name,
        PossibleQueries: queries,
    }

}