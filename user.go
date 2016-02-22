package slack

type User struct {
	ID        string
	Nick      string
	FirstName string
	LastName  string
}

func (user *User) FullName() (fullName string) {
	fullName = user.FirstName
	if user.LastName != "" {
		if user.FirstName != "" {
			fullName += " "
		}
		fullName += user.LastName
	}
	return
}

func UserFromJSON(data map[string]interface{}) *User {
	id := data["id"].(string)
	nick := data["name"].(string)

	profile := data["profile"].(map[string]interface{})
	var firstName, lastName string
	first, ok := profile["first_name"]
	if ok && first != nil {
		firstName = first.(string)
	}
	last, ok := profile["last_name"]
	if ok && last != nil {
		lastName = last.(string)
	}

	return &User{
		ID:        id,
		Nick:      nick,
		FirstName: firstName,
		LastName:  lastName,
	}
}
