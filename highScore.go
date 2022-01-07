package main

// import (
// 	"encoding/gob"
// 	"fmt"
// 	"os"
// )

// const (
// 	HighsScoreFilePath = "highscore.gob"
// )

// // type HighScore struct {
// // 	score int
// // }

// func SaveHighScore(hs int) {
// 	file, err := os.Open(HighsScoreFilePath)
// 	if err != nil {
// 		panic(fmt.Errorf("cant open file %v: %v", HighsScoreFilePath, err))
// 	}
// 	defer file.Close()

// 	file.WriteString(fmt.Sprintf("%v", hs))

// 	// encoder := gob.NewEncoder(file)
// 	// fmt.Println(hs)
// 	// encoder.Encode(hs)
// }

// func LoadHighScore() int {
// 	// hs := HighScore{}

// 	file, err := os.Open(HighsScoreFilePath)
// 	if err != nil {
// 		panic(fmt.Errorf("cant open file %v: %v", HighsScoreFilePath, err))
// 	}
// 	defer file.Close()

	

// 	// decoder := gob.NewDecoder(file)
// 	// decoder.Decode(&hs)
// 	// fmt.Println()
// 	// return hs
// }
