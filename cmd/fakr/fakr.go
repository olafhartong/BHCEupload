package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/schollz/progressbar/v3"
)

// Super-duper quick 'n dirrty script to create large stubs to test the split tool on.

func main() {

	var count uint64
	flag.Uint64Var(&count, "count", 10_000, "number of items to generate.")

	var fileName string
	flag.StringVar(&fileName, "out", "fake-source.json", "the file the result is written to")
	flag.Parse()

	var progress bool
	flag.BoolVar(&progress, "no-progress", false, "hide the progress bar")

	fh, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("unable to create file %q", fileName)
		os.Exit(1)
	}

	fmt.Printf("Generating %d items, writing to %q... \n", count, fileName)

	// {data:[{..},{..}], meta:{..}}
	_, err = io.WriteString(fh, `{"data":[`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// TODO Writing directly to disk is slow, we can improve here
	enc := json.NewEncoder(fh)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	pb := progressbar.Default(int64(count), "generating")
	defer pb.Close()

	// 0 defaults to PRNG
	fkr := gofakeit.New(0)
	for i := count; i > 0; i-- {
		pb.Add(1)

		di := DataItem{
			Kind: fkr.BuzzWord(),
			Data: DataItemData{
				Members: nil,
				GroupId: fkr.UUID(),
			},
		}

		if fkr.Bool() {
			memberCount := fkr.Number(10, 100)
			var members = make([]DataItemMember, 0, memberCount)
			for j := memberCount; j > 0; j-- {
				members = append(members, DataItemMember{
					GroupId: fkr.UUID(),
					Member: DataItemMemberMember{
						Id:   fkr.UUID(),
						Text: fkr.HackerPhrase(),
					},
				})
			}

			di.Data.Members = &members
		}

		err := enc.Encode(di)
		if err != nil {
			fmt.Printf("Failed to encode %q", err.Error())
			continue
		}

		if i > 1 {
			_, _ = io.WriteString(fh, ",\n")
		}
	}

	io.WriteString(fh, `], "meta":`)
	enc.Encode(Meta{
		Type:    "azure",
		Version: 5,
		Count:   int(count),
	})
	io.WriteString(fh, `}`)
	fh.Close()

	fmt.Println("Thanks for choppin' by!")
}
