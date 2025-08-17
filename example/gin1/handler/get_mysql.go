package handler

import (
	"net/http"

	"go-common/example/gin1/entity"

	"github.com/gin-gonic/gin"
)

// HandleMySQL demonstrates MySQL database operations including table creation,
// user creation, retrieval, and cleanup.
func (h *HTTPHandler) HandleMySQL(c *gin.Context) {
	// Create the user table in MySQL database
	err := h.mysqlDB.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}

	// Create a new user record
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	h.mysqlDB.Create(user1)

	// Retrieve a user record from database
	user2 := &entity.User{}
	h.mysqlDB.Take(user2)

	// Clean up by dropping the table
	err = h.mysqlDB.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
