# gotest

## Tests

Quelques tests sont disponibles :

```
$ go test ./...
```

## Documentation

Il faut au préalable lancer le serveur HTTP de l'API :

```
$ go run main.go
```

La documentation de l'API est désormais accessible ici [http://localhost:8001/doc](http://localhost:8001/doc)

## API

Une fois le serveur lancé via :

```
$ go run main.go
```

L'API est disponible sur le port 8001, exemple [http://localhost:8001/peoples](http://localhost:8001/peoples)

## Choix

- J'ai voulu normaliser les données provenant de la base pour avoir un rendu propre et homogène
- Pour ne pas surcharger le test il n'y a pas la possibilité d'attacher ou détacher un véhicule ou vaisseau d'un personnage, mais je peux le rajouter si nécessaire