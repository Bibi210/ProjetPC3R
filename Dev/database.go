package main

import (
	/* 	"database/sql" */
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
)

/* KeyUserName */
type UserLOGININFO struct {
	password       string
	currentSession int64
}

/* KeyUserName */
type UserProfile struct {
	Connection int
}

type SimpleClaims struct {
	Username  string
	ExpiresAt time.Time
	SessionId int64
}

func (c SimpleClaims) Valid() error {
	if c.ExpiresAt == (time.Time{}) {
		return OnlyServerError("Token don't have an expiration date")
	}
	if c.ExpiresAt.Before(time.Now()) {
		return OnlyServerError("Token is expired")
	}
	user, ok := loginMAP[c.Username]
	if !ok {
		return OnlyServerError(fmt.Sprintf("Invalid Username %s", c.Username))
	}
	if user.currentSession != c.SessionId {
		return OnlyServerError(fmt.Sprintf("The Session Used is Invalid %s", c.Username))
	}
	_, ok = profilMap[c.Username]
	if !ok {
		return OnlyServerError(fmt.Sprintf("User without profile %s", c.Username))
	}
	return nil
}

var loginMAP = map[string]UserLOGININFO{}
var profilMap = map[string]UserProfile{}

var serverKey = []byte("This is a fun serverkey")

func tokenToString(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString(serverKey)
	if err != nil {
		return tokenString, ServerRuntimeError("Error While Converting JWT Token to string", err)
	}
	return tokenString, nil
}

const sessionDuration = 1 * time.Hour

func createToken(username string) (*jwt.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, SimpleClaims{})
	expirationTime := time.Now().Add(sessionDuration)

	token.Claims = SimpleClaims{ExpiresAt: expirationTime, Username: username, SessionId: expirationTime.Unix()}
	user := loginMAP[username]
	user.currentSession = expirationTime.Unix()
	loginMAP[username] = user
	return token, token.Claims.Valid()
}

func loginAccount(auth AuthJSON) (string, error) {
	if loginMAP[auth.Login].password != auth.Mdp {
		return "", OnlyServerError("Invalid Password")
	}
	token, err := createToken(auth.Login)
	if err != nil {
		return "", err
	}
	tokenstr, err := tokenToString(token)
	if err != nil {
		return "", err
	}
	printDatabase()
	return tokenstr, nil
}

func addToDatabase(auth AuthJSON) error {
	_, ok := loginMAP[auth.Login]
	if ok {
		return OnlyServerError(fmt.Sprintf("User %s already exists", auth.Login))
	}
	loginMAP[auth.Login] = UserLOGININFO{password: auth.Mdp, currentSession: -1}
	profilMap[auth.Login] = UserProfile{Connection: 0}
	printDatabase()
	return nil
}

func deleteFromDatabase(username string) error {
	_, ok := loginMAP[username]
	if !ok {
		return OnlyServerError(fmt.Sprintf("User %s dont exists", username))
	}
	delete(loginMAP, username)
	delete(profilMap, username)
	printDatabase()
	return nil
}

func printDatabase() {
	log.Println("DatabaseState : ")
	for name, data := range loginMAP {
		log.Printf("User : %s | Data : %v\n", name, data)
	}
}

func verifySession(tokenString string) (string, error) {
	var claims SimpleClaims
	token, err := jwt.ParseWithClaims(tokenString, &SimpleClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid Token at parsing claims")
		}
		return serverKey, nil
	})

	if err != nil {
		return "", ServerRuntimeError("Invalid Token", err)
	}

	claims = *token.Claims.(*SimpleClaims)
	if err := claims.Valid(); err != nil {
		return "", err
	}
	return claims.Username, nil
}

func getUserData(username string) (UserProfile, error) {
	profile := profilMap[username]
	profile.Connection++
	profilMap[username] = profile
	return profile, nil
}

func logoutAccount(username string) error {
	user := loginMAP[username]
	user.currentSession = -1
	loginMAP[username] = user
	return nil
}

/*
func fn() {

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
} */
