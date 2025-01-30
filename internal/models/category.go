package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	MainCategory string `json:"main_category"`
}

func GetCategories(uid string) ([]Category, error) {
	rows, err := db.Query("SELECT id, name, main_category FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.MainCategory); err != nil {
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
	result, err := db.Exec("INSERT INTO categories (name, main_category) VALUES (?, ?)", category.Name, category.MainCategory)
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

func UpdateCategory(category Category) (Category, error) {
	_, err := db.Exec(
		"Update categories SET name = ?, main_category = ? WHERE id = ?",
		category.Name,
		category.MainCategory,
		category.ID,
	)

	if err != nil {
		log.Println("Warning: Could not retrieve last insert ID for category")
	}

	return category, nil
}

func DeleteCategory(category Category) error {
	_, err := db.Exec(
		"DELETE FROM categories WHERE id = ?",
		category.ID,
	)

	if err != nil {
		log.Println("Warning: Could not retrieve last insert ID for category")
	}

	return nil
}

func GetMainCategory(id int) (string, error) {
	var mainCategory string
	err := db.QueryRow("SELECT main_category FROM categories WHERE id = ?", id).Scan(&mainCategory)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("subcategory '%s' not found in categories table", fmt.Sprint(rune(id)))
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
			return "", fmt.Errorf("subcategory '%s' not found in categories table", fmt.Sprint(rune(id)))
		}
		return "", err
	}
	return subCategory, nil
}
