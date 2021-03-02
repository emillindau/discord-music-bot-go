package game

import (
	"io"
	"os"

	"github.com/jonas747/dca"
)

func downloadFile(name string, url string) error {
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

	err := convertFromURLtoDCA(url, name)
	return err
}

func convertFromURLtoDCA(url string, name string) error {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	encodingSession, err := dca.EncodeFile(url, options)

	if err != nil {
    return err
	}
	defer encodingSession.Cleanup()

	output, err := os.Create("temp/" + name + ".dca")

	if err != nil {
		return err
	}

	os.Remove("temp/" + name + ".mp3")

	_, err = io.Copy(output, encodingSession)
	return err
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