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
	ModelBase64 string `json:"model_base64,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	FileType    string `json:"file_type,omitempty"` // "stl", "3mf", "gcode.3mf"
}

func GetFileMetadata(fileName string, fileHandle multipart.File) (FileMetadata, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".stl":
		content, err := io.ReadAll(fileHandle)
		if err != nil {
			return FileMetadata{}, err
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		return FileMetadata{ModelBase64: encoded, FileType: "stl"}, nil

	case ".3mf", ".gcode.3mf":
		buf, err := io.ReadAll(fileHandle)
		if err != nil {
			return FileMetadata{}, err
		}
		reader := bytes.NewReader(buf)
		zipReader, err := zip.NewReader(reader, int64(len(buf)))
		if err != nil {
			return FileMetadata{}, err
		}

		var modelBase64, thumbnailBase64 string
		var isSliced bool

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
				isSliced = true
			}
			if strings.HasSuffix(file.Name, ".stl") && !isSliced {
				stlFile, err := file.Open()
				if err != nil {
					return FileMetadata{}, err
				}
				defer stlFile.Close()
				data, err := io.ReadAll(stlFile)
				if err != nil {
					return FileMetadata{}, err
				}
				modelBase64 = base64.StdEncoding.EncodeToString(data)
			}
			if strings.HasSuffix(file.Name, "metadata.json") {
				isSliced = true
			}
		}

		if isSliced && thumbnailBase64 != "" {
			return FileMetadata{Thumbnail: thumbnailBase64, FileType: "gcode.3mf"}, nil
		}
		if modelBase64 != "" {
			return FileMetadata{ModelBase64: modelBase64, FileType: "3mf"}, nil
		}
		return FileMetadata{}, nil
	}
	return FileMetadata{}, nil
}
