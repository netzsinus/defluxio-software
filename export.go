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

type DailyTsvExporter struct {
	BasePath string
}

func NewDailyTsvExporter(basepath string) (retval *DailyTsvExporter, err error) {
	retval = new(DailyTsvExporter)
	retval.BasePath = basepath
	// Check if the base path really exists
	if _, err := os.Stat(retval.BasePath); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("Base path doesn't exist, %s", err)
		} else {
			err = fmt.Errorf("Base path stat error: %s", err)
		}
	}
}

func (dte *DailyTsvExporter) ExportDataset(timeReadings []MeterReading) (err error) {
	//TODO
}
