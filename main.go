package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// SwimLane ...
type SwimLane struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Todo ...
type Todo struct {
	gorm.Model
	ID     int       `json:"id" gorm:"autoIncrement"`
	UUID   uuid.UUID `json:"uuid" gorm:"primaryKey"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
	//`json:"id" gorm:"unique;primaryKey;autoIncrement"`
}

func main() {
	initDB()

	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.Static("/static", "./static")

	e.GET("/", home)
	e.POST("/todos", postTodo)
	e.POST("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)

	e.Run(":8000")

}

func home(c *gin.Context) {
	var todos []Todo
	db.Find(&todos)

	var lanes = []SwimLane{
		{ID: "new", Name: "Todo"},
		{ID: "wip", Name: "Progress"},
		{ID: "done", Name: "Completed"},
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos": todos,
		"lanes": lanes,
	})

	// c.JSON(http.StatusOK, gin.H{"Message": "updated"})
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("DB.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Todo{})
}

func postTodo(c *gin.Context) {
	title := c.PostForm("title")
	status := c.PostForm("status")
	if status == "" {
		status = "new"
	}
	newTodo := Todo{Title: title, Status: status, UUID: uuid.New()}
	db.Create(&newTodo)
	fmt.Println(newTodo.ID, newTodo.UUID)

	c.HTML(http.StatusOK, "task.html", gin.H{
		"Title":  title,
		"Status": status,
		"Id":     newTodo.ID,
	})
}

func updateTodo(c *gin.Context) {
	uuid := c.PostForm("uuid")
	status := c.PostForm("status")

	var todo Todo
	// db.First(&todo, 1) // find with ID 1
	db.First(&todo, "UUID = ?", uuid) // find with UUID=uuid
	// Update - update todo's status
	db.Model(&todo).Update("Status", status)
	// multiple fields
	// db.Model(&todo).Updates(Todo{Title: "new title", Status: "new"})
	c.JSON(http.StatusOK, gin.H{"Message": "updated"})
}

func deleteTodo(c *gin.Context) {
	uuid := c.Param("uuid")
	var todo Todo
	db.Delete(&todo, "UUID=?", uuid)
}
