package main

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"os/exec"
	"strings"
)

const (
	BeginTableToken = "<table"
	EndTableToken   = "</table"
)

//EixfUtil reads output of exiftool -listx from stdout and converts it to json format
type EixfUtil struct {
	ctx context.Context
}

//NewEixfUtil is an ExifUtil constructor
func NewEixfUtil(ctx context.Context) *EixfUtil {
	return &EixfUtil{ctx: ctx}
}

func (s *EixfUtil) readExiftool(out chan<- string) {
	log.Println("reading exif")
	cmd := exec.Command("exiftool", "-listx")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	done := make(chan int)

	go func() {
		for scanner.Scan() {
			msg := scanner.Text()
			select {
			case <-s.ctx.Done():
				log.Println("reading interrupted")
				done <- 1
				return
			case out <- msg:
			}

		}
		done <- 0
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	doneStatus := <-done

	if doneStatus == 1 {
		log.Println("killing exiftool process...")
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err = cmd.Wait()
		if err != nil {
			log.Fatalln("wait error", err)
		}
	}

	log.Println("done reading exif")
}

func (s *EixfUtil) convertToJson(in <-chan string, out chan<- string) {
	isBuildingTable := false
	var sb strings.Builder
	tags := ""

	for {
		select {
		case newLine := <-in:
			if strings.HasPrefix(newLine, BeginTableToken) {
				isBuildingTable = true
			} else if strings.HasPrefix(newLine, EndTableToken) {
				sb.WriteString(newLine)
				isBuildingTable = false

				tableStr := sb.String()
				table, err := stringToTable(tableStr)
				if err != nil {
					log.Fatalln(err)
				}

				tagContainer := tableToTags(table)
				tagBytes, err := json.MarshalIndent(tagContainer, "", "  ")
				if err != nil {
					log.Fatalln(err)
				}

				tags = string(tagBytes)

			}

			if isBuildingTable {
				sb.WriteString(newLine)
			}
		case out <- tags:
		}
	}
}

func stringToTable(data string) (*Table, error) {
	table := &Table{}
	err := xml.Unmarshal([]byte(data), table)
	return table, err
}

func tableToTags(table *Table) *TagContainer {
	tagContainer := NewTagContainer(make([]Tag, 0))

	for _, tableTag := range table.Tag {
		tag := Tag{
			Writable:    tableTag.Writable,
			Path:        strings.Join([]string{table.Name, tableTag.Name}, ":"),
			Group:       table.Name,
			Description: map[string]string{},
			Type:        tableTag.Type,
		}

		for _, description := range tableTag.Desc {
			tag.Description[description.Lang] = description.Text
		}

		tagContainer.Tags = append(tagContainer.Tags, tag)
	}

	return tagContainer
}
