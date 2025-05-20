package auth

import (
	"context"
	"log"
	"time"

	"elsaidaliya/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var adminCollection *mongo.Collection
var userCollection *mongo.Collection

// Init initialise les collections du package auth.
func Init(db *mongo.Database) {
	adminCollection = db.Collection("admins")
	userCollection = db.Collection("users")
	createAdminUser()
}

// createAdminUser vérifie si un admin existe et en crée un sinon.
func createAdminUser() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingAdmin models.AdminUser
	err := adminCollection.FindOne(ctx, bson.M{"phone": "0549050018"}).Decode(&existingAdmin)

	if err == mongo.ErrNoDocuments {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Ned@0820"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Erreur de hash du mot de passe admin : %v", err)
			return
		}

		admin := models.AdminUser{
			ID:        primitive.NewObjectID(),
			Phone:     "0549050018",
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}

		if _, err := adminCollection.InsertOne(ctx, admin); err != nil {
			log.Printf("Erreur lors de la création de l'utilisateur admin : %v", err)
			return
		}

		log.Println("Utilisateur admin créé avec succès.")
	} else if err != nil {
		log.Printf("Erreur lors de la vérification de l'admin : %v", err)
	} else {
		log.Println("L'utilisateur admin existe déjà.")
	}
}
