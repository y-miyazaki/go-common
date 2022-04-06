package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/example/gin1/entity"
)

// GetPostgres handler
func (h *HTTPHandler) GetPostgres(c *gin.Context) {
	err := h.postgresDB.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	_ = h.postgresDB.Create(user1)

	user2 := &entity.User{}
	h.postgresDB.Take(user2)
	err = h.postgresDB.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}

	h.Logger.Infof("name = %s, email = %s", user2.Name, user2.Email)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
