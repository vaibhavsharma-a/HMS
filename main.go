package main

import (
	
	"fmt"
	"hms/handlers"
	"hms/models"
	"log"
	"github.com/joho/godotenv"
	"os"

	"github.com/labstack/echo/v4"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	if os.Getenv("APP_ENV") != "production"{
		if err := godotenv.Load(); err!=nil{
			log.Println(err)
		}
	}

	

	user_db := os.Getenv("postgres_user") 
	pass_db := os.Getenv("postgres_pass") 
	db_name := os.Getenv("postgres_db") 
	db_port := os.Getenv("postgres_port") 

	dns := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",user_db,pass_db,db_name,db_port)  //"user=? password=? dbname=? port=? sslmode=disable/enable"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("could not connect to the Sql DB object: %v", err)
	}
	defer sqlDB.Close()

	/*err = sqlDB.Ping()
	if err!=nil{
		log.Fatalf("Could not make the connection to the database: %v",err)
	}*/

	//!Auto migration of the inventory table
	db.AutoMigrate(&models.Inventory{})

	fmt.Println("Connection is successfull")

	//!creating a server with echo!!

	e := echo.New()

	


	//? route to load the html template
	e.GET("/", func(c echo.Context) error {
		return c.File("templates/home.html")
	})

	//todo: Passing the aritcle_id as body 
	e.GET("/inventory", handlers.GetAllInventory(db))
	e.GET("/inventory/", handlers.GetInventoryById(db))
	e.POST("/inventory/add", handlers.CreateNewInventory(db))
	e.PUT("/inventory",handlers.UpdateInventoryById(db))
	e.DELETE("/inventory/:artical_id",handlers.DeleteFromInventoryById(db)) //? passing the aricle_id as url param

	e.Logger.Fatal(e.Start(":8080"))

}
