package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olenedr/esamarathon/config"
	"github.com/olenedr/esamarathon/db"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	r "gopkg.in/gorethink/gorethink.v3"
)

const Table = "users"

type User struct {
	ID       string `gorethink:"id,omitempty" json:"id,omitempty"`
	Username string `gorethink:"username" json:"user_name,omitempty"`
}

type TwitchResponse struct {
	User       User `json:"token,omitempty"`
	Identified bool `json:"identified,omitempty"`
}

func Insert(username string) error {
	var data = map[string]interface{}{
		"username": username,
	}

	return db.Insert(Table, data)
}

func All() ([]User, error) {
	rows, err := db.GetAll(Table)
	var users []User
	if err != nil {
		return users, errors.Wrap(err, "user.All")
	}

	if err = rows.All(&users); err != nil {
		return users, errors.Wrap(err, "user.All")
	}

	return users, nil
}

// Exists checks the db for the user by Username
func (u *User) Exists() (bool, error) {
	res, err := r.Table(Table).Filter(map[string]interface{}{
		"username": u.Username,
	}).Run(db.Session)

	if err != nil {
		return false, err
	}

	defer res.Close()

	var rows []interface{}
	res.All(&rows)
	if len(rows) == 0 {
		return false, nil
	}

	return true, nil
}

func GetUserByUsername(username string) (User, error) {
	res, err := r.Table(Table).Filter(map[string]interface{}{
		"username": username,
	}).Run(db.Session)

	var u User
	if err != nil {
		return u, errors.Wrap(err, "user.GetUserByUsername")
	}

	defer res.Close()

	if err = res.One(&u); err != nil {
		return u, err
	}

	return u, nil
}

func RequestTwitchUser(token *oauth2.Token) (User, error) {
	c := &http.Client{}
	var res TwitchResponse

	req, err := http.NewRequest("GET", config.Config.TwitchAPIRootURL, nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", "OAuth "+token.AccessToken)
	req.Header.Add("Client-ID", config.Config.TwitchClientID)
	resp, err := c.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)

	// Check for err or empty struct
	if err != nil || res.User == (User{}) {
		fmt.Printf("%v", err)
		return User{}, err
	}

	return res.User, nil

}
