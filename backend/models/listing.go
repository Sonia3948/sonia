
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Medication represents a medication in a listing
type Medication struct {
	ID          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Dosage      string `bson:"dosage" json:"dosage"`
	Form        string `bson:"form" json:"form"` // pill, liquid, etc
	Quantity    int    `bson:"quantity" json:"quantity"`
	Price       string `bson:"price" json:"price"`
	Availability bool   `bson:"availability" json:"availability"`
}

// Listing represents a supplier's medicine listing
type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SupplierID  primitive.ObjectID `bson:"supplierID" json:"supplierID"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Medications []Medication       `bson:"medications" json:"medications"`
	PdfURL      string             `bson:"pdfUrl" json:"pdfUrl"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
