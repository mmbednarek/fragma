package terminal

import (
	"os"
)

type FileTerminal struct {
	output *os.File
}

func (t *FileTerminal) WriteColor() {

}
