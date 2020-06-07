package main

import (
	"fmt"
	"math"
	"testing"
)

func Test_lotteryOneUserFromUsers(t *testing.T) {
	// initialize
	userIDs := []string{"a", "b", "c", "d", "e"}
	counts := make(map[string]int)
	for _, id := range userIDs {
		counts[id] = 0
	}

	// sampling
	const loopCnt int = 1000
	var i int
	var userID string
	for i < loopCnt {
		userID = lotteryUsersFromUsers(userIDs, 1)[0]
		counts[userID]++
		i++
	}

	// calculation
	const allowErrRate = 0.05
	var per float64
	var errRate float64
	expectPer := float64(1) / float64(len(userIDs))

	fmt.Printf("ID: count\n")
	for k, v := range counts {
		per = float64(v) / float64(loopCnt)
		errRate = math.Abs(expectPer - per)
		fmt.Printf("%s: %d ", k, v)

		fmt.Printf("%v per \n", per*100.0)
		if errRate > allowErrRate {
			t.Errorf("Error Rate is big label: %s, Error Rate: %f", k, errRate)
		}
	}
}
