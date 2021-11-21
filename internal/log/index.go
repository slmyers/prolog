package log

import (
	"fmt"
	"github.com/tysonmote/gommap"
	"io"
	"os"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
	maxSize uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
		maxSize: c.Segment.MaxIndexBytes,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	// TODO: this might have to be changed to count the entries instead of relying on file size
	// due to the truncation bug
	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_WRITE|gommap.PROT_READ,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}
// Close is there going to be some nasty ass bug from the failing of files to truncate? Will it go away when
// not being on my windoze?
func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return fmt.Errorf("error syncing mmap: %w", err)
	}
	if err := i.file.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %w", err)
	}
	// on wsl if the mmap is not unmaped then the truncation will fail.
	i.mmap.UnsafeUnmap()
	if err := i.file.Truncate(int64(i.size)); err != nil {
		// if max size has been reached or exceeded then truncation is expected to fail
		fmt.Printf("index size (%d), index max size (%d), error truncating file: (%s)", i.size, i.maxSize, err.Error())
	}
	return i.file.Close()
}

func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if in == -1 {
		out = uint32((i.size / entWidth) - 1)
	} else {
		out = uint32(in)
	}
	pos = uint64(out) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += entWidth
	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}