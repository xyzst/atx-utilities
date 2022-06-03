package processor

import (
	"encoding/csv"
	"log"
	"os"
)

type CsvFile struct {
	input  *os.File
	output *os.File
	writer *csv.Writer
	reader *csv.Reader
}

func (receiver *CsvFile) Read() (record []string, err error) {
	return receiver.reader.Read()
}

func (receiver *CsvFile) Write(content []string) error {
	err := receiver.writer.Write(content)
	if err != nil {
		return err
	}

	defer receiver.writer.Flush()

	return nil
}

func (receiver *CsvFile) New(name string) interface{} {
	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	receiver.output = out
	receiver.writer = csv.NewWriter(receiver.output)

	if len(os.Args) != 2 {
		log.Fatal("did not provide path to input csv file")
	}

	input, err := os.Open(os.Args[1])
	receiver.input = input
	if err != nil {
		log.Fatal(err)
	}

	receiver.reader = csv.NewReader(input)
	receiver.reader.TrimLeadingSpace = true

	return receiver
}

func (receiver *CsvFile) Cleanup() {
	defer receiver.output.Close()
	defer receiver.input.Close()
}
