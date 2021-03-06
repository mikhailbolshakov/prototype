package cron

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"testing"
	"time"
)

var task = func() {
	fmt.Println("Task executed")
}

func Test_Every5Sec(t *testing.T) {

	s := gocron.NewScheduler(time.Local)
	j, err := s.Every(5).Second().Do(task)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(j.LastRun())

	s.StartBlocking()

}

func Test_SpecificTime(t *testing.T) {

	s := gocron.NewScheduler(time.Local)

	expectedTime := time.Now().Add(time.Second * 5)
	fmt.Println("expectedTime = ", expectedTime)

	quit := make(chan bool)

	j, _ := s.Every(1).Day().StartAt(expectedTime).Do(func() {
		fmt.Println("hello")
		quit <- true
	})
	fmt.Println("next run = ", j.NextRun())

	s.StartAsync()

	select {
		case <-time.After(time.Second * 10):
			fmt.Println("time out")
		case <-quit:
	}

}
//

func Test_SpecificTime2(t *testing.T) {

	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	expectedTime, _ := time.Parse(time.RFC3339, "2021-03-23T08:20:00.000Z")
	fmt.Println("expectedTime = ", expectedTime.String())
	//expectedTime := time.Date(2021, time.March, 23, 07, 46, 30, 0, time.UTC)

	quit := make(chan bool)

	j, _ := s.Every(1).Day().StartAt(expectedTime).Do(func() {
		fmt.Println("hello")
		quit <- true
	})
	fmt.Println("next run = ", j.NextRun())

	select {
	case <-time.After(time.Minute * 2):
		fmt.Println("time out")
	case <-quit:
	}

}

func Test_SpecificTimeMultiple(t *testing.T) {

	s := gocron.NewScheduler(time.Local)

	expectedTime1 := time.Now().Add(time.Second * 3)
	expectedTime2 := time.Now().Add(time.Second * 6)

	_, _ = s.Every(1).Day().StartAt(expectedTime1).Do(func() {
		fmt.Println("hello 1")
	})

	_, _ = s.Every(1).Day().StartAt(expectedTime2).Do(func() {
		fmt.Println("hello 2")
	})

	s.StartAsync()

	select {
		case <-time.After(time.Second * 10):
			fmt.Println("time out")
	}

}

func Test_SpecificTimeMultipleAfterStart(t *testing.T) {

	s := gocron.NewScheduler(time.Local)

	expectedTime1 := time.Now().Add(time.Second * 10)
	expectedTime2 := time.Now().Add(time.Second * 20)

	_, _ = s.Every(1).Day().StartAt(expectedTime1).Do(func() {
		fmt.Println("hello 1")
	})

	s.StartAsync()

	time.Sleep(time.Second * 2)

	_, _ = s.Every(1).Day().StartAt(expectedTime2).Do(func() {
		fmt.Println("hello 2")
	})

	select {
	case <-time.After(time.Minute * 1):
		fmt.Println("time out")
	}

}

func Test_SpecificTimeZones(t *testing.T) {

	s := gocron.NewScheduler(time.Local)

	//expectedTime := time.Date(2021, 01, 15, 15, 58, 00, 00, time.Local)
	expectedTime, _ := time.Parse(time.RFC3339, "2021-01-15T13:14:00.000Z")

	j, _ := s.Every(1).Day().StartAt(expectedTime).Do(func() {
		fmt.Println("hello 1")
	})
	log.Println(j.NextRun())

	s.StartAsync()

	select {
	case <-time.After(time.Minute * 10):
		fmt.Println("time out")
	}

}

func Test_SpecificTime_WithIncorrectParams(t *testing.T) {

	quit := make(chan bool)

	action := func(s1, s2 string) {
		fmt.Printf("fired: %s, %s", s1, s2)
		quit <- true
	}

	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	j, err := s.Every(1).Day().StartAt(time.Now().UTC().Add(time.Second * 10)).Do(action, "s1")
	if err != nil {
		t.Fatal()
	}
	fmt.Println("next run = ", j.NextRun())

	select {
	case <-time.After(time.Second * 20):
		fmt.Println("time out")
		t.Fatal()
	case <-quit:
		fmt.Println("passed")
	}

}