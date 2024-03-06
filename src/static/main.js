/**
 * Helper functions
 */

const isAnimationRunning = (el) => {
  let running = false
  el.getAnimations().forEach((a) => {
    if (a.playState != "finished") running = true
  })
  return running
}

const throttleClick = (fn, delay) => {
  let timerId = null
  return (...args) => {
    if (timerId == null) {
      fn(...args)
      timerId = setTimeout(() => {
        timerId = null
      }, delay)
    }
  }
}

/**
 * GUI classes
 */

const tiles = class {
  static #reset() {
    document.getElementById("introTile").classList.add("hidden")
    document.getElementById("permissionTile").classList.add("hidden")
    document.getElementById("mapTile").classList.add("hidden")
  }
  static viewIntro() {
    this.#reset()
    document.getElementById("introTile").classList.remove("hidden")
    document.getElementById("titleDiv").classList.remove("fadeinTitle")
    document.getElementById("titleDiv").classList.add("fadeinTitle")
  }
  static viewMap() {
    this.#reset()
    document.getElementById("mapTile").classList.remove("hidden")
  }
  static viewPermission() {
    this.#reset()
    document.getElementById("permissionTile").classList.remove("hidden")
    document.getElementById("permissionSeeking").classList.remove("hidden")
    document.getElementById("permissionDenied").classList.add("hidden")
  }
}
const closeButton = class {
  static show() {
    document.getElementById("closeLink").classList.remove("hidden")
  }
  static hide() {
    document.getElementById("closeLink").classList.add("hidden")
  }
}
const loadingIcon = class {
  static show() {
    document.getElementById("loadingIcon").classList.remove("hidden")
  }
  static hide() {
    document.getElementById("loadingIcon").classList.add("hidden")
  }
}
const copyButton = class {
  static show() {
    document.getElementById("copyLinkButton").classList.remove("hidden")
  }
  static hide() {
    document.getElementById("copyLinkButton").classList.add("hidden")
  }
  static light() {
    const btn = document.getElementById("copyLinkButton")
    const txt = document.getElementById("copyLinkText")
    txt.classList.remove("hidden")
    txt.classList.remove("animate-copyText")
    btn.classList.remove("animate-copyButton")
    btn.offsetWidth // Force css redraw, https://stackoverflow.com/a/63561659
    btn.classList.add("animate-copyButton")
    txt.classList.add("animate-copyText")
    setTimeout(() => {
      if (isAnimationRunning(txt)) return
      txt.classList.add("hidden")
      btn.classList.remove("animate-copyButton")
      txt.classList.remove("animate-copyText")
    }, 3000 + 500)
  }
}

const popupMessage = class {
  static open(msg) {
    document.getElementById("textPopup").innerHTML = msg
    document.getElementById("messagePopup").classList.remove("hidden")
    this.#fadeout()
  }
  static close() {
    document.getElementById("messagePopup").classList.add("hidden")
  }
  static #fadeout() {
    const el = document.getElementById("messagePopup")
    el.classList.add("popupMessageFade")
    setTimeout(() => {
      if (isAnimationRunning(el)) return
      el.classList.add("hidden")
      el.classList.remove("popupMessageFade")
    }, 10000 + 500)
  }
}

/**
 * MapIO dynamically loads map.js (heavy library dependencies) when its interfaces are used.
 */

class MapIO {
  constructor() {
    this.mapIO = undefined
  }
  async load() {
    if (this.mapIO == undefined) {
      this.mapIO = (await import("./map.js")).mapIO
      this.mapIO.logInit()
    }
  }
  async setPos(lon, lat, acc) {
    await this.load()
    this.mapIO.setPos(lon, lat, acc)
  }
  async recenter() {
    await this.load()
    this.mapIO.recenter()
  }
  async checkCenter() {
    await this.load()
    this.mapIO.checkCenter()
  }
  async resetMap() {
    await this.load()
    this.mapIO.resetMap()
  }
}
const mapIO = new MapIO()
window.mainMap = mapIO // allows console access to packed scope

/**
 * Host object projects the user's location data onto the local and remote map.
 * It manages the lifecycle of the published link.
 */

