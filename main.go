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

type Object struct {
	ObjectType string
	Data       string
}

type BlobObject struct {
	Size int
	Data string
}

func parseObject(data string) Object {
	tokens := strings.SplitN(data, " ", 2)
	return Object{tokens[0], tokens[1]}
}

func parseBlob(obj Object) (*BlobObject, error) {
	if obj.ObjectType != "blob" {
		return nil, fmt.Errorf("not blob %s", obj.ObjectType)
	}

	nullChar := string([]byte{0})
	tokens := strings.Split(obj.Data, nullChar)

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

	object := parseObject(string(b))

	blobObject, err := parseBlob(object)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", blobObject.Data)
}
