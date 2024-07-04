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
	Name     string  `form:"products[0].name"`
	Quantity int     `form:"products[0].quantity"`
	Price    float64 `form:"products[0].price"`
}

// Invoice - struct hóa đơn
type Invoice struct {
	Products    []Product `form:"products"`
	TotalAmount float64   `form:"total_amount"`
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

	// Load template HTML
	r.LoadHTMLFiles("./templates/index.html")

	// Define route for the root path
	r.GET("/", func(c *gin.Context) {
		// Render the index.html template when accessing the root path
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Define route for creating PDF
	r.POST("/create-pdf", func(c *gin.Context) {
		var invoice Invoice
		// Bind form data to the Invoice struct
		if err := c.ShouldBind(&invoice); err != nil {
			log.Printf("Error binding data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Generate PDF based on the invoice data
		pdfBytes, err := generatePDF(invoice)
		if err != nil {
			log.Printf("Error generating PDF: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Send the generated PDF as a response
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	})

	// Start the server on port 8080
	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
