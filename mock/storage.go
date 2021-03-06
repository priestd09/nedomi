package mock

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/ironsmile/nedomi/types"
)

//!TODO: Actually, looking at this, it is not a very good mock storage but it's
// an ok simple memory storage... maybe move this and get a real mock storage?

//!TODO: Make this storage thread safe; implementation details:
// - local per-key mutexes where possible? only for parts manipulation?
// - multiple maps as buckets to limit interference due to global locking:
//     - get the number of maps from config
//     - each map has its own rwmutex
//     - obj.hash mod N to determine which map is responsible for an object

// Storage implements the storage interface and is used for testing
type Storage struct {
	types.SyncLogger
	partSize uint64
	Objects  map[types.ObjectIDHash]*types.ObjectMetadata
	Parts    map[types.ObjectIDHash]map[uint32][]byte
}

// PartSize the maximum part size for the disk storage.
func (s *Storage) PartSize() uint64 {
	return s.partSize
}

// GetMetadata returns the metadata for this object, if present.
func (s *Storage) GetMetadata(id *types.ObjectID) (*types.ObjectMetadata, error) {
	if obj, ok := s.Objects[id.Hash()]; ok {
		return obj, nil
	}
	return nil, os.ErrNotExist
}

// GetPart returns an io.ReadCloser that will read the specified part of the object.
func (s *Storage) GetPart(idx *types.ObjectIndex) (io.ReadCloser, error) {
	if obj, ok := s.Parts[idx.ObjID.Hash()]; ok {
		if part, ok := obj[idx.Part]; ok {
			return ioutil.NopCloser(bytes.NewReader(part)), nil
		}
	}
	return nil, os.ErrNotExist
}

// GetAvailableParts returns an io.ReadCloser that will read the specified part of the object.
func (s *Storage) GetAvailableParts(oid *types.ObjectID) ([]*types.ObjectIndex, error) {
	var result = make([]*types.ObjectIndex, 0, len(s.Parts))
	if obj, ok := s.Parts[oid.Hash()]; ok {
		for partNum := range obj {
			result = append(result,
				&types.ObjectIndex{
					ObjID: oid,
					Part:  partNum,
				})
		}
	}
	return result, os.ErrNotExist
}

// SaveMetadata saves the supplied metadata.
func (s *Storage) SaveMetadata(m *types.ObjectMetadata) error {
	if _, ok := s.Objects[m.ID.Hash()]; ok {
		return os.ErrExist
	}

	s.Objects[m.ID.Hash()] = m

	return nil
}

// SavePart saves the contents of the supplied object part.
func (s *Storage) SavePart(idx *types.ObjectIndex, data io.Reader) error {
	objHash := idx.ObjID.Hash()
	if _, ok := s.Objects[objHash]; !ok {
		return errors.New("Object metadata is not present")
	}

	if _, ok := s.Parts[objHash]; !ok {
		s.Parts[objHash] = make(map[uint32][]byte)
	}

	if _, ok := s.Parts[objHash][idx.Part]; ok {
		return os.ErrExist
	}

	contents, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	s.Parts[objHash][idx.Part] = contents
	return nil
}

// Discard removes the object and its metadata.
func (s *Storage) Discard(id *types.ObjectID) error {
	if _, ok := s.Objects[id.Hash()]; !ok {
		return os.ErrNotExist
	}
	delete(s.Objects, id.Hash())
	delete(s.Parts, id.Hash())

	return nil
}

// DiscardPart removes the specified part of the object.
func (s *Storage) DiscardPart(idx *types.ObjectIndex) error {
	if obj, ok := s.Parts[idx.ObjID.Hash()]; ok {
		delete(obj, idx.Part)
		return nil
	}
	return os.ErrNotExist
}

// Iterate iterates over all the objects and passes them to the supplied callback
// function. If the callback function returns false, the iteration stops.
func (s *Storage) Iterate(callback func(*types.ObjectMetadata, ...*types.ObjectIndex) bool) error {
	for _, obj := range s.Objects {
		parts, _ := s.GetAvailableParts(obj.ID)
		if !callback(obj, parts...) {
			return nil
		}
	}
	return nil
}

// NewStorage returns a new mock storage that ready for use.
func NewStorage(partSize uint64) *Storage {
	return &Storage{
		partSize: partSize,
		Objects:  make(map[types.ObjectIDHash]*types.ObjectMetadata),
		Parts:    make(map[types.ObjectIDHash]map[uint32][]byte),
	}
}

// ChangeConfig does nothing
func (s *Storage) ChangeConfig(_ types.Logger) {
}
