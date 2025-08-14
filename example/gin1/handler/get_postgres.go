package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/example/gin1/entity"
)

// HandlePostgres demonstrates PostgreSQL database operations including table creation,
// user creation, retrieval, and cleanup with logging.
func (h *HTTPHandler) HandlePostgres(c *gin.Context) {
	// Create the user table in PostgreSQL database
	err := h.postgresDB.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}

	// Create a new user record
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	_ = h.postgresDB.Create(user1)

	// Retrieve a user record from database
	user2 := &entity.User{}
	h.postgresDB.Take(user2)

	// Clean up by dropping the table
	err = h.postgresDB.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}

	// Log the retrieved user information
	h.Logger.Infof("name = %s, email = %s", user2.Name, user2.Email)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
