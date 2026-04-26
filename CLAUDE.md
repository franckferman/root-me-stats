# Root-me Stats Project Documentation

Ce projet génère des badges SVG pour les profils Root-me avec support de plusieurs thèmes et comparaisons entre utilisateurs.

## Architecture du Projet

```
root-me-stats/
├── cmd/                    # Points d'entrée (binaires)
│   ├── server/            # Serveur HTTP API
│   ├── cli/               # CLI complet avec toutes les fonctions
│   └── badges/            # CLI simple pour badges uniquement
├── internal/              # Code interne (non exportable)
│   ├── fetcher/           # Extraction des données Root-me
│   ├── generator/         # Génération des badges SVG
│   ├── themes/            # Thèmes et icônes
│   └── cache/             # Cache fichier simple
├── pkg/                   # API publique Go
│   └── rootme/            # Interface publique pour les développeurs
├── docs/                  # Site web statique (GitHub Pages)
└── .github/workflows/     # Actions GitHub pour CI/CD
```

## Fichiers Clés

### `/internal/fetcher/rootme.go`
**Rôle :** Extraction des données depuis root-me.org
- Récupère les profils utilisateurs via scraping HTML
- Parse les statistiques, catégories, et challenges
- Utilise uniquement la stdlib Go (pas de dépendances externes)
- Gère la comparaison entre deux profils

**À modifier pour :**
- Changer les URLs Root-me
- Ajouter de nouvelles catégories
- Modifier le parsing des données

### `/internal/generator/svg.go`
**Rôle :** Génération des badges SVG
- Crée les badges avec animations CSS
- Gère les différents thèmes
- Positionne les éléments (icônes, barres de progression, texte)
- Génère les comparaisons entre utilisateurs

**À modifier pour :**
- Changer l'apparence des badges
- Modifier les positions des éléments (`x="350"` pour les pourcentages)
- Ajouter de nouveaux éléments visuels
- Modifier les animations CSS

### `/internal/themes/themes.go`
**Rôle :** Définition des thèmes et icônes
- Couleurs pour chaque thème (background, text, accent)
- Icônes SVG pour chaque catégorie Root-me
- Mapping des catégories vers leurs icônes

**À modifier pour :**
- Ajouter de nouveaux thèmes
- Changer les couleurs existantes
- Remplacer les icônes SVG
- Ajouter de nouvelles catégories

### `/cmd/server/main.go`
**Rôle :** Serveur HTTP API
- Endpoints : `/rm-gh`, `/compare`, `/api/profile`
- Gestion CORS pour utilisation web
- Cache automatique 24h
- Compatible avec les URLs Vercel existantes

**À modifier pour :**
- Ajouter de nouveaux endpoints
- Changer la logique de cache
- Modifier les paramètres d'API

### `/cmd/cli/main.go`
**Rôle :** Interface en ligne de commande
- Commandes : `badge`, `compare`, `profile`
- Options : `--nickname`, `--theme`, `--output`, `--stats`
- Génération locale de badges

**À modifier pour :**
- Ajouter de nouvelles commandes
- Changer les options CLI
- Modifier les formats de sortie

## Logique du Projet

### 1. Récupération des Données
```go
// Fetcher récupère les données depuis root-me.org
profile, err := fetcher.FetchProfile("franckferman")
// Parse automatiquement: rank, score, challenges, catégories
```

### 2. Génération SVG
```go
// Generator crée le badge SVG avec le thème choisi
opts := generator.BadgeOptions{
    Theme: "dark",
    ShowGlobalStats: true,
    Width: 380,
}
svg := generator.GenerateBadge(profile, opts)
```

### 3. Système de Thèmes
```go
// Themes fournit couleurs + icônes
theme := themes.GetTheme("midnight")  // Couleurs
icon := themes.GetCategoryIcon("Programming")  // Icône SVG
```

## Modificer le Projet

