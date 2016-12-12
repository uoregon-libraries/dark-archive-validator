package checksum

import (
	"hash"
	"io"
	"os"
)

// Checksum combines a hash method with a block writing function for easily
// checksumming files by path, and customizing the hash or the read/write
// method used
type Checksum struct {
	Hash       hash.Hash
	BlockWrite func(path string, w io.Writer) error
}

// New returns a new Checksum using the default block write method, which just
// uses io.Copy to send a stream of bytes from the file into the hash
func New(h hash.Hash) *Checksum {
	return &Checksum{h, defaultBlockWrite}
}

// Sum resets the hash, runs the given path through BlockWrite, and returns the
// final hash Sum
func (f *Checksum) Sum(path string) ([]byte, error) {
	f.Hash.Reset()
	var err = f.BlockWrite(path, f.Hash)
	return f.Hash.Sum(nil), err
}

func defaultBlockWrite(path string, w io.Writer) error {
	var f, err = os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}
