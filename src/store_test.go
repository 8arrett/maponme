package main

// There is a set of current tests that will run but they were not polished.
// This file is mostly updated for exploratory testing of large data sets.

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestConcurrentStorageActivity(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	startTime := time.Now().Unix()
	endTime := startTime + 1
	err := make(chan int, 100)
	done := make(chan int, 100000000)
	wg := new(sync.WaitGroup)
	location := UserLocation{lon: "123.45678", lat: "12.34567", acc: "100"}

	testSz := 1000
	for i := 0; i < testSz; i++ {
		user := createUser()
		setUserLocationExp(user.id, user.key, location, startTime+180)
		go testReadUser(done, wg, err, user.id, location)
		go testWriteUser(done, wg, user.id, user.key, location)
		wg.Add(2)

		if time.Now().Unix() > endTime {
			if i+1 != len(Database) {
				t.Fatalf(`Ended with %d users running and %d users in map`, i+1, len(Database))
			}

			endConcurrentStorageActivity(t, i+1, done, wg, err)
			return
		}
	}
	log.Printf("TestConcurrentStorageActivity: Generated full test users")
	endConcurrentStorageActivity(t, testSz, done, wg, err)
}

func endConcurrentStorageActivity(t *testing.T, n int, done chan<- int, wg *sync.WaitGroup, err chan int) {

	select {
	case <-err:
		t.Fatalf(`Collision detected`)
	default:
		err <- 0
	}

	log.Printf("Tested rapid concurrency on %d users running", n)

	benchUser := createUser()
	location := UserLocation{lon: "123.45678", lat: "12.34567", acc: "100"}
	benchStartS := time.Now()
	setUserLocationExp(benchUser.id, benchUser.key, location, 0)
	log.Print(time.Now().Sub(benchStartS), " to write to database")
	benchStartG := time.Now()
	getUserLocation(benchUser.id)
	log.Print(time.Now().Sub(benchStartG), " to read from database")

	for i := 0; i < n; i++ {
		done <- 1
		done <- 1
	}

	wg.Wait()
	return
}

func testReadUser(done <-chan int, wg *sync.WaitGroup, err chan<- int, user string, loc UserLocation) {

	for {
		select {
		case <-done:
			wg.Done()
			return
		case <-time.After(200 * time.Millisecond):
			res, _ := getUserLocation(user)
			if res.lon != loc.lon || res.lat != loc.lat || res.acc != loc.acc {
				err <- 1
			}
		}
		//time.Sleep(200 * time.Millisecond)
	}
}

func testWriteUser(done <-chan int, wg *sync.WaitGroup, user string, key string, loc UserLocation) {

	for {
		select {
		case <-done:
			deactivateUser(user, key)
			wg.Done()
			return
		case <-time.After(500 * time.Millisecond):
			setUserLocationExp(user, key, loc, time.Now().Unix()+180)
			//time.Sleep(500 * time.Millisecond)
		}
	}
}

func TestConcurrentGarbageActivity(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	//Database = make(map[string]UserStored)

	counterZ := 0
	for range Database {
		counterZ += 1
		//_, _ :=
		//log.Print(k,v)
	}
	log.Printf("Size of init database: %d", counterZ)

	location := UserLocation{lon: "123.45678", lat: "12.34567", acc: "100"}
	startTime := time.Now().Unix()

	sz := 2000
	userCreds := make([]UserCred, sz)

	for i := 0; i < sz; i++ {
		user := createUser()
		userCreds = append(userCreds, UserCred{id: user.id, key: user.key})

		if i%2 == 0 {
			setUserLocationExp(user.id, user.key, location, startTime+180)
		} else {
			setUserLocationExp(user.id, user.key, location, startTime-1)
		}
	}

	counterA := 0
	for range Database {
		counterA += 1
		//_, _ :=
		//log.Print(k,v)
	}
	log.Printf("Size of database: %d", counterA)

	for i := 0; i < sz; i++ {
		go func() {
			for {
				getUserLocation(userCreds[i].id)
				time.Sleep(150)
			}
		}()
		go func() {
			for {
				getUserLocation(userCreds[i].id)
				time.Sleep(150)
			}
		}()
		// go func () {
		// 	for {
		// 		getUserLocation(userCreds[i].id)
		// 		time.Sleep(100)
		// 	}
		// }()
	}

	removeGarbage(getInactiveUserIds())

	counter := 0
	counterB := 0
	for k, _ := range Database {
		counter += 1
		for _, v := range userCreds {
			if v.id == k {
				counterB += 1
			}
		}
	}
	log.Printf("Size of database: %d", counter)
	log.Printf("Size of database found: %d", counterB)
	// run reads/writes on all of them
	// garbage collect
	// test half have been deleted
	// test other half are intact
}
