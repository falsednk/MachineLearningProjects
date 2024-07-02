package main

import (
	"fmt"
	"log"
	"os"

	"example/test/preprocessing"

	"github.com/go-gota/gota/dataframe"
)

func main() {
	file, err := os.Open("weatherAUS.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	fmt.Printf("%T\n", df)

	preprocessing.GetInfoDf(df)
}
