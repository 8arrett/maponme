import VectorSource from "ol/source/Vector.js"
import VectorLayer from "ol/layer/Vector.js"
import Feature from "ol/Feature.js"
import Point from "ol/geom/Point.js"
import { fromLonLat } from "ol/proj"
import { circular } from "ol/geom/Polygon"

// --- This file is split from map.js to allow webpack optimization of dynamic imports

const getLayer = () => {
  return new VectorLayer()
}

const getSource = () => {
  return new VectorSource()
}

/** pos: [lon, lat] ( Coordinate{Array.<number>} ) */
const getPointFeature = (pos) => {
  return new Feature(new Point(fromLonLat(pos)))
}

const getCircleFeature = (pos, acc) => {
  return new Feature(circular(pos, acc).transform("EPSG:4326", "EPSG:3857"))
}

export { getLayer, getSource, getPointFeature, getCircleFeature }
