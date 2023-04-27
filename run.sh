#!/usr/bin/bash

cd Dev/frontend || exit
if [ ! -d "node_modules" ]; then # folder does not exist
    printf "\n\e[32m###############\e[34m npm install \e[32m###############\e[0m\n\n"
    npm install
    printf "\n"
elif [ -z "$(ls -A node_modules)" ]; then # node_modules is empty
    printf "\n\e[32m###############\e[34m npm install \e[32m###############\e[0m\n\n"
    npm install
    printf "\n"
fi
printf "\e[32m##############\e[34m npm run build \e[32m##############\e[0m\n\n"
npm run build
printf "\n"

cd ../backend || exit
printf "\e[32m#################\e[34m go run \e[32m##################\e[0m\n\n"
go run .
