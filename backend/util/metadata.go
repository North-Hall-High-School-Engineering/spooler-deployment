package util

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type FileMetadata struct {
	StlFile   string `json:"stl_file"`
	Thumbnail string `json:"thumbnail"`
}

func GetFileMetadata(fileName string, fileHandle multipart.File) (FileMetadata, error) {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".stl":
		content, err := io.ReadAll(fileHandle)
		if err != nil {
			return FileMetadata{}, err
		}

		encoded := base64.StdEncoding.EncodeToString(content)

		return FileMetadata{StlFile: encoded}, nil

	case ".3mf":
		buf, err := io.ReadAll(fileHandle)
		if err != nil {
			return FileMetadata{}, err
		}

		reader := bytes.NewReader(buf)
		zipReader, err := zip.NewReader(reader, int64(len(buf)))
		if err != nil {
			return FileMetadata{}, err
		}

		var thumbnailBase64 string
		for _, file := range zipReader.File {
			if strings.HasSuffix(file.Name, "plate_1.png") {
				imgFile, err := file.Open()
				if err != nil {
					return FileMetadata{}, err
				}
				defer imgFile.Close()

				data, err := io.ReadAll(imgFile)
				if err != nil {
					return FileMetadata{}, err
				}

				thumbnailBase64 = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(data))

			}
		}

		return FileMetadata{Thumbnail: thumbnailBase64}, nil
	}

	return FileMetadata{}, nil
}
