package load

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

// NewCSVIterator constructs an iterator that load CSV data from an io.Reader.
func NewCSVIterator(ioReader io.Reader) (itr *CSVIterator) {
	reader := csv.NewReader(ioReader)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	itr = &CSVIterator{
		reader: reader,
	}

	return
}

type CSVIterator struct {
	reader *csv.Reader
	err    error
}

func (it *CSVIterator) Next() (ok bool, key string, vals []string) {
	if it.err != nil {
		return
	}

	var record []string

	for {
		var err error
		record, err = it.reader.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			it.err = err

			return
		}

		if len(record) > 1 {
			break
		}
	}

	ok = true
	key = strings.TrimSpace(record[0])
	vals = record[1:]

	return ok, key, vals
}

func (it *CSVIterator) Error() (err error) {
	return it.err
}
