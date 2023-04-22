package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SimpleClaims struct {
	Username  string
	ExpiresAt string
}

func parseTime(timeStr string) time.Time {
	time, err := time.Parse(time.ANSIC, timeStr)
	ServerRuntimeError("Error while parsing time", err)
	return time
}

func formatTime(t time.Time) string {
	return t.Format(time.ANSIC)
}

func (c SimpleClaims) Exp() time.Time {
	exp, err := time.Parse(time.ANSIC, c.ExpiresAt)
	ServerRuntimeError("Error while parsing time", err)
	return exp
}

func (c SimpleClaims) Valid() error {
	db := openDatabase()
	defer closeDatabase(db)
	if c.Exp().IsZero() {
		OnlyServerError("Token don't have an expiration date")
	}
	if c.Exp().Before(time.Now()) {
		OnlyServerError("Token is expired")
	}
	user := getUser(db, c.Username)
	if formatTime(user.Session) != c.ExpiresAt {
		OnlyServerError(fmt.Sprintf("The Session Used is Invalid %s", c.Username))
	}
	return nil
}

var serverKey = []byte("This is a fun serverkey")

func tokenToString(token *jwt.Token) tokenString {
	out, err := token.SignedString(serverKey)
	if err != nil {
		ServerRuntimeError("Error While Converting JWT Token to string", err)
	}
	return tokenString(out)
}

func tokenFromString(tokenString tokenString) *jwt.Token {
	token, err := jwt.ParseWithClaims(string(tokenString), &SimpleClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func claimsFromString(tokenString tokenString) SimpleClaims {
	token := tokenFromString(tokenString)
	claims, ok := token.Claims.(*SimpleClaims)
	if !ok {
		OnlyServerError("Error While Parsing JWT Token from string")
	}
	return *claims
}

const sessionDuration = 1 * time.Hour

func createToken(db *sql.DB, username string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, SimpleClaims{})
	expirationTime := time.Now().Add(sessionDuration)
	expirationTimeStr := expirationTime.Format(time.ANSIC)
	claims := SimpleClaims{ExpiresAt: expirationTimeStr, Username: username}
	token.Claims = claims
	user := getUser(db, username)
	user.Session = expirationTime
	user.UpdateRow(db)
	token.Claims.Valid()
	return token
}

func verifySession(tokenString tokenString) username {
	if tokenString == "" {
		OnlyServerError("User is not logged in")
	}
	claims := claimsFromString(tokenString)
	claims.Valid()
	return username(claims.Username)
}

func isLogged(db *sql.DB, username username) bool {
	user := getUser(db, string(username))
	return user.Session.After(time.Now())
}

func extendSession(db *sql.DB, username string) tokenString {
	log.Println("Extending Session for user : ", username)
	token := createToken(db, username)
	return tokenToString(token)
}

func loginAccount(db *sql.DB, auth AuthJSON) tokenString {
	user := getUser(db, auth.Login)
	if user.Password != auth.Mdp {
		OnlyServerError("Invalid Password")
	}
	token := createToken(db, auth.Login)
	tokenstr := tokenToString(token)
	return tokenstr
}

func logoutAccount(db *sql.DB, username username) {
	user := getUser(db, string(username))
	user.Session = time.Now()
	user.UpdateSession(db)
}
