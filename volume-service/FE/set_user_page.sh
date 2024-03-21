#!/usr/bin/sh

mv user/ user_ref/
npm create vite@latest user -- --template svelte
cd user
rm -r public
rm -r src/*
mv ../user_ref/* src/
npm i bootstrap
npm i svelte-spa-router
cd ..
rm -d user_ref