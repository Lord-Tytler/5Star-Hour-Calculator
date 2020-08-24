package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

var output [][]string

type student struct {
	lastName  string
	firstName string
	time      int64
}

var students []student

func main() {
	s := []string{"Last:", "First:", "Hours:"}
	output = append(output, s)
	files, err := ioutil.ReadDir("./logs")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(files); i++ {
		csvFile, err := os.Open("./logs/" + files[i].Name())
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
		toStruct(csvFile)
	}
	sortABC()
	toArray()
	toCSV(output)
}
func sortABC() {

	var names []string
	for i := 0; i < len(students); i++ {
		names = append(names, students[i].lastName)
	}
	sort.Strings(names)
	var a []student
	for i := 0; i < len(names); i++ {
		a = append(a, students[findName(names[i])])
	}
	students = a
}
func findName(last string) int {
	for i := 0; i < len(students); i++ {
		if students[i].lastName == last {
			return i
		}
	}
	return -1
}

func millisToHours(millis int64) float64 {
	return float64(millis) / 1000 / 60 / 60
}

func deltaT(arrived string, left string) int64 {
	layout := "1/2/2006 3:04:05 PM"
	in, err := time.Parse(layout, arrived)
	if err != nil {
		fmt.Println(err)
	}
	out, err1 := time.Parse(layout, left)
	if err1 != nil {
		fmt.Println(err)
	}
	return out.Sub(in).Milliseconds()
}
func reportError(datetime string, last string, first string) {
	fmt.Printf("\nData Error: %s, %s, %s", datetime, last, first)
}
func entryPresent(last string, first string) (bool, int) {
	name := first + last
	for i := 0; i < len(students); i++ {
		if students[i].firstName+students[i].lastName == name {
			return true, i
		}
	}
	return false, -1
}
func toArray() {
	for i := 0; i < len(students); i++ {
		last := students[i].lastName
		first := students[i].firstName
		hours := fmt.Sprintf("%f", millisToHours(students[i].time))
		s := []string{last, first, hours}
		output = append(output, s)
	}
}

func toCSV(s [][]string) {
	file, err := os.Create("./output.csv")
	if err != nil {
		fmt.Println("Error while creating the file::", err)
		return
	}
	writer := csv.NewWriter(file)
	err = writer.WriteAll(s)
	if err != nil {
		fmt.Println("Error while writing to the file ::", err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println("Error while closing the file ::", err)
		return
	}

}

func toStruct(csvFile *os.File) {

	r := csv.NewReader(csvFile)

	for i := 0; i < 5; i++ {
		_, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for record[1] == "" {
			record, err = r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		if present, index := entryPresent(record[1], record[2]); present == true {
			if deltaT(record[4], record[5]) <= 0 {
				reportError(record[4], record[1], record[2])
			} else {
				students[index].time += deltaT(record[4], record[5])
			}
		} else {
			students = append(students, student{record[1], record[2], deltaT(record[4], record[5])})
		}
	}
}
