package model

//Movie this movie model
// swagger:model Movie
type Movie struct {
	// the Uid for this user
	// required: true
	Uid string `bson:"uid,omitempty" json:"uid"`
	// the name for this user
	// required: true
	Name string `bson:"name,omitempty" json:"name"`
	// the budget for this user
	// required: true
	Budget int `bson:"budget,omitempty" json:"budget"`
	// the director for this user
	// required: true
	Director string `json:"director"`
}
