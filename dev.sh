#!/bin/sh

cd /app/

npx nodemon '--watch' src/routes.go '--watch' src/store.go '--ext' go '--exec' 'go run src/main.dev.go src/routes.go src/store.go' '--signal' SIGTERM '--legacy-watch' '--polling-interval' 300 &

npx nodemon '--watch' 'src/static/main.js' '--watch' 'src/static/map.js' '--watch' 'src/static/point.js' '--ext' js '--exec' 'npx webpack --config webpack.dev.config.js --devtool eval-source-map' '--legacy-watch' '--polling-interval' 300 &

npx nodemon '--watch' 'src/static/index.htm' '--watch' 'src/static/input.css' '--ext' 'htm,css' '--exec' 'npx tailwindcss -c tailwind.config.js -i src/static/input.css -o src/static/tail.css' '--legacy-watch' '--polling-interval' 300 &
