@echo off

cd web
npm i
npx browserify .\index.js > ../assets/bundle.js
npx uglify-js --compress  -o ..\assets\bundle.js ..\assets\bundle.js
