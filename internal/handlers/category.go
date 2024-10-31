package handlers

import (
  "guilliman/internal/models"
  "github.com/gin-gonic/gin"
  "net/http"
)

func GetCategoriesHandler(c *gin.Context) {
  categories, err := models.GetCategories() // Fetch categories from storage
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, categories)
}

func CreateCategoryHandler(c *gin.Context) {
  var newCategory models.Category
  if err := c.ShouldBindJSON(&newCategory); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  models.AddCategory(newCategory) // Add category to storage
  c.JSON(http.StatusCreated, newCategory)
}
