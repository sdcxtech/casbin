package load

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/sdcxtech/casbin/core"
)

// NewCSVIterator constructs an iterator that load CSV data from an io.Reader.
func NewCSVIterator(ioReader io.Reader) (itr core.AssertionIterator) {
	reader := csv.NewReader(ioReader)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	itr = &csvIterator{
		reader: reader,
	}

	return
}

type csvIterator struct {
	reader *csv.Reader
	err    error
}

func (it *csvIterator) Next() (ok bool, key string, vals []string) {
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

	return
}

func (it *csvIterator) Error() (err error) {
	return it.err
}
