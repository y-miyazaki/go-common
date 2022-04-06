package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/example/gin1/entity"
)

// GetMySQL handler
func (h *HTTPHandler) GetMySQL(c *gin.Context) {
	err := h.mysqlDB.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	h.mysqlDB.Create(user1)

	user2 := &entity.User{}
	h.mysqlDB.Take(user2)
	err = h.mysqlDB.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
