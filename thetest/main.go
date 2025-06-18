package main

import "fmt"

func main() {
  
}

func firstUniqChar(s string) int {
  saveIndex:= -1
      checkerMap := make(map[rune]int)
  for index, _ := range s{
  checkerMap[s[index]]= +1
  if checkerMap[s[index]] == 1 
  saveIndex = index
  }
  
  return saveIndex
  }