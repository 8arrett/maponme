import { Map, View } from "ol"
import OSM from "ol/source/OSM.js"
import TileLayer from "ol/layer/Tile.js"
import { fromLonLat, transform } from "ol/proj.js"

// Some likely unrelated good reading:
//  https://openlayers.org/en/latest/apidoc/module-ol_proj.html#.useGeographic
//  https://openlayers.org/en/latest/examples/geographic.html (Great guide for popups too)
//  https://openlayers.org/workshop/en/mobile/geolocation.html

class MapIO {
  constructor() {
    const olCss = document.createElement("link")
    olCss.rel = "stylesheet"
    olCss.href = "/s/ol.css"
    document.querySelector("head").appendChild(olCss)

    this.map = new Map({
      target: "mapTile",
      layers: [
        new TileLayer({
          source: new OSM(),
        }),
      ],
      view: new View({
        center: [0, 0],
        zoom: 2,
      }),
    })
    this.tileID = this.map.getLayers().getArray()[0].ol_uid

    this.map.on("moveend", () => {
      this.checkCenter()
    })

    this.resetMap()
  }

  logInit() {
    console.log("MapIO imported")
  }

  setPos(lon, lat, acc) {
    // prevent sending map into space when user hasn't posted position yet
    if (acc == 0) return

    this.pos = [parseFloat(lon), parseFloat(lat)]
    this.acc = acc

    if (!this.hasZoomed) {
      this.recenter()
      this.hasZoomed = true
    }

    import("./point.js").then((lazyLoad) => {
      if (!this.hasPoint) {
        this.sourceVector = lazyLoad.getSource()
        this.layerVector = lazyLoad.getLayer()
        this.layerVector.setSource(this.sourceVector)
        this.map.addLayer(this.layerVector)
        this.hasPoint = true
      }

      this.sourceVector.clear()
      this.sourceVector.addFeatures([
        lazyLoad.getPointFeature(this.pos),
        lazyLoad.getCircleFeature(this.pos, this.acc),
      ])
    })
  }

  // subjectively zoom 17 looks like a nice max
  //   this.map.getView().setZoom(17)
  // Estimated zoom in meters table:
  //   https://wiki.openstreetmap.org/wiki/Zoom_levels

  recenter() {
    const tile = document.getElementById("mapTile")
    const sz = Math.min(tile.offsetWidth, tile.offsetHeight)

    this.map.getView().setCenter(fromLonLat(this.pos))

    const res = Math.max((this.acc * 3) / sz, 1.2)
    this.map.getView().setResolution(res)

    if (this.pos[0] == 0 && this.pos[1] == 0)
      // Position unset, so default map of the world
      this.map.getView().setZoom(2)
  }

  checkCenter() {
    const center = transform(
      this.map.getView().getCenter(),
      "EPSG:3857",
      "EPSG:4326"
    )

    const nearCenter =
      center[0] > this.pos[0] + 1e-6 ||
      center[0] < this.pos[0] - 1e-6 ||
      center[1] > this.pos[1] + 1e-6 ||
      center[1] < this.pos[1] - 1e-6

    if (nearCenter) {
      document.getElementById("centerMap").classList.remove("hidden")
      return
    }
    document.getElementById("centerMap").classList.add("hidden")
  }

  resetMap() {
    const layers = [...this.map.getLayers().getArray()]
    layers.forEach((layer) => {
      if (this.tileID !== layer.ol_uid) {
        this.map.removeLayer(layer)
      }
    })

    this.pos = [0, 0]
    this.acc = 0

    this.hasZoomed = false
    this.hasPoint = false

    this.sourceVector = undefined
    this.layerVector = undefined
  }
}

const mapIO = new MapIO()
window.mapIO = mapIO
export { mapIO }
