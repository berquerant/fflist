// Code generated by "dataclass -type Entry -field Path string|Info fs.FileInfo -output entry_dataclass_generated.go"; DO NOT EDIT.

package walk

import "io/fs"

type Entry interface {
	Path() string
	Info() fs.FileInfo
}
type entry struct {
	path string
	info fs.FileInfo
}

func (s *entry) Path() string      { return s.path }
func (s *entry) Info() fs.FileInfo { return s.info }
func NewEntry(
	path string,
	info fs.FileInfo,
) Entry {
	return &entry{
		path: path,
		info: info,
	}
}