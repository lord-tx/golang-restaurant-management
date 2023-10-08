package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lord-tx/golang-restaurant-management/database"
	helper "github.com/lord-tx/golang-restaurant-management/helpers"
	"github.com/lord-tx/golang-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex, err := strconv.Atoi(c.Query("startIndex"))
		if err != nil || startIndex < 0 {
			startIndex = (page - 1) * recordPerPage
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allUsers bson.M

		err = result.All(ctx, &allUsers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User

		filter := bson.M{"user_id": userId}
		err := userCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		/// STEPS:
		// - Convert JSON from user (postman) to GoObject (Struct)
		var user models.User

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// - Validate GoStruct
		validationError := validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		// - Check if userEmail already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// - Check if phoneNumber already exists

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return Email or username matters
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone already exists"})
			return
		}

		// - Hash Password
		password := HashPassword(*user.Password)
		user.Password = &password

		// - Create Extra Details (timestamps... etc)
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// - GenerateTokens
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// - Insert into database
		resultInsertionNumber, insertionError := userCollection.InsertOne(ctx, user)
		if insertionError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// - Return httpStatus and result
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		/// STEPS:
		// - Convert JSON from user (postman) to GoObject (Struct)

		var user models.User
		var foundUser models.User

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// - Find the user with the email provided
		filter := bson.M{"user_id": user.User_id}
		err = userCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// - Verify Password
		isPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// - Generate Token

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *&foundUser.User_id)

		// - Update Tokens
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// - Return httpStatus and result
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	msg := ""
	check := false
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}

	return check, msg
}
