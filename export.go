// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package defluxio

import (
	"bufio"
	"fmt"
	"os"
)

type TsvExporter struct {
	Path string
}

func NewTsvExporter(path string) (*TsvExporter, error) {
	retval := new(TsvExporter)
	retval.Path = path
	return retval, nil
}

func (tsve *TsvExporter) ExportDataset(timeReadings []MeterReading) error {
	file, err := os.Create(tsve.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, "timestamp\treading")
	for _, element := range timeReadings {
		exportLine := fmt.Sprintf("%d\t%f", element.Reading.Timestamp.Unix(),
			element.Reading.Value)
		fmt.Fprintln(w, exportLine)
	}
	return w.Flush()
}
