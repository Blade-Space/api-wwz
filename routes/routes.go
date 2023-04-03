package wwf

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ZipRequest struct {
	Dst     string   `json:"dst"`
	Sources []string `json:"sources"`
}

type UnzipRequest struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func RegisterRoutes(api *gin.RouterGroup) {
	api.POST("/zip", ZipFilesHendler)
	api.POST("/unzip", UnZipHendler)
}

func ZipFilesHendler(c *gin.Context) {
	var zipRequest ZipRequest
	err := c.BindJSON(&zipRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = ZipFiles(zipRequest.Dst, zipRequest.Sources...)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ZIP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ZIP created successfully"})
}

func UnZipHendler(c *gin.Context) {
	var unzipRequest UnzipRequest
	err := c.BindJSON(&unzipRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = Unzip(unzipRequest.Src, unzipRequest.Dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unzip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unzipped successfully"})
}

// ZipFiles создает архив с указанными файлами и папками
func ZipFiles(dst string, sources ...string) error {
	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, src := range sources {
		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(path, string(filepath.Separator))

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip извлекает архив в указанную папку
func Unzip(src, dst string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dst, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}
