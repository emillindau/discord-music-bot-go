package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/jonas747/dca"
)

func DownloadFile(name string, url string) (string, error) {
	// resp, err := http.Get(url)

	// if err != nil {
	// 	return err
	// }

	// defer resp.Body.Close()

	// fileName := "temp/" + name + ".mp3"
	// // Create the file
	// out, err := os.Create(fileName)
	// if err != nil {
	// 	return err
	// }

	// defer out.Close()

	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	return err
	// }

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for _, s := range sg {
	// 		downloadFile(s.id, s.downloadUrl)
	// 	}
	// }()
	// wg.Wait()

	path, err := convertFromURLtoDCA(url, name)
	if err != nil {
		return "", err
	}

	return path, nil
}

func convertFromURLtoDCA(url string, name string) (string, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	encodingSession, err := dca.EncodeFile(url, options)

	filepath := "temp/" + name + ".dca"

	if err != nil {
    return "", err
	}
	defer encodingSession.Cleanup()

	output, err := os.Create(filepath)

	if err != nil {
		return "", err
	}

	_, err = io.Copy(output, encodingSession)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

// func convertMP3toDCA(filePath string, name string) error {
// 	encodeSession, err := dca.EncodeFile(filePath, dca.StdEncodeOptions)

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer encodeSession.Cleanup()

// 	output, err := os.Create("temp/" + name + ".dca")

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	// os.Remove("temp/" + name + ".mp3")

// 	io.Copy(output, encodeSession)
// 	return nil
// }


func LoadSound(filePath string) ([][]byte, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	var opuslen int16
	var buffer = make([][]byte, 0)
	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// eof
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from file");
			return nil, err
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)

		if err != nil {
			fmt.Println("Error reading pcm")
			return nil, err
		}

		buffer = append(buffer, inBuf)
	}
}