package users

import (
	"datastore-gin-example/common"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var client = common.DatastoreClient
var ctx = common.Ctx

type User struct {
	Key  *datastore.Key `datastore:"__key__" json:"id"`
	Name string         `json:"name"`
	Age  int            `json:"age"`
}

type LoginAccount struct {
	Account  string `datastroe:"account" json:"account"`
	Password string `datastore:"password" json:"password"`
	UserID   string `datastore:"user_id" json:"userID"`
	Actived  bool   `datastore:"actived" json:"actived"`
}

func GetAllUsers() ([]*User, error) {
	users := make([]*User, 0)
	q := datastore.NewQuery("User")
	_, err := client.GetAll(ctx, q, &users)
	return users, err
}

func GetUserByID(userID string) (*User, error) {
	var err error
	key, err := datastore.DecodeKey(userID)
	user := &User{}
	err = client.Get(ctx, key, user)
	return user, err
}

func InsertUser(user *User) error {
	var err error
	var keys []*datastore.Key
	keys = append(keys, datastore.IncompleteKey("User", nil))
	keys, err = client.AllocateIDs(ctx, keys)
	_, err = client.Mutate(ctx,
		datastore.NewInsert(keys[0], &User{
			Name: user.Name,
			Age:  user.Age,
		}),
	)
	return err
}

func UpdateUser(userID string, user *User) error {
	var err error
	key, err := datastore.DecodeKey(userID)
	_, err = client.Mutate(ctx,
		datastore.NewUpdate(key, &User{
			Name: user.Name,
			Age:  user.Age,
		}),
	)
	return err
}

func DeleteUser(userID string) error {
	var err error
	key, err := datastore.DecodeKey(userID)
	err = client.Delete(ctx, key)
	return err
}

func UserLogin(loginAccount *LoginAccount) (*User, *LoginAccount, error) {
	query := datastore.NewQuery("LoginAccount").
		Filter("account =", loginAccount.Account).
		Filter("password =", loginAccount.Password).
		Filter("actived =", true)
	it := client.Run(ctx, query)
	var searchAccount = LoginAccount{}
	var err error
	for {
		_, err := it.Next(&searchAccount)
		if err == iterator.Done {
			break
		}

	}
	if (LoginAccount{}) == searchAccount || searchAccount.Actived == false {
		return nil, nil, err
	}
	common.LogInfo("search user by id" + searchAccount.UserID)
	user, err := GetUserByID(searchAccount.UserID)
	if (&User{}) == user {
		return nil, nil, err
	}
	return user, loginAccount, err
}
