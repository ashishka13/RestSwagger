package model

//Movie this movie model
// swagger:model Movie
type Movie struct {
	// the name for this user
	// required: true
	Name string `json:"name"`
	// the budget for this user
	// required: true
	Budget int `json:"budget"`
	// the director for this user
	// required: true
	Director string `json:"director"`
}
