package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// Product - struct sản phẩm
type Product struct {
	Name     string
	Quantity int
	Price    float64
}

// Invoice - struct hóa đơn
type Invoice struct {
	Products    []Product
	TotalAmount float64
}

// Hàm generatePDF để tạo PDF từ dữ liệu hóa đơn
func generatePDF(invoice Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Product Name")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(30, 10, "Price")
	pdf.Cell(30, 10, "Total")
	pdf.Ln(10)
	for _, product := range invoice.Products {
		pdf.Cell(40, 10, product.Name)
		pdf.Cell(30, 10, strconv.Itoa(product.Quantity))
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", product.Price))
		total := product.Price * float64(product.Quantity)
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", total))
		pdf.Ln(10)
	}

	pdf.Ln(10)
	pdf.Cell(40, 10, "Total Amount")
	pdf.Cell(30, 10, fmt.Sprintf("%.2f", invoice.TotalAmount))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("./templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/create-pdf", func(c *gin.Context) {
		var products []Product
		for i := 0; ; i++ {
			name := c.PostForm(fmt.Sprintf("products[%d].name", i))
			if name == "" {
				break
			}
			quantity, _ := strconv.Atoi(c.PostForm(fmt.Sprintf("products[%d].quantity", i)))
			price, _ := strconv.ParseFloat(c.PostForm(fmt.Sprintf("products[%d].price", i)), 64)
			products = append(products, Product{Name: name, Quantity: quantity, Price: price})
		}
		totalAmount, _ := strconv.ParseFloat(c.PostForm("total_amount"), 64)

		invoice := Invoice{
			Products:    products,
			TotalAmount: totalAmount,
		}

		log.Printf("Invoice: %+v", invoice)
		pdfBytes, err := generatePDF(invoice)
		if err != nil {
			log.Printf("Error generating PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	})

	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
