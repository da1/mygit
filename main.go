package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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

func saveBlobObject(blob BlobObject) error {
	data := fmt.Sprintf("blob %d%s%s", blob.Size, string([]byte{0}), blob.Data)

	h := sha1.New()
	io.WriteString(h, data)
	sha1 := h.Sum(nil)
	fmt.Printf("blob %v, sha1: %s", blob, sha1)

	r := strings.NewReader(blob.Data)
	buf, err := compress(r)
	if err != nil {
		return fmt.Errorf("compress error %v", err)
	}

	filePath := fmt.Sprintf(".git/objects/%s/%s", sha1[0:2], sha1[2:])
	fmt.Printf("filePath %s\n", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create error %v", err)
	}
	defer file.Close()

	file.Write(buf.Bytes())
	return nil
}

func compress(r io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	defer zw.Close()

	if _, err := io.Copy(zw, r); err != nil {
		return buf, err
	}
	return buf, nil
}

func extract(zr io.Reader) (io.Reader, error) {
	return zlib.NewReader(zr)
}

func catFile(hash string, debug bool) {
	filePath := fmt.Sprintf(".git/objects/%s/%s", hash[0:2], hash[2:])
	if debug {
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

	if debug {
		fmt.Printf("read %s\n", b)
	}

	object := parseObject(string(b))

	blobObject, err := parseBlob(object)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if debug {
		fmt.Printf("%v\n", blobObject)
	}

	fmt.Printf("%s", blobObject.Data)
}

func addIndex(fileName string, debug bool) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("file can't open %s %v\n", fileName, err)
		return
	}
	defer f.Close()

	bufferSize := 1024
	data := ""
	buf := make([]byte, bufferSize)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		data += string(buf[:n])
	}

	if debug {
		fmt.Println(data)
	}

	blob := BlobObject{utf8.RuneCountInString(data), data}
	if debug {
		fmt.Printf("%v\n", blob)
	}
}

func main() {
	p := flag.String("p", "", "p")
	add := flag.String("add", "", "add file")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()

	if *p != "" {
		catFile(*p, *debug)
	}

	if *add != "" {
		addIndex(*add, *debug)
	}
}
