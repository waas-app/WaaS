package model

import "strings"

type User struct {
	ID                uint   `gorm:"primary_key;autoIncrement" json:"-"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Slug              string `gorm:"default:md5((random())::text)" json:"id"`
	Admin             bool   `json:"admin,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (u *User) GetPID() string {
	//if email is present we use email else we use phone
	if u.Email != "" {
		return u.Email
	}
	return ""
}

func (u *User) PutPID(pid string) {
	//need to check if it is email or phone
	if strings.Contains(pid, "@") {
		u.PutEmail(pid)
	}
}

func (u *User) GetPassword() (password string) {
	return u.EncryptedPassword
}

func (u *User) PutPassword(password string) {
	u.EncryptedPassword = password
}

// GetArbitrary is used only to display the arbitrary data back to the user
// when the form is reset.
func (u *User) GetArbitrary() (arbitrary map[string]string) {
	return
}

// PutArbitrary allows arbitrary fields defined by the authboss library
// consumer to add fields to the user registration piece.
func (u *User) PutArbitrary(arbitrary map[string]string) {
	if val, ok := arbitrary["name"]; ok {
		u.Username = val
	}
}

func (u *User) PutEmail(email string) { u.Email = email }
