package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const sessionDuration = 1 * time.Hour

type simple_claims struct {
	username  string
	expiresAt string
}

func parseTime(timeStr string) time.Time {
	time, err := time.Parse(time.ANSIC, timeStr)
	ServerRuntimeError("Error while parsing time", err)
	return time
}

func formatTime(t time.Time) string {
	return t.Format(time.ANSIC)
}

func (c simple_claims) Exp() time.Time {
	exp, err := time.Parse(time.ANSIC, c.expiresAt)
	ServerRuntimeError("Error while parsing time", err)
	return exp
}

func (c simple_claims) Valid() error {
	db := openDatabase()
	defer closeDatabase(db)
	if c.Exp().IsZero() {
		OnlyServerError("Token don't have an expiration date")
	}
	if c.Exp().Before(time.Now()) {
		OnlyServerError("Token is expired")
	}
	user := getUser(db, username(c.username))
	if formatTime(user.session) != c.expiresAt {
		OnlyServerError(fmt.Sprintf("The Session Used is Invalid %s", c.username))
	}
	return nil
}

var serverKey = []byte("This is a fun serverkey")

func tokenToString(token *jwt.Token) token_string {
	out, err := token.SignedString(serverKey)
	ServerRuntimeError("Error While Converting JWT Token to string", err)
	return token_string(out)
}

func tokenFromString(tokenString token_string) *jwt.Token {
	token, err := jwt.ParseWithClaims(string(tokenString), &simple_claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			OnlyServerError("Unexpected signing method")
			return nil, nil
		}
		return serverKey, nil
	})
	ServerRuntimeError("Error While Parsing JWT Token from string", err)
	return token
}

func claimsFromString(tokenString token_string) simple_claims {
	token := tokenFromString(tokenString)
	claims, ok := token.Claims.(*simple_claims)
	if !ok {
		OnlyServerError("Error While Parsing Claims from string")
	}
	return *claims
}

func createToken(db *sql.DB, name string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, simple_claims{})
	expirationTime := time.Now().Add(sessionDuration)
	expirationTimeStr := expirationTime.Format(time.ANSIC)
	claims := simple_claims{expiresAt: expirationTimeStr, username: name}
	token.Claims = claims
	user := getUser(db, username(name))
	user.session = expirationTime
	user.Update(db)
	token.Claims.Valid()
	return token
}

func verifySession(tokenString token_string) username {
	if tokenString == "" {
		OnlyServerError("User is not logged in")
	}
	claims := claimsFromString(tokenString)
	claims.Valid()
	return username(claims.username)
}

func isLogged(db *sql.DB, username username) bool {
	user := getUser(db, username)
	return user.session.After(time.Now())
}

func extendSession(db *sql.DB, username string) token_string {
	log.Println("Extending Session for user : ", username)
	token := createToken(db, username)
	return tokenToString(token)
}

func loginAccount(db *sql.DB, auth AuthJSON) token_string {
	user := getUser(db, username(auth.Login))
	if user.password != auth.Mdp {
		OnlyServerError("Invalid Password")
	}
	token := createToken(db, auth.Login)
	tokenstr := tokenToString(token)
	return tokenstr
}

func logoutAccount(db *sql.DB, username username) {
	user := getUser(db, username)
	user.session = time.Now()
	user.UpdateSession(db)
}
