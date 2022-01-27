//package main
//
//import (
//	"encoding/csv"
//	"fmt"
//	"os"
//)
//
////var driversdata [][]string
////
////
//func ReadCSV(file string) [][]string {
//	f, err := os.Open(file)
//	if err != nil {
//		fmt.Sprintf("Unable to open CSV File, %s : E: %v", file, err)
//	}
//	defer f.Close()
//
//	csvReader := csv.NewReader(f)
//	records, err := csvReader.ReadAll()
//
//	if err != nil {
//		fmt.Sprintf("Unable to open CSV File, %s : E: %v", file, err)
//	}
//	return records
//}
//
//func main()  {
//
//	driversdata := ReadCSV("./drivers.csv")
//
//	driver4 := driversdata[12]
//	fmt.Println(driver4)
//}
