package main

import "os"
import "bufio"
import "strings"
import "strconv"
import "reflect"
import "fmt"

//
// Utility functions because I don't know where to find these elsewhere
//

// Splits a slice into [head | tail]
func shift(queue []int) (int, []int) {
	switch len(queue) {
	case 0:
		return -1, []int{}
	case 1:
		return queue[0], []int{}
	}
	return queue[0], queue[1:]
}

// Prepend an element to a slice and return it
func unshift(ball int, queue []int) []int {
	queue = append([]int{ball}, queue...)
	return queue
}

//
// Functions that model the operation of a ball clock
//

// Return a queue slice initialized with `depth` balls
func makeQueue(depth int) []int {
	queue := []int{}
	for i := 0; i < depth; i++ {
		queue = append(queue, i)
	}
	return queue
}

// Add a ball to a "track" slice, and return both the track and any balls that should fall off
func passBallToTrack(ball int, track []int, depth int) (int, []int, []int) {
	track = unshift(ball, track)
	returning := []int{}
	if len(track) == depth {
		returning = append(returning, track...)
		track = []int{}
	}
	ball, returning = shift(returning)
	return ball, returning, track
}

// Set up and then run a ball clock for twelve hours to figure out the ball pattern
// This pattern will repeat every twelve hours, so we only have to run it once.
func getTwelveHourTransform(balls int) []int {
	var queue, oneMinuteTrack, fiveMinuteTrack, oneHourTrack, returned []int
	var ball int
	queue = makeQueue(balls)
	for {
		ball, queue = shift(queue)
		ball, returned, oneMinuteTrack = passBallToTrack(ball, oneMinuteTrack, 5)
		if ball >= 0 {
			// minuteTrack returned ball for fiveMinuteTrack
			queue = append(queue, returned...)
			ball, returned, fiveMinuteTrack = passBallToTrack(ball, fiveMinuteTrack, 12)
			if ball >= 0 {
				// fiveMinuteTrack returned ball for oneHourTrack
				queue = append(queue, returned...)
				ball, returned, oneHourTrack = passBallToTrack(ball, oneHourTrack, 12)
				if ball >= 0 {
					// oneHourTrack returned ball - twelve hours have passed
					returned = append(returned, ball)
					queue = append(queue, returned...)
					break
				} // oneHourTrack
			} // fiveMinuteTrack
		} // oneMinuteTrack
	}
	return queue
}

// Transform the queue using the twelve hour pattern provided
func transformQueue(queue, pattern []int) []int {
	var transformedQueue []int
	for _, v := range pattern {
		transformedQueue = append(transformedQueue, queue[v])
	}
	return transformedQueue
}

// Use the twelve hour pattern to transform the queue repeatedly until it returns to the original pattern
func mainLoop(balls int) int {
	var halfDays int
	var queue, originalQueue, transformPattern []int
	queue, originalQueue = makeQueue(balls), makeQueue(balls)
	transformPattern = getTwelveHourTransform(balls)
	for {
		// Count the days until the queue returns to its original pattern
		halfDays++
		queue = transformQueue(queue, transformPattern)
		if reflect.DeepEqual(queue, originalQueue) {
			break
		}
	}
	return halfDays / 2
}

//
// Ball Clock problem
// Read number of balls (one number per line) from stdin, 0 to end
// Output according to the required format

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		balls, err := strconv.Atoi(strings.TrimRight(line, "\n"))
		if balls == 0 && err == nil {
			// exit on zero input, other cases are not defined
			return
		}
		if balls >= 27 && balls <= 127 {
			fmt.Println(balls, "balls cycle after", mainLoop(balls), "days.")
		}
	}
}
