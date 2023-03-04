package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Product struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       image.Image
}

const ImageSize = 64 * 64

var products = []Product{
	{Id: "meet", Name: "Meet of Cow", Description: "juicy"},
	{Id: "fish", Name: "Fish", Description: "caught in the sea"},
	{Id: "bacon", Name: "Bacon of Pig", Description: "fresh"},
}

func main() {

	router := gin.Default()

	router.GET("/list", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, products)
		for _, product := range products {
			buf := new(bytes.Buffer)
			err := jpeg.Encode(buf, product.Image, nil)
			if err != nil {
				fmt.Println("Failed to encode image")
			}
			bytes := buf.Bytes()
			ctx.Writer.Write(bytes)
		}
	})

	router.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		for _, product := range products {
			if product.Id == id {
				buf := new(bytes.Buffer)
				err := jpeg.Encode(buf, product.Image, nil)
				if err != nil {
					fmt.Println("Failed to encode image")
				}
				bytes := buf.Bytes()
				ctx.IndentedJSON(http.StatusOK, product)
				ctx.Writer.Write(bytes)
				return
			}
		}
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
	})

	router.PUT("/add", func(ctx *gin.Context) {
		var product Product

		product.Id = uuid.NewString()
		product.Name = ctx.Query("name")
		product.Description = ctx.Query("description")

		product.Image, _, _ = image.Decode(ctx.Request.Body)

		products = append(products, product)

		ctx.IndentedJSON(http.StatusOK, product)
	})

	router.POST("/update", func(ctx *gin.Context) {
		var got Product
		err := ctx.BindJSON(&got)
		fmt.Printf("got: %v", got)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "incorrect JSON"})
		}
		for i, product := range products {
			if product.Id == got.Id {
				products[i] = got
				ctx.IndentedJSON(http.StatusOK, gin.H{"message": "updated"})
				return
			}
		}
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
	})

	router.DELETE("/remove", func(ctx *gin.Context) {
		id := ctx.Query("id")
		for i, product := range products {
			if product.Id == id {
				removed := products[i]
				products = append(products[:i], products[i+1:]...)
				ctx.IndentedJSON(http.StatusOK, removed)
				return
			}
		}
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
	})

	err := router.Run()
	if err != nil {
		panic("Bad server")
	}
}