class Host {
  constructor() {
    this.active = false
    this.expiredCount = 0
    this.acc = null
    this.mod = 0
  }
  async #getHostID() {
    fetch("/api/", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.id != undefined && data.key != undefined) {
          this.active = true
          this.id = data.id
          this.key = data.key
          document.location.hash = this.id
        } else {
          console.log("POST did not return success")
        }
      })
      .catch((err) => {
        err.dev = "fetch failed"
        console.log("POST: " + JSON.stringify(err))
      })
  }
  async #queueUpdate(lon, lat, acc) {
    if (this.acc != null && this.acc < acc) return

    this.timestamp = Math.floor(Date.now() / 1000)
    this.lon = lon
    this.lat = lat
    this.acc = acc
  }
  async #sendUpdate() {
    if (!this.active) {
      this.#getHostID()
      if (!this.active) return
    }

    if (this.acc == null) return

    fetch("/api/" + this.id, {
      method: "PUT",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        Key: this.key,
        Lat: this.lat,
        Lon: this.lon,
        Acc: this.acc,
      }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.error !== undefined) {
          loadingIcon.show()
          if (data.error == "data-expired") {
            this.expiredCount += 1
            if (this.expiredCount > 1) {
              popupMessage.open(
                "Your map has disappeared!" +
                  " <br />This could be caused by not visiting for over an hour, or we made a mistake. (Sorry!)" +
                  " <br />Feel free to create a new one!"
              )
              this.stop()
            } else {
              this.#sendUpdate()
            }
          }
          return
        }
        this.expiredCount = 0
        loadingIcon.hide()
        this.mod = data.mod
      })
      .catch((err) => {
        loadingIcon.show()
        err.dev = "fetch failed"
        console.log("PUT: " + JSON.stringify(err))
      })

    mapIO.setPos(this.lon, this.lat, this.acc)
    mapIO.checkCenter()
    this.lon = null
    this.lat = null
    this.acc = null
  }
  async #sendEndSignal(id, key) {
    if (id == undefined) return

    fetch("/api/" + id, {
      method: "DELETE",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        Key: key,
      }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.error !== undefined) {
          console.log("DELETE: " + JSON.stringify(data))
        }
      })
      .catch((err) => {
        err.dev = "fetch failed"
        console.log("DELETE: " + JSON.stringify(err))
      })
  }
  async #updateLoop() {
    this.#sendUpdate()
    this.syncID = setTimeout(() => {
      this.#updateLoop()
    }, 2500)
  }
  start() {
    mapIO.load()
    tiles.viewPermission()

    if (!("geolocation" in navigator)) {
      document.getElementById("permissionSeeking").classList.add("hidden")
      document.getElementById("permissionDenied").classList.remove("hidden")
    }

    closeButton.show()
    copyButton.show()

    this.watchID = navigator.geolocation.watchPosition(
      (pos) => {
        tiles.viewMap()
        mapIO.checkCenter()
        const lon = pos.coords.longitude.toString()
        const lat = pos.coords.latitude.toString()
        const acc = Math.ceil(pos.coords.accuracy).toString()
        this.#queueUpdate(lon, lat, acc)
      },
      (err) => {
        if (err.code == 1 || err.code == 2) {
          // PERMISSION_DENIED or POSITION_UNAVAILABLE
          tiles.viewPermission()
          document.getElementById("permissionSeeking").classList.add("hidden")
          document.getElementById("permissionDenied").classList.remove("hidden")
          navigator.geolocation.clearWatch(this.watchID)
          this.watchID = undefined
        }
        // TIMEOUT drops
      },
      {
        enableHighAccuracy: true,
        timeout: 300,
        maximumAge: 0,
      }
    )

    this.#updateLoop()
  }
  stop() {
    this.#sendEndSignal(this.id, this.key)
    navigator.geolocation.clearWatch(this.watchID)
    this.watchID = undefined
    clearTimeout(this.syncID)
    this.syncID = undefined
    this.id = undefined
    this.key = undefined
    this.lon = undefined
    this.lat = undefined
    this.acc = null
    this.mod = 0
    tiles.viewIntro()
    closeButton.hide()
    copyButton.hide()
    loadingIcon.hide()
    window.location.hash = ""
    mapIO.resetMap()
    this.active = false
  }
}
const host = new Host()
window.host = host // allows console access to packed scope

/**
 * Reader object pulls in another user's location data and projects it onto the local map.
 */

