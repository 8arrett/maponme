#!/bin/sh

if [[ ${PWD##*/} != "app" ]]
then
    echo ""
    echo "  Please run from the app directory"
    echo ""
    exit 1
fi

echo ""
echo "###  Cleaning distribution directory.."
rm -R dist/*

echo ""
echo "###  Starting go build.."
go build -o dist/apiServer -ldflags="-s -w" src/main.go src/routes.go src/store.go

echo ""
echo "###  Starting tailwind compile.."
npx tailwindcss -i ./src/static/input.css -o ./src/static/tail.css --minify

echo ""
echo "###  Starting webpack.."
rm ./src/static/pack/*
NODE_ENV=production npm run build

echo ""
echo "### Moving.."
mkdir dist/static/
cp -R src/static/* dist/static/

echo ""
echo "###  Cleaning.."
rm dist/static/input.css
rm dist/static/*.js

echo ""
echo "###  Minifying.."
npx html-minifier dist/static/index.htm -o dist/static/index.htm --remove-comments --collapse-whitespace --minify-js --minify-css
node -e "const CleanCSS = require('clean-css'); new CleanCSS({returnPromise: true}).minify(['dist/static/ol.css']).then((out) => console.log(out.styles))" > dist/static/new.css
sed -i '/^$/d' dist/static/new.css
mv dist/static/new.css dist/static/ol.css  # brace expansion not working on alpine

echo ""
echo "###  Zipping.."
for f in $(find dist/static -type f)
do
    gzip -k -9 "$f"
done

echo ""
echo "### Zip cleaning.."  # tiny files are larger when compressed
rm dist/static/robots.txt.gz
rm dist/static/error.htm.gz

echo ""
echo "### Creating index directory.."
mkdir dist/static/index/
mv dist/static/index.htm dist/static/index/index.htm
mv dist/static/index.htm.gz dist/static/index/index.htm.gz
mv dist/static/favicon.ico dist/static/index/favicon.ico
mv dist/static/favicon.ico.gz dist/static/index/favicon.ico.gz
mv dist/static/robots.txt dist/static/index/robots.txt

echo ""



### Decided against encoding page literal, more performant for reverse proxy to serve static

# echo ""
# echo "Cleaning last html migration"
# rm -f src/index.go

# echo ""
# echo "Starting to write new index string"
# cat <<END >> src/index.go
# package main

# // THIS FILE IS AUTOMATICALLY REGENERATED

# func renderIndex() []byte {
# END
# html=$(cat src/static/index.htm | tr -d '\n' | tr -d '\r' | tr -s "\"" | sed "s/\"/\\\\\"/g")
# echo "  return []byte(\""$html"\")" >> src/index.go
# echo "}" >> src/index.go