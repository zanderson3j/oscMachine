package main

import (
	"math/rand"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

var all_freqs = [...]float32{
	16.35, 17.32, 18.35, 19.45, 20.60, 21.83, 23.12, 24.50, 25.96, 27.50, 29.14, 30.87,
	32.70, 34.65, 36.71, 38.89, 41.20, 43.65, 46.25, 49.00, 51.91, 55.00, 58.27, 61.74,
	65.41, 69.30, 73.42, 77.78, 82.41, 87.31, 92.50, 98.00, 103.83, 110.00, 116.54, 123.47,
	130.81, 138.59, 146.83, 155, 56, 164.81, 174.61, 185.00, 196.00, 207.65, 220.00, 233.08, 246.94,
	261.63, 277.18, 293.66, 311.13, 329.63, 349.23, 369.99, 392.00, 415.30, 440.00, 466.16, 493.88,
	523.25, 554.37, 587.33, 622.25, 659.25, 698.46, 739.99, 783.99, 830.61, 880.00, 932.33, 987.77,
	1046.50, 1108.73, 1174.66, 1244.51, 1318.51, 1396.91, 1479.98, 1567.98, 1661.22, 1760.00, 1864.66, 1975.53,
	2093.00, 2217.46, 2349.32, 2489.02, 2637.02, 2793.83, 2959.96, 3135.96, 3322.44, 3520.00, 3729.31, 3951.07,
	4186.01, 4434.92, 4698.63, 4978.03, 5274.04, 5587.65, 5919.91, 6271.93, 6644.88, 7040.00, 7458.62, 7902.13,
}

// Harmonic Minor in D (D, E, F, G, A, Bb, C#, D)
var harm_min_d = [...]float32{
	293.66, 329.63, 349.23,
	392.00, 440.00, 466.16,
	554.37, 587.33, 659.25,
	698.46, 783.99, 880.00,
	932.33, 1108.73, 1174.66,
}

// Major C
var maj_c = [...]float32{
	261.63, 293.66, 329.63,
	349.23, 392.00, 440.00,
	493.88, 523.25, 587.33,
	659.25, 698.46, 783.99,
	880.00, 987.77, 1046.50,
}

func getNextNote(nextNote int, arr_length int) int {
	number := rand.Intn(100) + 1

	if number < 10 {
		nextNote -= 2
	} else if number < 45 {
		nextNote -= 1
	} else if number < 65 {

	} else if number < 90 {
		nextNote += 1
	} else {
		nextNote += 2
	}

	if nextNote < 0 {
		nextNote = 0
	} else if nextNote >= arr_length {
		nextNote = arr_length - 1
	}

	return nextNote
}

func getNextRingAndFreqMod() (float32, float32) {
	number := rand.Intn(100) + 1
	nextRing := 0.0
	nextFreqMod := 0.0

	if number < 33 {
		nextRing = 1.5
		nextFreqMod = 4.5
	} else if number < 66 {
		nextRing = 3.0
		nextFreqMod = 9.0
	} else {
		nextRing = 0.75
		nextFreqMod = 1.5
	}

	return float32(nextRing), float32(nextFreqMod)
}

func getNextTiming() int {
	number := rand.Intn(100) + 1

	if number < 25 {
		return 32
	} else if number < 50 {
		return 16
	} else if number < 90 {
		return 8
	} else {
		return 4
	}
}

func getNextAmp() float32 {
	number := rand.Intn(100) + 1

	if number < 70 {
		return 440.0
	} else if number < 85 {
		return 660.0
	} else {
		return 220.0
	}
}

func sendTracRingFreq(client osc.Client) {
	nextNote := rand.Intn(len(maj_c))
	nextTiming := 1
	for true {
		nextRing, nextFreqMod := getNextRingAndFreqMod()

		arg := maj_c[nextNote]
		msg := osc.NewMessage("/juce/tracheaFrequency")
		msg.Append(float32(arg))
		client.Send(msg)
		nextNote = getNextNote(nextNote, len(maj_c))
		nextTiming = getNextTiming()

		msg = osc.NewMessage("/juce/ringFrequency")
		msg.Append(float32(nextRing))
		client.Send(msg)

		msg = osc.NewMessage("/juce/freqModFrequency")
		msg.Append(float32(nextFreqMod))
		client.Send(msg)

		time.Sleep(time.Second / time.Duration(nextTiming))
	}
}

func sendFreqAmp(client osc.Client) {
	for true {
		nextAmp := getNextAmp()
		msg := osc.NewMessage("/juce/freqModAmplitude")
		msg.Append(float32(nextAmp))
		client.Send(msg)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
}

func sendLevel(client osc.Client) {
	nextLevel := float32(0.0)
	for true {
		number := rand.Intn(100) + 1
		if number <= 75 {
			nextLevel = float32(0.0)
		} else {
			nextLevel = float32(0.5)
		}
		msg := osc.NewMessage("/juce/level")
		msg.Append(float32(nextLevel))
		client.Send(msg)
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	client := osc.NewClient("localhost", 9001)

	go sendTracRingFreq(*client)
	go sendFreqAmp(*client)
	go sendLevel(*client)

	for true {
	}
}
