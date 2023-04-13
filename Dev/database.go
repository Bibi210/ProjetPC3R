package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
)

/* KeyUserName */
type UserLOGININFO struct {
	password       string
	currentSession string
}

/* KeyUserName */
type UserProfile struct {
	connection int
}

var loginMAP = map[string]UserLOGININFO{}
var profilMap = map[string]UserProfile{}

var serverKey = []byte("This is a fun serverkey")

func loginAccount(auth Auth) (string, error) {
	user, ok := loginMAP[auth.Login]
	if !ok {
		return "", OnlyServerError(fmt.Sprintf("Invalid Username %s", auth.Login))
	}
	if user.password != auth.Mdp {
		return "", OnlyServerError("Invalid Password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	claims["username"] = auth.Login
	tokenString, err := token.SignedString(serverKey)
	if err != nil {
		return "", ServerRuntimeError("Error While Generating JWT Token", err)
	}

	user.currentSession = tokenString
	loginMAP[auth.Login] = user
	printDatabase()
	return tokenString, nil
}

func addToDatabase(auth Auth) error {
	_, ok := loginMAP[auth.Login]
	if ok {
		return OnlyServerError("User already registered")
	}
	loginMAP[auth.Login] = UserLOGININFO{password: auth.Mdp, currentSession: ""}
	profilMap[auth.Login] = UserProfile{connection: 0}
	printDatabase()
	return nil
}

func printDatabase() {
	log.Println("DatabaseState : ")
	for name, data := range loginMAP {
		log.Printf("User : %s | Data : %s\n", name, data)
	}
}

func getUserData(tokenString string) (UserProfile, error) {
	var defaultUser = UserProfile{connection: 0}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, OnlyServerError("Invalid Token")
		}
		return token, nil
	})

	if err != nil {
		return defaultUser, OnlyServerError("Invalid Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return defaultUser, OnlyServerError("Can't Parse JWT Token")
	}
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return defaultUser, OnlyServerError("Token Expired")
	}

	username := claims["username"].(string)
	user, ok := loginMAP[username]
	if !ok {
		return defaultUser, OnlyServerError(fmt.Sprintf("Invalid Username %s", username))
	}
	if user.currentSession != tokenString {
		return defaultUser, OnlyServerError(fmt.Sprintf("The Session Used is Invalid %s", username))
	}
	profile, ok := profilMap[username]
	if !ok {
		return defaultUser, OnlyServerError(fmt.Sprintf("User without profile %s", username))
	}
	profile.connection++
	return profile, nil
}

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
}
