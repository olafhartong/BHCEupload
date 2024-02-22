package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bcicen/jstream"
)

// Meta is the expected meta property of the resulting documents. This might need updating depending on the Blood Hound spec
type Meta struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Count   int    `json:"count,omitempty"`
	Methods int    `json:"methods,omitempty"`
}

type Config struct {
	SourceFileName string
	ChunkSize      uint
}

var (
	reMetaCapturePattern = regexp.MustCompile(`(?s)"meta":\s*({.*?})`)
)

func main() {
	config := Config{}

	flag.StringVar(&config.SourceFileName, "source", "", "The path of the source file (e.g.: export.json)")
	flag.UintVar(&config.ChunkSize, "chunk-size", 100, "Data Items to split on")
	flag.Parse()

	if config.SourceFileName == "" {
		fmt.Println("The -source flag cannot be empty")
		os.Exit(1)
	}

	fh, err := os.Open(config.SourceFileName)
	if err != nil {
		fmt.Printf("Error opening %s: %s\n", config.SourceFileName, err.Error())
		os.Exit(1)
	}

	metaSource, err := getMetaTagFromFile(fh)
	if err != nil {
		fmt.Printf("Error reading meta tag from the source: %s", err.Error())
		_ = fh.Close()
		os.Exit(1)
	}

	// Our working buffer. We start a bit larger, so that we allocate less often
	buf := &strings.Builder{}
	buf.Grow(2024)
	buf.Reset()

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")

	// Since we're manually constructing the JSON document, to save some memory and overhead, we define some wrappers here.
	finalizeDoc := func(result string, meta Meta) string {

		// we naively concatenate a comma after each array element to our buffer, here we chop the last one off,
		// if there is one. This way we only have one if statement and hopefully a lot less branch miss-predictions.
		if ll := len(result); result[ll-1:] == "," {
			result = result[:ll-1]
		}

		// Wrapping result in the expected structure: {data:[{..},{..}], meta:{..}}
		return `{"data":[` + result + fmt.Sprintf(
			`],"meta":{"type":"%s","version":%d,"count":%d}}`,
			meta.Type, meta.Version, meta.Count,
		)
	}

	// addData JSON re-encodes the type we read from the JSON decoder and adds it to the buffer with a trailing comma
	addData := func(buf *strings.Builder, enc *json.Encoder, d any) error {

		// TODO: We only care about file offset and range, converting to types is a lot of unnecessary overhead
		err := enc.Encode(d)

		if err != nil {
			return err
		}
		buf.WriteByte(',')
		return nil
	}

	// writeDocToFile creates a new file and writes the result to it
	writeDocToFile := func(runSeed int64, docCnt int64, res string) (string, error) {
		fh, err := os.Create(fmt.Sprintf("./res_%d_%04d.json", runSeed, docCnt))
		if err != nil {
			return "", err
		}

		_, err = io.WriteString(fh, res)
		if err != nil {
			return "", err
		}

		return fh.Name(), fh.Close()
	}

	i := uint(0)
	docCnt := int64(0)
	runSeed := time.Now().UnixMilli() // Only used for sortable grouping of the chunk files.

	items := jstream.NewDecoder(fh, 2)
	for mv := range items.Stream() {
		if mv.ValueType != jstream.Object {
			continue
		}

		err := addData(buf, enc, mv.Value)
		if err != nil {
			fmt.Printf("Error while marshalling %q", err.Error())
			continue
		}

		i++
		if i >= config.ChunkSize {

			result := finalizeDoc(buf.String(), Meta{
				Type:    metaSource.Type,
				Version: metaSource.Version,
				Count:   int(i),
			})
			buf.Reset()

			fileName, err := writeDocToFile(runSeed, docCnt, result)
			if err != nil {
				fmt.Printf("Error while writing document %d, %s", docCnt, err.Error())
				continue
			}

			fmt.Printf("Wrote file: %s\n", fileName)

			i = 0
			docCnt++
		}
	}

	// If we still have dangling items, we wrap them up in a final document here.
	if buf.Len() > 0 {
		result := finalizeDoc(buf.String(), Meta{
			Type:    metaSource.Type,
			Version: metaSource.Version,
			Count:   int(i),
		})
		buf.Reset()

		fileName, err := writeDocToFile(runSeed, docCnt, result)
		if err != nil {
			fmt.Printf("Error while writing document %d, %s", docCnt, err.Error())
			return
		}

		fmt.Printf("Wrote file: %s\n", fileName)
	}
}

// getMetaTag tries to capture the (expected) "meta" element from the source document. For efficiency's sake it acts on
// a small chunk of the input. Typically, the "meta" element is the last section of the source document.
func getMetaTag(chunk []byte) (Meta, error) {
	matches := reMetaCapturePattern.FindStringSubmatch(string(chunk))
	// 0 is always the entire match and 1 is the first capture group
	if len(matches) < 2 {
		return Meta{}, fmt.Errorf("could not find the meta tag in this chunk")
	} else if len(matches) > 2 {
		return Meta{}, fmt.Errorf("found multiple meta tags in this chunk, this is unsupported")
	}

	var meta Meta
	return meta, json.Unmarshal([]byte(matches[1]), &meta)
}

// getMetaTagFromFile accepts a file and returns the Meta type, or an error if it can't find one. It reset the cursor
// back at the beginning of the file.
func getMetaTagFromFile(fh *os.File) (meta Meta, err error) {
	if fh == nil {
		err = fmt.Errorf("no file handle specified")
		return
	}

	// The number of bytes we look back from the end of the file
	const chunkSize = 128

	// Placing the cursor near the end
	_, err = fh.Seek(-chunkSize, 2)
	if err != nil {
		return
	}

	// Make sure we are back at the beginning when we bail here.
	defer func() {
		_, err = fh.Seek(0, 0)
	}()

	// Reading the last n bytes of the file, this should typically contain the meta tag
	var end = make([]byte, chunkSize)
	_, err = io.ReadFull(fh, end)
	if err != nil {
		return
	}

	// Extract the meta tag, we'll be needing some of the data
	meta, err = getMetaTag(end)

	return
}
