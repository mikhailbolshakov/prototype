package cron

import (
	"fmt"
	"github.com/go-co-op/gocron"
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

	expectedTime, _ := time.Parse(time.RFC3339, "2021-01-14T17:20:00+03:00")
	fmt.Println("expectedTime = ", expectedTime)

	quit := make(chan bool)

	j, _ := s.Every(1).Day().StartAt(expectedTime).Do(func() {
		fmt.Println("hello")
		quit <- true
	})
	fmt.Println("next run = ", j.NextRun())

	s.StartAsync()

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