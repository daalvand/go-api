package main

import (
	"github.com/daalvand/go-api/db"
	"github.com/daalvand/go-api/routes"
	"github.com/daalvand/go-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	storage.InitLogger()

	database, err := initializeDatabase()
	if err != nil {
		logrus.Fatal("Error initializing database", err)
	}

	migrateTables(database)

	server := gin.Default()
	routes.SetupRoutes(server, database)

	err = server.Run(":8080")
	if err != nil {
		logrus.Fatal("Server Error", err)
	}
}

func migrateTables(db *db.Database) {
	if err := db.MigrateTables(); err != nil {
		logrus.Fatal("Error migrating tables", err)
	}
}

func initializeDatabase() (*db.Database, error) {
	dialector := sqlite.Open("database.sqlite")
	config := &gorm.Config{}
	return db.InitializeDB(dialector, config)
}
