// (C) 2014 Mathias Dalheimer <md@gonium.net>. See LICENSE file for
// license.
package defluxio

import (
	"bufio"
	"fmt"
	"os"
	"text/template"
	"time"
)

const (
	licenseTemplate = `# The https://netzsin.us grid frequency measurements are 
# (c) Mathias Dalheimer, <md@gonium.net>.
#
# This database is made available under the Open Database License:
#   http://opendatacommons.org/licenses/odbl/1.0/.
# Any rights in individual contents of the database are licensed under
# the Database Contents License:
#   http://opendatacommons.org/licenses/dbcl/1.0/
# Please see http://opendatacommons.org/licenses/odbl/summary/ for a
# human-readable explanation of the license.
#
# Generated on {{.GenerationDate}}
timestamp	reading
`
)

type LicenseData struct {
	GenerationDate string
}

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
	gd := time.Now().String()
	ld := LicenseData{GenerationDate: gd}
	t :=
		template.Must(template.New("licenseTemplate").Parse(licenseTemplate))
	err = t.Execute(w, ld)
	if err != nil {
		return fmt.Errorf("Failed to create license header: %s", err)
	}
	for _, element := range timeReadings {
		exportLine := fmt.Sprintf("%d\t%f", element.Reading.Timestamp.Unix(),
			element.Reading.Value)
		fmt.Fprintln(w, exportLine)
	}
	return w.Flush()
}
