package system

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_LoadFile(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		// Test with invalid file format
		manager := NewManager()

		_, err := manager.LoadFile("./testdata/invalid_file_format.txt")

		assert.Error(t, err)
	})

	t.Run("create new manager", func(t *testing.T) {
		manager := NewManager()

		assert.NotNil(t, manager)
	})

	t.Run("invalid format", func(t *testing.T) {
		// create a temporary file and write test data to it
		file, err := ioutil.TempFile("", "testfile.txt")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		testData := "Foo north=Bar west=Baz south=Qu-ux\nBar south=Foo west=Bee"
		_, err = file.WriteString(testData)
		require.NoError(t, err)

		// create a new Manager with default OpenFunc
		manager := NewManager()

		// call LoadFile with the temporary file path
		records, err := manager.LoadFile(file.Name())
		require.NoError(t, err)

		// check that records were parsed correctly
		expected := LoadFileRecords{
			"Foo": {
				"north": "Bar",
				"west":  "Baz",
				"south": "Qu-ux",
			},
			"Bar": {
				"south": "Foo",
				"west":  "Bee",
			},
		}
		require.Equal(t, expected, records)

		// test error handling: invalid format
		invalidData := "Foo north=Bar\nBar south=Foo=west=Bee"
		_, err = file.WriteString(invalidData)
		require.NoError(t, err)

		_, err = manager.LoadFile(file.Name())
		require.Error(t, err)
		require.EqualError(t, err, ErrInvalidFormat.Error())

		// test error handling: duplicated key
		duplicateData := "Foo north=Bar\nFoo west=Baz"
		_, err = file.WriteString(duplicateData)
		require.NoError(t, err)

		_, err = manager.LoadFile(file.Name())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidFormat)
	})
}
