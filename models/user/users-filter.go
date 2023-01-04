package models

type UsersFilter struct {
	IDFilter        []ID        `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstNameFilter []FirstName `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastNameFilter  []LastName  `json:"last_name,omitempty" bson:"last_name,omitempty"`
	NicknameFilter  []Nickname  `json:"nickname,omitempty" bson:"nickname,omitempty"`
	EmailFilter     []Email     `json:"email,omitempty" bson:"email,omitempty"`
	CountryFilter   []Country   `json:"country,omitempty" bson:"country,omitempty"`
}
