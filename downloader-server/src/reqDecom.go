package src

import (
	"errors"
	"strings"
)

type RData struct {
  Command string
  URL string
}

func RDataWrapper(rawData string) (RData, error) {
  // Filtering the requested data.
  dataArr := strings.Split(rawData, " ")

  if len(dataArr) != 2 {
    return RData{},errors.New("Invalid Requested Format: '<test | lock> <URL>'")
  }

  return RData{Command: dataArr[0], URL: dataArr[1]}, nil
}


