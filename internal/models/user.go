package models

import "golang.org/x/crypto/bcrypt"

type User struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

type UserLogged struct {
  Id int       `json:"id"`
  Password string `json:"password"`
}

func RegisterUser(user User) (User, error) {
  hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
  _, err := db.Exec(`INSERT INTO users (
    email,
    password
  ) VALUES (?, ?)`,
    user.Email,
    hashedPassword,
  )
  if err != nil {
    return User{}, err
  }

  return user, nil
}

func LoginUser(user User) (UserLogged, error) {
  var userLogged UserLogged
  err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&userLogged.Id, &userLogged.Password)
  if err != nil {
    return UserLogged{}, err
  }

  error := bcrypt.CompareHashAndPassword([]byte(userLogged.Password), []byte(user.Password)) 
  
  if error != nil {
    return UserLogged{}, error
  }

  return userLogged, nil
}

