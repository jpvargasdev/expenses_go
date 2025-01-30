package controller

import (
	"guilliman/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetCategoriesController(c *gin.Context) {
	categories, err := models.GetCategories() // Fetch categories from storage
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *Controller) CreateCategoryController(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := models.AddCategory(newCategory) // Add category to storage
	if err != nil {
		log.Printf("Error adding category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *Controller) UpdateCategoryController(c *gin.Context) {
  var updatedCategory models.Category
  if err := c.ShouldBindJSON(&updatedCategory); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
	category, err := models.UpdateCategory(updatedCategory)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed updating category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *Controller) DeleteCategoryController(c *gin.Context) {
	var deletedCategory models.Category
	if err := c.ShouldBindJSON(&deletedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DeleteCategory(deletedCategory); err != nil {
		log.Printf("Error deleting category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed deleting category"})
		return
	}

	c.JSON(http.StatusOK, "Category deleted succesfully")
}
