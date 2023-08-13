package apis

import (
	"github.com/francoispqt/gojay"
)

func (pm *PaginationMarks) NKeys() int {
	return 0
}

func (pm *PaginationMarks) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {
	case "chunkNumber":
		err = dec.Int(&(pm.ReportNumber))

	case "isLastChunk":
		err = dec.Bool(&(pm.IsLastReport))
	}

	return err

}