### Ajouter un Nouveau Thème
1. Ouvrir `/internal/themes/themes.go`
2. Ajouter dans `themes` map :
```go
"monTheme": {
    Background: "#123456",
    Bar:        "#789abc", 
    Accent:     "#def012",
    Text:       "#ffffff",
    Title:      "#eeeeee",
},
```

### Changer les Icônes
1. Aller sur root-me.org et récupérer les vraies icônes SVG
2. Extraire le `<path d="...">` de chaque icône
3. Remplacer dans `/internal/themes/themes.go` :
```go
"Programming": "M8 3A2 2 0...",  // Nouveau path SVG
```

### Modifier l'Apparence des Badges
1. Ouvrir `/internal/generator/svg.go`
2. Modifier les positions dans `generateCategories()` :
   - `x="350"` : Position des pourcentages
   - `width="150"` : Largeur des barres de progression
   - `y` : Positions verticales

### Ajouter une Nouvelle Catégorie
1. `/internal/themes/themes.go` : Ajouter l'icône
2. `/internal/fetcher/rootme.go` : Ajouter dans `categoryMap`
3. Tester avec un profil qui a cette catégorie

### Modifier le Serveur API
1. `/cmd/server/main.go` : Ajouter de nouveaux endpoints
2. Exemple :
```go
mux.HandleFunc("/mon-endpoint", handleMonEndpoint)
```

## Tests et Développement

### Compiler et Tester
```bash
# Compiler le CLI
go build -o cli cmd/cli/main.go

# Générer un badge de test
./cli badge --nickname=franckferman --theme=dark --output=test.svg

# Démarrer le serveur de dev
go run cmd/server/main.go
# Test: http://localhost:3000/rm-gh?nickname=franckferman&style=dark
```

### Débogage
- Vérifier les données : `./cli profile --nickname=USERNAME`
- Tester les thèmes : Générer avec différents `--theme`
- Valider SVG : Ouvrir les fichiers `.svg` dans un navigateur

## Déploiement

### Local/VPS
```bash
go build -o server cmd/server/main.go
./server  # Démarre sur port 3000
```

### GitHub Pages (Site Web)
- Le dossier `/docs/` contient le site web statique
- GitHub Actions génère automatiquement les exemples
- Déployé sur : `https://username.github.io/root-me-stats/`

### API en Production
- Compiler le serveur pour la plateforme cible
- Configurer les variables d'environnement `PORT` et `HOST`
- Déployer le binaire unique (zéro dépendances)

## Problèmes Fréquents

### "Rank 0" dans les Badges
- Problème : API Root-me ne répond pas ou profil introuvable
- Solution : Vérifier que le username existe sur root-me.org

### Pourcentages Tronqués
- Problème : `x="360"` trop à droite pour la largeur du badge
- Solution : Changer à `x="350"` ou moins dans `svg.go`

### Icônes qui ne s'affichent pas
- Problème : Path SVG invalide ou viewBox incorrect
- Solution : Vérifier les paths dans `themes.go`

### Cache Problems
- Le cache dure 24h dans `/tmp/.cache/`
- Pour forcer refresh : Supprimer le dossier cache

## Performance

### Optimisations Appliquées
- Cache fichier 24h pour éviter de spammer root-me.org
- SVG généré en mémoire (pas de fichiers temporaires)
- Stdlib uniquement (démarrage instantané)
- Binaires statiques (pas de dépendances runtime)

### Métriques Typiques
- Première génération : ~2-3 secondes (récupération données)
- Générations suivantes : ~50ms (depuis cache)
- Taille binaire : ~8-12MB (statique)
- RAM utilisée : ~5-10MB par processus

## Contributions

### Avant de Modifier
1. Tester localement avec `go run cmd/cli/main.go`
2. Vérifier que tous les thèmes fonctionnent
3. Valider les SVG générés dans un navigateur
4. Tester l'API avec différents usernames

### Structure des Commits
- Garder les commits atomiques (un changement = un commit)
- Messages descriptifs : "Fix percentage positioning in badges"
- Tester avant de pousser

Ce projet privilégie la simplicité et les performances. Éviter d'ajouter des dépendances externes sauf absolue nécessité.