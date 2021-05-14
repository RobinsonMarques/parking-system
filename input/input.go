package input

import "github.com/RobinsonMarques/parking-system/database"

type CreateUserInput struct {
	Person   database.Person
	Document string  `json:"Document" binding:"required"`
	Balance  float64 `json:"Balance"`
}
