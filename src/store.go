package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type UserStored struct {
	key    string // Private Key
	lon    string // Longitude
	lat    string // Latitude
	acc    string // Accuracy
	expire int64  // Expire
	mod    int64  // Modified
}

type UserCred struct {
	id  string
	key string
}

type UserLocation struct {
	lon string
	lat string
	acc string
	mod int64
}

var Database = make(map[string]UserStored)
var DatabaseMux = &sync.RWMutex{}

const DatabaseCooldownTimer = int64(2) // 2 seconds

// ------------------------------------------------------------
//   Database access functions
// ------------------------------------------------------------

func setUser(id string, user UserStored) {
	DatabaseMux.Lock()
	defer DatabaseMux.Unlock()
	Database[id] = user
}

func getUser(id string) UserStored {
	DatabaseMux.RLock()
	defer DatabaseMux.RUnlock()
	return Database[id]
}

func removeUser(id string) {
	DatabaseMux.Lock()
	defer DatabaseMux.Unlock()
	delete(Database, id)
}

func userExists(id string) bool {
	DatabaseMux.RLock()
	defer DatabaseMux.RUnlock()
	_, exists := Database[id]
	return exists
}

// ------------------------------------------------------------
//   Random generator functions
// ------------------------------------------------------------

func generateUserId() string {

	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Incoming! Caught system error in store.go::generateUserId() ->")
		fmt.Fprintln(os.Stderr, err)
	}
	base := base64.StdEncoding.EncodeToString(buf)
	replacer := strings.NewReplacer("/", "-", "+", "=")
	return replacer.Replace(base)
}

func generateUserKey() string {

	buf := make([]byte, 18)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Incoming! Caught system error in store.go::generateUserKey() ->")
		fmt.Fprintln(os.Stderr, err)
	}
	return base64.StdEncoding.EncodeToString(buf)
}

// ------------------------------------------------------------
//   Garbage collector
// ------------------------------------------------------------

func DatabaseGarbageCollector() {
	for {
		time.Sleep(10 * time.Minute)
		removeGarbage(getInactiveUserIds())
	}
}

func getInactiveUserIds() []string {
	var expiredUserIds []string
	timestamp := time.Now().Unix()

	DatabaseMux.RLock()
	defer DatabaseMux.RUnlock()

	for k, v := range Database {
		if v.expire < timestamp {
			expiredUserIds = append(expiredUserIds, k)
		}
	}
	return expiredUserIds
}

func removeGarbage(expiredUserIds []string) {

	for _, id := range expiredUserIds {
		removeUser(id)
	}
}

// ------------------------------------------------------------
//   API call functions
// ------------------------------------------------------------

func createUser() UserCred {

	key := generateUserKey()
	reqTime := time.Now().Unix()
	expire := reqTime + 180 // 3 min init expiration

	id := ""
	for {
		id = generateUserId()
		if !userExists(id) {
			break
		}
		fmt.Fprintln(os.Stderr, "CreateUser actually found a collision!")
		fmt.Fprintln(os.Stderr, id)
	}

	user := UserStored{key: key, expire: expire, mod: 0}
	setUser(id, user)

	return UserCred{id: id, key: key}
}

func deactivateUser(id string, key string) error {

	user := getUser(id)
	if user.expire < time.Now().Unix() {
		return errors.New("expired")
	}
	if user.key != key {
		return errors.New("auth")
	}

	removeUser(id)
	return nil
}

func setUserLocation(id string, key string, loc UserLocation) (int64, error) {

	reqTime := time.Now().Unix()
	return setUserLocationExp(id, key, loc, reqTime + 7200) // 2 hours
}

func setUserLocationExp(id string, key string, loc UserLocation, expire int64) (int64, error) {

	reqTime := time.Now().Unix()
	user := getUser(id)
	if user.expire < reqTime {
		return 0, errors.New("expired")
	}
	if user.mod > reqTime-DatabaseCooldownTimer {
		return 0, errors.New("cooldown")
	}
	if user.key != key {
		return 0, errors.New("auth")
	}

	user.lon = loc.lon
	user.lat = loc.lat
	user.acc = loc.acc
	user.expire = expire
	user.mod = reqTime
	setUser(id, user)

	return reqTime, nil
}

func getUserLocation(id string) (UserLocation, bool) {

	user := getUser(id)
	if user.expire < time.Now().Unix() {
		null := UserLocation{}
		return null, true
	}

	location := UserLocation{lon: user.lon, lat: user.lat, acc: user.acc, mod: user.mod}
	return location, false
}
