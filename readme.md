# Github User Repositories Cloner

Le GitHub User Repositories Cloner est un programme en Go conçu pour simplifier le processus de clonage de tous les dépôts d'un utilisateur GitHub. Ce script utilise les fonctionnalités Git de base pour cloner, récupérer et mettre à jour tous les dépôts d'un utilisateur spécifié, facilitant ainsi la gestion de multiples projets GitHub. L'outil est utile pour les développeurs qui souhaitent garder tous les dépôts d'un utilisateur localement à jour sur leur machine.

## Fonctionnalités

- Fonctionnalité 1 : Fetch all the repositories of a GitHub user
- Fonctionnalité 2 : Creating a csv file with all the informations of the repositories of a GitHub user
- Fonctionnalité 3 : Clone all the repositories of a GitHub user
- Fonctionnalité 4 : Update all the repositories of a GitHub user
- Fonctionnalité 5 : Zip all the repositories of a GitHub user

## Prérequis

Assurez-vous d'avoir Go installé sur votre machine. Si ce n'est pas le cas, vous pouvez le télécharger depuis [le site officiel de Go](https://golang.org/dl/).

## Installation

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/Teyik0/golang-cc.git
   cd votre-projet
   ```
2. Ajouter un fichier .env à la racine du projet :
   ```bash
   touch .env
   ```
   Remplissez le avec les information du .env.exemple
3. Runner le projet :
   ```bash
    docker-compose up --build
   ```
   ou
   ```bash
    go run .
   ```
