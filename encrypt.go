package main

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

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

// func DeleteJSONContent() {
// 	f, err := os.OpenFile(".vilo/stage.json", os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		fmt.Println("Error deleting content of the File")
// 	}
// 	f.Close()
// }

// func CreateFile(file string) {
// 	if _, err := os.Stat(file); err != nil {
// 		if os.IsNotExist(err) {
// 			_, err = os.Create(file)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		} else {
// 			fmt.Println(file + "already exists")
// 		}
// 	}
// }

// func SendFile(url, filePath string) error {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	var requestBody bytes.Buffer
// 	writer := multipart.NewWriter(&requestBody)

// 	part, err := writer.CreateFormFile("file", filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create form file: %v", err)
// 	}

// 	_, err = io.Copy(part, file)
// 	if err != nil {
// 		return fmt.Errorf("failed to copy file data: %v", err)
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		return fmt.Errorf("failed to close writer: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", url, &requestBody)
// 	if err != nil {
// 		return fmt.Errorf("failed to create HTTP request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to send HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		return fmt.Errorf("server error: %d, response: %s", resp.StatusCode, string(bodyBytes))
// 	}

// 	fmt.Println("File uploaded successfully!")
// 	return nil
// }
