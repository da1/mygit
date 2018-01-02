package main

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type BlobObject struct {
	Size int
	Data string
}

func parseBlob(data string) (*BlobObject, error) {
	tokens := strings.SplitN(data, " ", 2)

	blob := tokens[0]
	if blob != "blob" {
		return nil, fmt.Errorf("not blob %s", blob)
	}

	data1 := tokens[1]

	nullChar := string([]byte{0})
	tokens = strings.Split(data1, nullChar)

	contentSize, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}
	data2 := tokens[1]

	return &BlobObject{contentSize, data2}, nil
}

func extract(zr io.Reader) (io.Reader, error) {
	return zlib.NewReader(zr)
}

func main() {
	p := flag.String("p", "", "p")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()

	filePath := fmt.Sprintf(".git/objects/%s/%s", (*p)[0:2], (*p)[2:])
	if *debug {
		fmt.Printf("file path %s\n", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("file can't open %s %v\n", filePath, err)
		return
	}

	r, err := extract(f)
	if err != nil {
		fmt.Printf("extract err %v\n", err)
		return
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	if *debug {
		fmt.Printf("read %s\n", b)
	}

	blobObject, err := parseBlob(string(b))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", blobObject.Data)
}
