package handlers

import (
	"fmt"
	"hms/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

// ! Gettin all the entires that are present in  your inventory
func GetAllInventory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var inventory []models.Inventory
		if err := db.Find(&inventory).Error; err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, "Error while fetching the data")
		}

		html := ""

		for _, invrow := range inventory {
			html += fmt.Sprintf(
				`<tr class="hover:bg-gray-900 hover:text-white">
          <th class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%s</th>
          <th class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%s</th>
          <th class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%d</th>
          <th class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%.2f</th>
        </tr>`, invrow.ArticalId, invrow.Name, invrow.Quantity, invrow.Price,
			)
		}
		return c.HTML(http.StatusOK, html)
	}
}

// ! Get the details of the item by specific article id
func GetInventoryById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		//? extracting the parameter from the url of the request
		articalID := c.QueryParam("artical_id")
		var inventory models.Inventory

		if err := db.First(&inventory, "artical_id = ?", articalID).Error; err != nil {
			
			log.Error(err)
			fmt.Println("no match id")

			errorMessage := "No results found for the given Article ID"
			tableHTML := fmt.Sprintf(`
				<tr>
					<td colspan="4" class="text-center text-red-500 font-bold py-4">
							%s
					</td>
				</tr>`, errorMessage)

			return c.HTML(http.StatusOK, tableHTML)

		}

		//log.Printf("Fetched Inventory: %+v", inventory)
		//return c.String(http.StatusOK, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%f</td></tr>", inventory.ArticalId, inventory.Name, inventory.Quantity, inventory.Price))
		return c.HTML(http.StatusOK, fmt.Sprintf(`
    <tr class="border border-black hover:bg-gray-100">
        <td class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%s</td>
        <td class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%s</td>
        <td class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%d</td>
        <td class="border border-gray-900 px-4 py-2 text-left cursor-pointer">%.2f</td>
    </tr>
		`, inventory.ArticalId, inventory.Name, inventory.Quantity, inventory.Price))

		//return c.Render(http.StatusOK,"searchingrow.html",inventory)

		/*
			return c.JSON(http.StatusOK, echo.Map{
				"Item Name": inventory.Name,
				"Quantity":  inventory.Quantity,
			})
		*/
	}
}

//! Create the record of the new article in the inventory

func CreateNewInventory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		//getting the values from the form and coverting it to required datatype
		inventory := models.Inventory{
			ArticalId: c.FormValue("artical_id"),
			Name:      c.FormValue("name"),
			Quantity:  parseInt(c.FormValue("quantity")),
			Price:     parseFloat(c.FormValue("price")),
		}

		if inventory.ArticalId == "" || inventory.Name == "" || inventory.Quantity <= 0 || inventory.Price <= 0 {

			return c.String(http.StatusBadRequest, "invalid input data")
		}

		if err := db.Create(&inventory).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Failed to create new inventory item")
		}

		return c.Render(http.StatusCreated, "added.html", inventory)

		//todo: this works when we are testing or sending the data from POSTMAN as a JSON
		/*
			//? Binding the incoming request in the form of json
			if err := c.Bind(&inventory); err != nil {
				log.Error(err)
				return c.JSON(http.StatusBadRequest, "Invalid Input data for the article")
			}

			if err := db.Create(&inventory).Error; err != nil {
				log.Error(err)
				return c.JSON(http.StatusInternalServerError, "Can not put the article in the inventory")
			}

			return c.JSON(http.StatusOK, inventory)
		*/

	}
}

//? hepler functions to convert the form elements to numbers

func parseInt(value string) int32 {
	v, _ := strconv.Atoi(value)
	return int32(v)
}

func parseFloat(value string) float64 {
	v, _ := strconv.ParseFloat(value, 64)
	return v
}

//! Updating the inventory of already existing the articals

func UpdateInventoryById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.Req

		//? Binding the incoming request in the form of json
		if err := c.Bind(&req); err != nil {
			c.Logger().Error(err, "Error in processing the incoming request")

			return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Invalid Input"})
		}

		//? to fetch the entire structure of the inventory
		var inventory models.Inventory

		//? to fetch using specific artical id
		if err := db.Where("artical_id = ?", req.ArticalId).First(&inventory).Error; err != nil {
			c.Logger().Error(err, "Error in matching the article id")

			return c.JSON(http.StatusNotFound, echo.Map{"Error": "No such Article found"})
		}

		switch req.Operation {
		case "In":
			inventory.Quantity += int32(req.Quantity)
		case "Out":
			if inventory.Quantity < int32(req.Quantity) {
				return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Insufficient number of items"})
			}
			inventory.Quantity -= int32(req.Quantity)
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Please use valid Operations"})
		}

		if err := db.Save(&inventory).Error; err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"Error": "There is error while saving the database"})
		}

		return c.JSON(http.StatusOK, echo.Map{"Updated Article": inventory})
	}
}

//! delete the already existing item in the inventory

func DeleteFromInventoryById(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		articleId := c.Param("artical_id")
		if err := db.Delete(&models.Inventory{}, "artical_id = ?", articleId).Error; err != nil {
			c.Logger().Error(err, "Error in deleting the article id")

			return c.JSON(http.StatusNotFound, echo.Map{"Error": "No such artile id found"})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"Succsess": "Operation Completed",
		})
	}
}
