package models

import (
	"log"
  "fmt"
  "database/sql"
)

type Category struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

func GetCategories() ([]Category, error) {
  rows, err := db.Query("SELECT id, name FROM categories")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var categories []Category
  for rows.Next() {
    var category Category
    if err := rows.Scan(&category.ID, &category.Name); err != nil {
      return nil, err
    }
    categories = append(categories, category)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return categories, nil
}

func AddCategory(category Category) (Category, error) {
  result, err := db.Exec("INSERT INTO categories (name) VALUES (?)", category.Name)
  if err != nil {
    return Category{}, err
  }

  lastID, err := result.LastInsertId()
  if err != nil {
    log.Println("Warning: Could not retrieve last insert ID for category")
  } else {
    category.ID = int(lastID)
  }

  return category, nil
}

func GetMainCategory(id int) (string, error) {
  var mainCategory string
  err := db.QueryRow("SELECT main_category FROM categories WHERE id = ?", id).Scan(&mainCategory)
  if err != nil {
    if err == sql.ErrNoRows {
      return "", fmt.Errorf("subcategory '%s' not found in categories table", string(id))
    }
    return "", err
  }
  return mainCategory, nil
}

func GetSubCategory(id int) (string, error) {
  var subCategory string
  err := db.QueryRow("SELECT name FROM categories WHERE id = ?", id).Scan(&subCategory)
  if err != nil {
    if err == sql.ErrNoRows {
      return "", fmt.Errorf("subcategory '%s' not found in categories table", string(id))
    }
    return "", err
  }
  return subCategory, nil
}
