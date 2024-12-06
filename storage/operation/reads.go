package operation

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/vmihailenco/msgpack"

	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/storage"
)

// CheckFunc is a function that checks if the value should be read and decoded.
// return (true, nil) to read the value and pass it to the CreateFunc and HandleFunc for decoding
// return (false, nil) to skip reading the value
// return (false, err) if running into any exception, the iteration should be stopped.
type CheckFunc func(key []byte) (bool, error)

// CreateFunc returns a pointer to an initialized entity that we can potentially
// decode the next value into during a badger DB iteration.
type CreateFunc func() interface{}

// HandleFunc is a function that starts the processing of the current key-value
// pair during a badger iteration. It should be called after the key was checked
// and the entity was decoded.
// No errors are expected during normal operation. Any errors will halt the iteration.
type HandleFunc func() error
type IterationFunc func() (CheckFunc, CreateFunc, HandleFunc)

// IterateKeysInPrefixRange will iterate over all entries in the database, where the key starts with a prefixes in
// the range [startPrefix, endPrefix] (both inclusive). We require that startPrefix <= endPrefix (otherwise this
// function errors). On every such key, the `check` function is called. If `check` errors, iteration is aborted.
// No errors expected during normal operations.
func IterateKeysInPrefixRange(startPrefix []byte, endPrefix []byte, check func(key []byte) error) func(storage.Reader) error {
	return Iterate(startPrefix, endPrefix, func() (CheckFunc, CreateFunc, HandleFunc) {
		return func(key []byte) (bool, error) {
			err := check(key)
			if err != nil {
				return false, err
			}
			return false, nil
		}, nil, nil
	}, storage.IteratorOption{IterateKeyOnly: true})
}

// Iterate will iterate over all keys with prefixes in the given range [startPrefix, endPrefix] (both inclusive)
func Iterate(startPrefix []byte, endPrefix []byte, iterFunc IterationFunc, opt storage.IteratorOption) func(storage.Reader) error {
	return func(r storage.Reader) error {

		if len(startPrefix) == 0 {
			return fmt.Errorf("startPrefix prefix is empty")
		}

		if len(endPrefix) == 0 {
			return fmt.Errorf("endPrefix prefix is empty")
		}

		// Reverse iteration is not supported by pebble
		if bytes.Compare(startPrefix, endPrefix) > 0 {
			return fmt.Errorf("startPrefix key must be less than or equal to endPrefix key")
		}

		it, err := r.NewIter(startPrefix, endPrefix, opt)
		if err != nil {
			return fmt.Errorf("can not create iterator: %w", err)
		}
		defer it.Close()

		for it.First(); it.Valid(); it.Next() {
			item := it.IterItem()
			key := item.Key()

			// initialize processing functions for iteration
			check, create, handle := iterFunc()

			keyCopy := make([]byte, len(key))

			// The underlying database may re-use and modify the backing memory of the returned key.
			// Tor safety we proactively make a copy before passing the key to the upper layer.
			copy(keyCopy, key)

			// check if we should process the item at all
			shouldReadValue, err := check(keyCopy)
			if err != nil {
				return err
			}
			if !shouldReadValue { // skip reading value
				continue
			}

			err = item.Value(func(val []byte) error {

				// decode into the entity
				entity := create()
				err = msgpack.Unmarshal(val, entity)
				if err != nil {
					return irrecoverable.NewExceptionf("could not decode entity: %w", err)
				}

				// process the entity
				err = handle()
				if err != nil {
					return fmt.Errorf("could not handle entity: %w", err)
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("could not process value: %w", err)
			}
		}

		return nil
	}
}

// Traverse will iterate over all keys with the given prefix
func Traverse(prefix []byte, iterFunc IterationFunc, opt storage.IteratorOption) func(storage.Reader) error {
	return Iterate(prefix, prefix, iterFunc, opt)
}

// Exists returns true if a key exists in the database.
// No errors are expected during normal operation.
func Exists(key []byte, keyExists *bool) func(storage.Reader) error {
	return func(r storage.Reader) error {
		_, closer, err := r.Get(key)
		if err != nil {
			// the key does not exist in the database
			if errors.Is(err, storage.ErrNotFound) {
				*keyExists = false
				return nil
			}
			// exception while checking for the key
			return irrecoverable.NewExceptionf("could not load data: %w", err)
		}
		defer closer.Close()

		// the key does exist in the database
		*keyExists = true
		return nil
	}
}

// Retrieve will retrieve the binary data under the given key from the database
// and decode it into the given entity. The provided entity needs to be a
// pointer to an initialized entity of the correct type.
// Error returns:
//   - storage.ErrNotFound if the key does not exist in the database
//   - generic error in case of unexpected failure from the database layer, or failure
//     to decode an existing database value
func Retrieve(key []byte, entity interface{}) func(storage.Reader) error {
	return func(r storage.Reader) error {
		val, closer, err := r.Get(key)
		if err != nil {
			return err
		}

		defer closer.Close()

		err = msgpack.Unmarshal(val, entity)
		if err != nil {
			return irrecoverable.NewExceptionf("could not decode entity: %w", err)
		}
		return nil
	}
}

// FindHighestAtOrBelow is for database entries that are indexed by block height. It is suitable to search
// keys with the format prefix` + `height` (where "+" denotes concatenation of binary strings). The height
// is encoded as Big-Endian (entries with numerically smaller height have lexicographically smaller key).
// The function finds the *highest* key with the given prefix and height equal to or below the given height.
func FindHighestAtOrBelow(prefix []byte, height uint64, entity interface{}) func(storage.Reader) error {
	return func(r storage.Reader) error {
		if len(prefix) == 0 {
			return fmt.Errorf("prefix must not be empty")
		}

		key := append(prefix, EncodeKeyPart(height)...)
		it, err := r.NewIter(prefix, key, storage.DefaultIteratorOptions())
		if err != nil {
			return fmt.Errorf("can not create iterator: %w", err)
		}
		defer it.Close()

		var highestKey []byte

		// find highest value below the given height
		for it.First(); it.Valid(); it.Next() {
			// copy the key to avoid the underlying slices of the key
			// being modified by the Next() call
			highestKey = it.IterItem().KeyCopy(highestKey)
		}

		if len(highestKey) == 0 {
			return storage.ErrNotFound
		}

		// read the value of the highest key
		val, closer, err := r.Get(highestKey)
		if err != nil {
			return err
		}

		defer closer.Close()

		err = msgpack.Unmarshal(val, entity)
		if err != nil {
			return irrecoverable.NewExceptionf("failed to decode value: %w", err)
		}

		return nil
	}
}
