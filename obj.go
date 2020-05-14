package simplify

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func LoadOBJ(path string) (*Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return importOBJ(file)
}

func firstWords(value string, count int) (string, string) {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count--
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i], value[i+1:]
			}
		}
	}
	// Return the entire string.
	return value, ""
}

func strToVector(str string) (*Vector, error) {
	components := strings.Split(str, " ")

	if len(components) != 3 {
		return nil, errors.New("unable to parse: " + str)
	}

	xParse, err := strconv.ParseFloat(components[0], 64)
	if err != nil {
		return nil, errors.New("unable to parse X componenent: " + components[0])
	}

	yParse, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return nil, errors.New("unable to parse Y componenent: " + components[1])
	}

	zParse, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return nil, errors.New("unable to parse Z componenent: " + components[2])
	}

	v := Vector{xParse, yParse, zParse}
	return &v, nil
}

func strToFaceIndexes(str string) (int, int, int, error) {
	components := strings.Split(str, " ")

	if len(components) != 3 {
		return -1, -1, -1, fmt.Errorf("unable to parse: (%s)", str)
	}

	v1Components := strings.Split(components[0], "/")
	v1Parse, err := strconv.Atoi(v1Components[0])
	if err != nil {
		return -1, -1, -1, errors.New("unable to parse X componenent: " + v1Components[0])
	}

	v2Components := strings.Split(components[1], "/")
	v2Parse, err := strconv.Atoi(v2Components[0])
	if err != nil {
		return -1, -1, -1, errors.New("unable to parse Y componenent: " + v2Components[1])
	}

	v3Components := strings.Split(components[2], "/")
	v3Parse, err := strconv.Atoi(v3Components[0])
	if err != nil {
		return -1, -1, -1, errors.New("unable to parse Z componenent: " + v3Components[0])
	}

	return v1Parse, v2Parse, v3Parse, nil
}

func importOBJ(objStream io.Reader) (*Mesh, error) {
	if objStream == nil {
		return nil, errors.New("Need a reader to read obj from")
	}

	vertices := make([]Vector, 0)

	faces := make([]*Triangle, 0)

	scanner := bufio.NewScanner(objStream)
	for scanner.Scan() {
		ln := scanner.Text()
		firstWord, rest := firstWords(ln, 1)
		if firstWord == "v" {
			vector, err := strToVector(strings.TrimSpace(rest))
			if err != nil {
				return nil, err
			}
			vertices = append(vertices, *vector)
		}

		if firstWord == "f" {
			v1, v2, v3, err := strToFaceIndexes(strings.TrimSpace(rest))
			if err != nil {
				return nil, err
			}

			faces = append(faces, NewTriangle(vertices[v1-1], vertices[v2-1], vertices[v3-1]))
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return NewMesh(faces), nil
}
