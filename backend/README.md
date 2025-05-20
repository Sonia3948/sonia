
# El Saidaliya Backend API

Ce dossier contient l'API backend pour la plateforme El Saidaliya, écrite en Go avec Gin et MongoDB.

## Prérequis

- Go (version 1.16 ou supérieure)
- MongoDB (installé localement ou accessible via URI)

## Configuration

1. Copiez le fichier `.env.example` en `.env` et modifiez les variables selon votre environnement.
2. Assurez-vous que MongoDB est en cours d'exécution.

## Installation

```bash
# Installer les dépendances
./setup.sh
# ou manuellement
go mod tidy
```

## Démarrage du serveur

```bash
go run main.go
```

Le serveur démarrera par défaut sur le port 8080. Vous pouvez vérifier que tout fonctionne en accédant à `http://localhost:8080/api/health`.

## Structure du projet

- `main.go` - Point d'entrée de l'application et configuration du serveur
- `models/` - Définitions des structures de données (à venir)
- `controllers/` - Logique de traitement des requêtes (à venir)
- `middlewares/` - Middlewares pour l'authentification, etc. (à venir)
- `routes/` - Définitions des routes API (à venir)
- `utils/` - Fonctions utilitaires (à venir)

## API Endpoints

### Santé
- `GET /api/health` - Vérifier l'état de l'API

### Authentification
- `POST /api/auth/register` - Enregistrer un nouvel utilisateur
- `POST /api/auth/login` - Connecter un utilisateur
- `POST /api/auth/forgot-password` - Réinitialiser le mot de passe

### Utilisateurs
- `GET /api/users` - Récupérer tous les utilisateurs
- `GET /api/users/:id` - Récupérer un utilisateur par ID
- `POST /api/users` - Créer un utilisateur
- `PUT /api/users/:id` - Mettre à jour un utilisateur
- `DELETE /api/users/:id` - Supprimer un utilisateur

### Listings (PDF)
- `GET /api/listings` - Récupérer tous les listings
- `GET /api/listings/:id` - Récupérer un listing par ID
- `POST /api/listings` - Créer un listing
- `PUT /api/listings/:id` - Mettre à jour un listing
- `DELETE /api/listings/:id` - Supprimer un listing
- `GET /api/listings/search` - Rechercher des listings

### Offres
- `GET /api/offers` - Récupérer toutes les offres
- `GET /api/offers/:id` - Récupérer une offre par ID
- `POST /api/offers` - Créer une offre
- `PUT /api/offers/:id` - Mettre à jour une offre
- `DELETE /api/offers/:id` - Supprimer une offre
