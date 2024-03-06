const webpack = require("webpack")

module.exports = {
  mode: process.env.NODE_ENV === "development" ? "development" : "production",
  entry: {
    main: {
      import: "./src/static/main.js",
    },
    map: {
      import: "./src/static/map.js",
    },
    point: {
      import: "./src/static/point.js",
    },
  },
  output: {
    path: __dirname + "/src/static/pack/",
    filename: "[name].bundle.js",
  },
  optimization: {
    minimize: false, // dev only

    splitChunks: {
      chunks: "all",
    },
  },
  //watch: true,  // dev only
  //watchOptions: {  // dev only
  //  poll: 300
  //}
}
