#!/usr/bin/sh

mv admin/ admin_ref/
npm create vite@latest admin -- --template svelte
cd admin
rm -r public
rm -r src/*
mv ../admin_ref/* src/
npm i bootstrap
npm i svelte-spa-router
cd ..
rm -d admin_ref/