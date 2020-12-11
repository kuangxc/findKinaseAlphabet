package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	defaultFileName = "kinase_human_domain.fasta"

	defaultRefTag = ">SRC_Hs_Src_SrcA"

	paramFile = "param.txt"
)

var (
	fileName string = ""
	refTag   string = ""
)

func readFile() map[string]string {
	data := make(map[string]string)
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("read file err:", err)
		return nil
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var key, value string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Contains(line, ">") {
			if len(key) > 0 && len(value) > 0 {
				data[key] = value
				key = line
				value = ""
			} else {
				key = line
			}
		} else {
			value += line
		}
		if err != nil {
			if err == io.EOF {
				data[key] = value
			}
			break
		}
	}
	//fmt.Println(data)
	return data
}

func getTagPrefix(s string) string {
	index := strings.Index(s, "_")
	return s[1:index]
}

func countAlphabet(data map[string]string, refAlphabet rune) {
	file, err := os.Create(fmt.Sprintf("result_%s.txt", string(refAlphabet)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	refLine := data[refTag]
	index := 0
	for i, a := range refLine {
		//fmt.Println(refLine)
		if a >= 'A' && a <= 'Z' {
			index++
			_, err := file.Write([]byte(fmt.Sprintf("SRC=%d\n", index)))
			if err != nil {
				fmt.Println(err)
				return
			}
			alphabetCount := 0
			sameAlphabetCount := 0
			sameAlphabetTags := ""
			for tag, line := range data {
				if tag != refTag {
					if line[i] >= 'A' && line[i] <= 'Z' {
						alphabetCount++
						if line[i] == byte(refAlphabet) {
							sameAlphabetCount++
							sameAlphabetTags += getTagPrefix(tag) + "|"
						}
					}
				}
			}
			file.Write([]byte(fmt.Sprintf("%s=%d\n", string(refAlphabet), sameAlphabetCount)))
			file.Write([]byte(fmt.Sprintf("res=%d\n", alphabetCount)))
			file.Write([]byte(fmt.Sprintf("%s\n", sameAlphabetTags)))
			file.Write([]byte(fmt.Sprintln()))
		}
	}
}

func readParamFromFile() {
	f, err := os.Open(paramFile)
	if err != nil {
		fmt.Println("no param file,use default:", err)
		fileName = defaultFileName
		refTag = defaultRefTag
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	line, err := buf.ReadString('\n')
	if err != nil {
		fmt.Println("read param err ", err)
	}
	fileName = strings.TrimSpace(line)

	line, err = buf.ReadString('\n')
	if err != nil {
		fmt.Println("read param ", err)
	}
	refTag = ">" + strings.TrimSpace(line)
}

func main() {
	// 1.read param from file
	readParamFromFile()
	// 2. read file,create map [key=tag] value=lines info
	data := readFile()
	// 3. count every alphabet and write file
	for a := 'A'; a <= 'Z'; a++ {
		countAlphabet(data, a)
	}

}
