package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	/* _ "github.com/mattn/go-sqlite3" */)

/* KeyUserName */
type UserLOGININFO struct {
	password  string
	ExpiresAt time.Time
}

/* Table of Users  -> UserID : PK(int),Username : string , Password : string , Session : Blob , LastSeen : string */
/* Table of ShitPost -> ShitPostID : PK(int) | Poster : FK(UserID) | Caption : string | URL : string | Date : string*/
/* Table of Msg -> MsgID : PK(int) | Sender :  FK(UserID) | Content : string | Date : string*/
/* Table of DM -> Receiver : FK(UserID) | Message : FK(MsgID)
/* Table of Comments -> Post : FK(ShitPost) | Message : FK(MsgID)  */

/* Get User Profile -> Select *(!Password) from Users */
/* Get User ID ->  Select UserID from Users where UserID = $1
/* Get User ShitPosts ->  Select * from ShitPost where Poster = Get User ID */
/* Get User Comments -> Select * from Comments where Poster = Get User ID 

/* KeyUserName */
type UserProfile struct {
	Connection int
	LastSeen   time.Time
	NBMessages int
}

type SimpleClaims struct {
	Username  string
	ExpiresAt string
}

func (c SimpleClaims) Exp() time.Time {
	exp, err := time.Parse(time.ANSIC, c.ExpiresAt)
	if err != nil {
		ServerRuntimeError("Error while parsing time", err)
	}
	return exp
}

func (c SimpleClaims) Valid() error {
	if c.Exp().IsZero() {
		OnlyServerError("Token don't have an expiration date")
	}
	if c.Exp().Before(time.Now()) {
		OnlyServerError("Token is expired")
	}
	user, ok := loginMAP[c.Username]
	if !ok {
		OnlyServerError(fmt.Sprintf("Invalid Username %s", c.Username))
	}
	if user.ExpiresAt.Format(time.ANSIC) != c.ExpiresAt {
		OnlyServerError(fmt.Sprintf("The Session Used is Invalid %s", c.Username))
	}
	_, ok = profilMap[c.Username]
	if !ok {
		OnlyServerError(fmt.Sprintf("User without profile %s", c.Username))
	}
	return nil
}

var loginMAP = map[string]UserLOGININFO{}
var profilMap = map[string]UserProfile{}

var serverKey = []byte("This is a fun serverkey")

func tokenToString(token *jwt.Token) string {
	tokenString, err := token.SignedString(serverKey)
	if err != nil {
		ServerRuntimeError("Error While Converting JWT Token to string", err)
	}
	return tokenString
}

func tokenFromString(tokenString string) *jwt.Token {
	token, err := jwt.ParseWithClaims(tokenString, &SimpleClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			OnlyServerError("Unexpected signing method")
			return nil, nil
		}
		return serverKey, nil
	})
	if err != nil {
		ServerRuntimeError("Error While Parsing JWT Token from string", err)
	}
	return token
}

func claimsFromString(tokenString string) SimpleClaims {
	token := tokenFromString(tokenString)
	claims, ok := token.Claims.(*SimpleClaims)
	if !ok {
		OnlyServerError("Error While Parsing JWT Token from string")
	}
	return *claims
}

const sessionDuration = 1 * time.Hour

func createToken(username string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, SimpleClaims{})
	expirationTime := time.Now().Add(sessionDuration)
	expirationTimeStr := expirationTime.Format(time.ANSIC)
	claims := SimpleClaims{ExpiresAt: expirationTimeStr, Username: username}
	token.Claims = claims
	user := loginMAP[username]
	user.ExpiresAt = expirationTime
	loginMAP[username] = user

	printDatabase()
	token.Claims.Valid()
	return token
}

func loginAccount(auth AuthJSON) string {
	if loginMAP[auth.Login].password != auth.Mdp {
		OnlyServerError("Invalid Password")
	}
	token := createToken(auth.Login)
	tokenstr := tokenToString(token)
	return tokenstr
}

func addToDatabase(auth AuthJSON) {
	_, ok := loginMAP[auth.Login]
	if ok {
		OnlyServerError(fmt.Sprintf("User %s already exists", auth.Login))
	}
	loginMAP[auth.Login] = UserLOGININFO{password: auth.Mdp}
	profilMap[auth.Login] = UserProfile{Connection: 0}
	printDatabase()
}

func deleteFromDatabase(username string) {
	_, ok := loginMAP[username]
	if !ok {
		OnlyServerError(fmt.Sprintf("User %s dont exists", username))
	}
	delete(loginMAP, username)
	delete(profilMap, username)
	printDatabase()
}

func printDatabase() {
	log.Println("DatabaseState : ")
	for name, data := range loginMAP {
		log.Printf("User : %s | Data : %v\n", name, data)
	}
}

func verifySession(tokenString string) string {
	claims := claimsFromString(tokenString)
	claims.Valid()
	return claims.Username
}

func getUserData(username string) UserProfile {
	profile := profilMap[username]
	profile.Connection++
	profilMap[username] = profile
	return profile
}

func logoutAccount(username string) {
	user := loginMAP[username]
	user.ExpiresAt = time.Now()
	loginMAP[username] = user

}

func isUserConnected(username string) bool {
	user := loginMAP[username]
	return user.ExpiresAt.After(time.Now())
}

func extendSession(username string) string {
	log.Println("Extending Session for user : ", username)
	token := createToken(username)
	return tokenToString(token)
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
