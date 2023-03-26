package system

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
)

// LoadFileRecords represents a file record in which each key is the first word of a new line,
// and the internal map represents key values joint by '=' separated by empty spaces.
//
//	 Example:
//	 Foo north=Bar west=Baz south=Qu-ux
//		Bar south=Foo west=Bee
//
//	 map[string]map[string]string{
//			"Foo": {"north": "Bar", "west": "Baz", "south": "Qu-ux"},
//			"Bar": {"south": "Foo", "west": "Bee"},
//		}
type LoadFileRecords = map[string]map[string]string

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrDuplicatedKey = errors.New("duplicated key")
)

// LoadFile must point to a file that has a valid LoadFileRecords format. If not it is going to return ErrInvalidFormat.
func (manager *Manager) LoadFile(path string) (LoadFileRecords, error) {
	output := make(LoadFileRecords)

	file, err := manager.OpenFunc(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Println("WARNING: failed to close file:", path)
		}
	}()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		key, values, err := processLoadFileRecord(scanner.Text())
		if err != nil {
			return nil, err
		}

		if _, alreadyExist := output[key]; alreadyExist {
			return nil, fmt.Errorf("%w: %q", ErrDuplicatedKey, key)
		}

		output[key] = values
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func processLoadFileRecord(s string) (string, map[string]string, error) {
	recs := make(map[string]string)

	line := strings.Split(s, " ")

	if len(line) < 1 {
		return "", nil, ErrInvalidFormat
	}

	key := line[0]

	// Index start at 1 since we already read the line header
	for i := 1; i < len(line); i++ {
		lineSplit := strings.Split(line[i], "=")
		if len(lineSplit) != 2 {
			return "", nil, ErrInvalidFormat
		}

		lineKey, lineValue := lineSplit[0], lineSplit[1]

		if _, alreadyExist := recs[lineKey]; alreadyExist {
			return "", nil, fmt.Errorf("%w: %q -> %q", ErrDuplicatedKey, key, lineKey)
		}

		recs[lineKey] = lineValue
	}

	return key, recs, nil
}
