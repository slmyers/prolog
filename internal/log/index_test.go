package log

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	idx, err := newIndex(f, c)
	require.NoError(t, err)
	_, _, err = idx.Read(-1)
	require.Error(t, err)
	require.Equal(t, f.Name(), idx.Name())

	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
	}

	for _, want := range entries {
		err = idx.Write(want.Off, want.Pos)
		require.NoError(t, err)

		_, pos, err := idx.Read(int64(want.Off))
		require.NoError(t, err)
		require.Equal(t, want.Pos, pos)
	}

	// index and scanner should error when reading past existing entries
	_, _, err = idx.Read(int64(len(entries)))
	require.Equal(t, io.EOF, err)
	err = idx.Close()
	require.NoError(t, err)

	// index should build its state from the existing file
	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	idx, err = newIndex(f, c)
	require.NoError(t, err)
	off, pos, err := idx.Read(-1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), off)
	require.Equal(t, entries[1].Pos, pos)
}

//func TestIndex(t *testing.T) {
//	c := Config{}
//	c.Segment.MaxIndexBytes = 1024
//
//	t.Run("sanity", func(t *testing.T) {
//		f, err := ioutil.TempFile(os.TempDir(), "index_test")
//		require.NoError(t, err)
//		defer os.Remove(f.Name())
//		idx, err := newIndex(f, c)
//		require.NoError(t, err)
//		_, _, err = idx.Read(-1)
//		require.Error(t, err)
//		require.Equal(t, f.Name(), idx.Name())
//		idx.Close()
//	})
//
//	t.Run("multiple read-writes", func(t *testing.T) {
//		f, err := ioutil.TempFile(os.TempDir(), "index_test")
//		require.NoError(t, err)
//		defer os.Remove(f.Name())
//		idx, err := newIndex(f, c)
//
//		entries := []struct {
//			Off uint32
//			Pos uint64
//		}{
//			{Off: 0, Pos: 0},
//			{Off: 1, Pos: 10},
//		}
//
//		for _, want := range entries {
//			err = idx.Write(want.Off, want.Pos)
//			require.NoError(t, err)
//
//			_, pos, err := idx.Read(int64(want.Off))
//			require.NoError(t, err)
//			require.Equal(t, want.Pos, pos)
//		}
//		idx.Close()
//	})
//
//	t.Run("index and scanner should error when reading past existing entries", func(t *testing.T) {
//		entries := []struct {
//			Off uint32
//			Pos uint64
//		}{
//			{Off: 0, Pos: 0},
//			{Off: 1, Pos: 10},
//		}
//		f, err := ioutil.TempFile(os.TempDir(), "index_test_sanity")
//		require.NoError(t, err)
//		defer os.Remove(f.Name())
//		idx, err := newIndex(f, c)
//		// index and scanner should error when reading past existing entries
//		_, _, err = idx.Read(int64(len(entries)))
//		require.Equal(t, io.EOF, err)
//		idx.Close()
//	})
//
//	t.Run("index should build its state from the existing file", func(t *testing.T) {
//		entries := []struct {
//			Off uint32
//			Pos uint64
//		}{
//			{Off: 0, Pos: 0},
//			{Off: 1, Pos: 10},
//		}
//
//		f, err := ioutil.TempFile(os.TempDir(), "index_test")
//		require.NoError(t, err)
//		defer os.Remove(f.Name())
//
//		idx, err := newIndex(f, c)
//		require.NoError(t, err)
//
//		for _, want := range entries {
//			err = idx.Write(want.Off, want.Pos)
//
//			require.NoError(t, err)
//		}
//		idx.Close()
//
//		// index should build its state from the existing file
//		f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
//		fmt.Printf("%+v\n", f)
//		idx, err = newIndex(f, c)
//		require.NoError(t, err)
//		off, pos, err := idx.Read(-1)
//		require.NoError(t, err)
//		require.Equal(t, entries[1].Pos, pos)
//		require.Equal(t, uint32(1), off)
//
//	})
//
//
//}
