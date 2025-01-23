package controller

import (
	"guilliman/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetCategoriesController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }


	categories, err := models.GetCategories(uid) // Fetch categories from storage
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *Controller) CreateCategoryController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }

	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  newCategory.UserID = uid 

	category, err := models.AddCategory(newCategory) // Add category to storage
	if err != nil {
		log.Printf("Error adding category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *Controller) UpdateCategoryController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }

  var updatedCategory models.Category
  if err := c.ShouldBindJSON(&updatedCategory); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  updatedCategory.UserID = uid 
  
  category, err := models.UpdateCategory(updatedCategory)
  if err != nil {
    log.Printf("Error updating category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed updating category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}


func (h *Controller) DeleteCategoryController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }

  var deletedCategory models.Category
  if err := c.ShouldBindJSON(&deletedCategory); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  deletedCategory.UserID = uid 
  
  err := models.DeleteCategory(deletedCategory)
  if err != nil {
    log.Printf("Error deleting category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed deleting category"})
		return
	}

	c.JSON(http.StatusOK, "Category deleted succesfully")
}
