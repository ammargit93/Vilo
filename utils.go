type Node struct {
	hash      string
	data      interface{}
	commitMsg string
	next      *Node
	currTime  time.Time
	nodeName  string
}

func NewNode(hash string, data interface{}, commitMsg string, next *Node, currTime time.Time, nodeName string) Node {
	return Node{
		hash:      hash,
		data:      data,
		commitMsg: commitMsg,
		next:      next,
		currTime:  currTime,
		nodeName:  nodeName,
	}
}

func AddNode(head *Node, newnode Node) *Node {
	h := head
	for {
		if head.next == nil {
			head.next = &newnode
			break
		}
		head = head.next
		fmt.Println("loop")
	}
	head = h
	return head
}















package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {
	inputFile := "script.py"           // Python script file to compress and encrypt
	encryptedFile := "script.enc"      // Encrypted and compressed output
	decryptedFile := "script_decrypted.py" // Decrypted output file

	key := []byte("thisis32bytekeythisis32bytekey!") // Replace with a secure key
	EncryptAndCompress(inputFile, encryptedFile, key)
	fmt.Println("File successfully compressed and encrypted:", encryptedFile)
	DecryptAndDecompress(encryptedFile, decryptedFile, key)
	fmt.Println("File successfully decrypted and decompressed:", decryptedFile)
}

func EncryptAndCompress(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}
	if _, err := outputFile.Write(iv); err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	gzipWriter := gzip.NewWriter(cipher.StreamWriter{S: stream, W: outputFile})
	defer gzipWriter.Close()
	_, err = io.Copy(gzipWriter, inputFile)
	return err
}

func DecryptAndDecompress(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	gzipReader, err := gzip.NewReader(cipher.StreamReader{S: stream, R: inputFile})
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, gzipReader)
	return err
}
