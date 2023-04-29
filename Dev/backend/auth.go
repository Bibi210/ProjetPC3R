package main

import (
	"Backend/Database"
	"Backend/Helpers"
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const sessionDuration = 1 * time.Hour

type simpleClaims struct {
	Username  string
	ExpiresAt string
}

func (c simpleClaims) Exp() time.Time {
	exp, err := time.Parse(time.ANSIC, c.ExpiresAt)
	Helpers.ServerRuntimeError("Error while parsing time", err)
	return exp
}

func (c simpleClaims) Valid() error {
	db := Database.OpenDatabase()
	defer Database.CleanCloser(db)
	if c.Exp().IsZero() {
		Helpers.OnlyServerError("Token don't have an expiration date")
	}
	if c.Exp().Before(time.Now()) {
		Helpers.OnlyServerError("Token is expired")
	}
	user := Database.GetUser(db, c.Username)
	if Helpers.FormatTime(user.Session) != c.ExpiresAt {
		Helpers.OnlyServerError("The Session Used is Invalid : " + c.Username)
	}
	return nil
}

var serverKey = []byte("This is a fun serverKey")

func tokenToString(token *jwt.Token) tokenString {
	out, err := token.SignedString(serverKey)
	Helpers.ServerRuntimeError("Error While Converting JWT Token to string", err)
	return tokenString(out)
}

func tokenFromString(tokenString tokenString) *jwt.Token {
	token, err := jwt.ParseWithClaims(string(tokenString), &simpleClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			Helpers.OnlyServerError("Unexpected signing method")
			return nil, nil
		}
		return serverKey, nil
	})
	Helpers.ServerRuntimeError("Error While Parsing JWT Token from string", err)
	return token
}

func claimsFromString(tokenString tokenString) simpleClaims {
	token := tokenFromString(tokenString)
	claims, ok := token.Claims.(*simpleClaims)
	if !ok {
		Helpers.OnlyServerError("Error While Parsing Claims from string")
	}
	return *claims
}

func createToken(db *sql.DB, name string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, simpleClaims{})
	expirationTime := time.Now().Add(sessionDuration)
	expirationTimeStr := Helpers.FormatTime(expirationTime)
	claims := simpleClaims{ExpiresAt: expirationTimeStr, Username: name}
	token.Claims = claims
	user := Database.GetUser(db, name)
	user.Session = expirationTime
	user.LastSeen = time.Now()
	user.UpdateUserSession(db)
	token.Claims.Valid()
	return token
}

func verifySession(tokenString tokenString) username {
	if tokenString == "" {
		Helpers.OnlyServerError("User is not logged in")
	}
	claims := claimsFromString(tokenString)
	claims.Valid()
	return username(claims.Username)
}

func isLogged(db *sql.DB, username username) bool {
	if !Database.IsUserExist(db, string(username)) {
		return false
	}
	user := Database.GetUser(db, string(username))
	return user.Session.After(time.Now())
}

func extendSession(db *sql.DB, username string) tokenString {
	token := createToken(db, username)
	return tokenToString(token)
}

func loginAccount(db *sql.DB, auth Helpers.RequestAuthJSON) tokenString {
	user := Database.GetUser(db, auth.Login)
	if user.Password != auth.Mdp {
		Helpers.OnlyServerError("Invalid Password")
	}
	token := createToken(db, auth.Login)
	tokenStr := tokenToString(token)
	return tokenStr
}

func logoutAccount(db *sql.DB, username username) {
	user := Database.GetUser(db, string(username))
	user.Session = time.Time{}
	user.UpdateUserSession(db)
}