class Reader {
  constructor() {
    this.active = false
    this.id = ""
    this.expiredCount = 0
    this.mod = 0
  }
  async #update() {
    let headers = { "Content-type": "application/json" }
    if (this.mod != 0) {
      const digits =
        "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
      const tagIdx = this.mod % 262143
      headers["If-None-Match"] =
        digits[tagIdx % 64] +
        digits[(tagIdx >> 6) % 64] +
        digits[(tagIdx >> 12) % 64]
    }
    fetch("/api/" + this.id, {
      method: "GET",
      headers,
    })
      .then((res) => {
        if (res.status == 304) return { cache: true }
        return res.json()
      })
      .then((data) => {
        if (data.cache) return
        if (data.error !== undefined) {
          loadingIcon.show()
          if (data.error == "data-expired") {
            this.expiredCount += 1
            if (this.expiredCount > 1) {
              popupMessage.open(
                "Your friend's map has disappeared!" +
                  " <br />This could be caused by them closing it, or simply not having looked at it for awhile." +
                  " <br />Of course, it could also be our mistake. (Sorry!)"
              )
              this.stop()
            } else {
              this.#update()
            }
          }
          return
        }
        this.expiredCount = 0
        loadingIcon.hide()
        mapIO.setPos(data.lon, data.lat, data.acc)
        mapIO.checkCenter()
        this.mod = data.mod
      })
      .catch((err) => {
        loadingIcon.show()
        err.dev = "fetch failed"
        console.log("GET: " + JSON.stringify(err))
      })
  }
  setUser(id) {
    this.id = id
  }
  start() {
    mapIO.load()
    tiles.viewMap()
    this.#startLoop()
    this.active = true
  }
  #startLoop() {
    this.#update()
    this.timeoutID = setTimeout(() => {
      this.#startLoop()
    }, 3000)
  }
  stop() {
    this.id = undefined
    if (this.timeoutID != undefined) clearTimeout(this.timeoutID)
    this.timeoutID = undefined
    this.expired = false
    this.mod = 0

    window.location.hash = ""
    tiles.viewIntro()
    loadingIcon.hide()
    mapIO.resetMap()
    this.active = false
  }
}
const reader = new Reader()
window.reader = reader // allows console access to packed scope

/**
 *  Click listeners and link autoload
 */

const modifiedTimerLoop = () => {
  time = 0
  if (reader.active) {
    time = Math.floor(Date.now() / 1000) - reader.mod
  } else if (host.active) {
    time = Math.floor(Date.now() / 1000) - host.mod
  }

  words = ""
  if (time < 5) {
    words = "moments"
  } else if (time < 60) {
    words = time + " seconds"
  } else if (time < 120) {
    words = "a minute"
  } else if (time < 3600) {
    words = Math.floor(time / 60) + " minutes"
  } else if (time < 7200) {
    words = "an hour"
  } else if (time < 10800) {
    words = Math.floor(time / 3600) + " hours"
  } else {
    words = "long"
  }
  document.getElementById("accessTimerTime").innerText = words
  setTimeout(modifiedTimerLoop, 1000)
}

const initActions = () => {
  if (window.location.hash.length > 8) {
    reader.setUser(window.location.hash.substr(1, 8))
    reader.start()
  }

  modifiedTimerLoop()

  document.getElementById("startButtonContainer").addEventListener(
    "click",
    throttleClick(() => {
      host.start()
    }, 400)
  )
  document.getElementById("startButton2").addEventListener(
    "click",
    throttleClick(() => {
      host.start()
    }, 400)
  )
  document.getElementById("closeLink").addEventListener("click", () => {
    host.stop()
  })

  document.getElementById("centerMap").addEventListener("click", () => {
    mapIO.recenter()
    mapIO.checkCenter()
  })
  document.getElementById("copyLinkButton").addEventListener("click", () => {
    const trail = host.id !== undefined ? "#" + host.id : ""
    const url = "https://mapon.me/" + trail
    navigator.clipboard.writeText(url).then(() => {
      document.getElementById("copyLinkUrl").innerText = url
      copyButton.light()
    })
  })

  document.getElementById("messagePopup").addEventListener("click", () => {
    popupMessage.close()
  })

  addEventListener("beforeunload", () => {
    if (host.active) host.stop()
  })
}

if (document.readyState === "complete") {
  initActions()
} else {
  window.addEventListener("DOMContentLoaded", initActions)
}
