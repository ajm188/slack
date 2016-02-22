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
		fullName += " " + user.LastName
	}
	return
}

func UserFromJSON(data map[string]interface{}) *User {
	id := data["id"].(string)
	nick := data["name"].(string)

	profile := data["profile"].(map[string]interface{})
	firstName := profile["first_name"].(string)
	lastName := profile["last_name"].(string)

	return &User{
		ID:        id,
		Nick:      nick,
		FirstName: firstName,
		LastName:  lastName,
	}
}
