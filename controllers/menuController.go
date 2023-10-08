package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lord-tx/golang-restaurant-management/database"
	"github.com/lord-tx/golang-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := menuCollection.Find(context, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching menu items"})
			return
		}

		var allMenus []bson.M
		if err = result.All(context, &allMenus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}

		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := foodCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(menu)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		menu.Created_at = time.Now()
		menu.Updated_at = time.Now()
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_date != nil && menu.End_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "kindly retype time"})
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_date})
		}

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Name})
		}
		menu.Updated_at = time.Now()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, updateErr := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		}

		c.JSON(http.StatusOK, result)

	}
}
