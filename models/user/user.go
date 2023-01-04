package models

type User struct {
	ID        `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Nickname  `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Password  `json:"password,omitempty" bson:"password,omitempty"`
	Email     `json:"email,omitempty" bson:"email,omitempty"`
	Country   `json:"country,omitempty" bson:"country,omitempty"`
	CreatedAt `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (u *User) GetID() interface{} {
	return u.ID
}

func (u *User) SetID() interface{} {
	return u.ID
}
