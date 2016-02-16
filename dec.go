package main

import (
	"bufio"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const ()

type PEM struct {
	raw     string
	scanner *bufio.Scanner
	Headers map[string]string
}

func (p *PEM) getLine() (string, error) {
	if p.scanner.Scan() {
		line := p.scanner.Text()
		if err := p.scanner.Err(); err != nil {
			return "", err
		}
		return line, nil
	}
	return "", ERRGetLine
}

func (p *PEM) NewFromString(s string) {
	p.scanner = bufio.NewScanner(strings.NewReader(s))
}

func (p *PEM) BlockType() (string, error) {
	blkLine, err := p.getLine()
	if err != nil {
		return "", err
	}

	blockType, err := getBlockType([]byte(blkLine))

	if err != nil {
		return "", err

	}

	return blockType, nil

}

func (p *PEM) ExtractHeaders() (map[string]string, error) {

	var hdrs map[string]string = map[string]string{}

	for {
		hdrLine, err := p.getLine()
		if err != nil {
			return map[string]string{}, err
		}

		if hdrLine == "" {
			break
		}

		// Split the string
		keyVal := strings.Split(hdrLine, ":")
		if len(keyVal) == 2 {
			fmt.Printf("Adding key: [%s] => Val: [%s]\n", keyVal[0], keyVal[1])
			hdrs[strings.Trim(keyVal[0], " ")] = strings.Trim(keyVal[1], " ")
		}

	}

	// TODO = need to scan lines until you find a blank line
	return hdrs, nil
}

func (p *PEM) ExtractBody() (string, error) {

	var body string
	var hasEnd bool = false

	for p.scanner.Scan() {
		bline := p.scanner.Text()
		if err := p.scanner.Err(); err != nil {
			return "", err
		}
		// See if we've hit the last line
		re := regexp.MustCompile("^-----END ([a-zA-Z ]+)-----")
		matches := re.FindSubmatch([]byte(bline))

		// no match - add to body
		if len(matches) < 2 {
			body = body + bline
		} else {
			hasEnd = true
		}
	}

	if !hasEnd {
		return "", ERRNoBody
	}
	return body, nil

}

var (
	ERRNoBlock = errors.New("Error No Block Found")
	ERRGetLine = errors.New("getLine: out of data to scan")
	ERRNoBody  = errors.New("Failed to find end of body")
)

func main() {

	var p PEM
	secretMsg, err := ioutil.ReadFile("enc.pem")
	if err != nil {
		fmt.Printf("ReadFile: %s\n", err)
		os.Exit(1)
	}

	p.NewFromString(string(secretMsg))
	blockType, err := p.BlockType()
	if err != nil {
		fmt.Printf("failed to find blocktype")
		os.Exit(1)
	}

	headers, err := p.ExtractHeaders()
	if err != nil {
		fmt.Printf("failed to extract headers")
		os.Exit(1)
	}

	body, err := p.ExtractBody()
	if err != nil {
		fmt.Printf("failed to extract body")
		os.Exit(1)
	}

	fmt.Printf("body:\n%q\n", body)

	password := []byte("password")

	sDec, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Printf("secret is\n%X\n\n", sDec)

	var b pem.Block = pem.Block{
		Type:    blockType,
		Headers: headers,
		Bytes:   sDec,
	}

	DecryptPEMBlock, err := x509.DecryptPEMBlock(&b, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", DecryptPEMBlock)
}

func stripNewLines(s []byte) []byte {
	ns := make([]byte, len(s), len(s))

	j := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			ns[j] = s[i]
			j++
		}
	}
	return ns
}

func getBlockType(s []byte) (string, error) {

	re := regexp.MustCompile("^-----BEGIN ([a-zA-Z ]+)-----")
	matches := re.FindSubmatch(s)
	if len(matches) < 2 {
		fmt.Printf("err: %d\n", len(matches))
		return "", ERRNoBlock
	} else {
		return string(matches[1]), nil
	}

}
