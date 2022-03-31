package protoutil

import (
	"fmt"
	"testing"

	v1 "github.com/mmbednarek/fragma/api/fragma/core/v1"
)

func TestExtractValueByFieldName(t *testing.T) {
	app := v1.Volume{
		Status: &v1.VolumeStatus{
			Size:      150,
			FreeSpace: 0,
		},
	}

	something := ExtractValueByFieldName[int64](&app, "status.size")
	fmt.Println(something)
}
