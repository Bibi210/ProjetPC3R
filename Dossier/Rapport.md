---
title: Rapport Projet PC3R - Implémentation d'une application AJAX
subtitle: ShitPostLand - Un mini réseau social pour shitposter
author: Dibassi Brahima, Said Mohammad Zuhair
date: 27 Avril, 2023
lang: fr
geometry:
  - margin = 1.2in
mainfont: Palatino
sansfont: Helvetica
monofont: Menlo
fontsize: 12pt
urlcolor: NavyBlue
include-before: | # Texte avant la table des matières
    \newpage
numbersections: true # Numéros de sections
toc: true # Table des matières
tableofcontents: true # Table des matières
---
\newpage

# Lancement du projet en local

Afin de lancer le projet en local, il faut se placer a la racine du projet et lancer la commande suivante:

```bash
  bash RunProject.sh
```

Et ensuite se rendre sur l'adresse suivante: [http://localhost:25565/](http://localhost:25565/)

# Description de l'archive

Le code du projet se trouve dans le dossier `Dev/` et est divisé en 2 sous-dossiers:

- `Frontend/` : Contient le code du client écrit en ReactJS.
- `Backend/` : Contient le code du serveur écrit en GoLang.
  - `Helpers` :
    - `errors.go` : Contient le code de gestion des erreurs.
    - `inoutFormats.go` : Contient les formats d'entrée et de sortie des requêtes en JSON encodés en structure GoLang.
  - `Database/` : Contient le code manipulant la base de données SQLite.




\newpage