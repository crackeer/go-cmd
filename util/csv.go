package util

import (
	"encoding/csv"
	"io"
	"os"
)

// ReadCSVWithHeader
//  @param fileName
//  @return []map[string]interface{}
//  @return error
func ReadCSVWithHeader(fileName string) ([]map[string]interface{}, error) {

	list, err := ReadCSV(fileName)

	if err != nil {
		return nil, err
	}

	keys := list[0]
	retData := []map[string]interface{}{}

	for _, item := range list[1:] {
		tmp := map[string]interface{}{}
		for i, v := range item {
			tmp[keys[i]] = v
		}
		retData = append(retData, tmp)
	}
	return retData, nil

}

// ReadCSV
//  @param fileName
//  @return [][]string
//  @return error
func ReadCSV(fileName string) ([][]string, error) {

	list := [][]string{}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		list = append(list, record)
	}
	return list, nil
}
