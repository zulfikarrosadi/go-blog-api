package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zulfikarrosadi/go-blog-api/web"
)

type UploadedResponse struct {
	FileName string `json:"fileName"`
}

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

func FileUploadHandler(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
		return err
	}

	files := form.File["files"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer src.Close()

		// check the file size
		if file.Size > MAX_UPLOAD_SIZE {
			fmt.Println("file size too big")
			return errors.New("uploaded file size is too big, 1MB is max")
		}

		// check the file format
		buff := make([]byte, 512)
		_, err = src.Read(buff)
		if err != nil {
			return err
		}

		fileType := http.DetectContentType(buff)
		if fileType != "image/jpeg" && fileType != "image/png" {
			fmt.Println("file format is unsupported: ", fileType)
			return errors.New("uploaded file format is not supported, please use png or jpeg format only")
		}

		_, err = src.Seek(0, io.SeekStart)
		if err != nil {
			fmt.Println(err)
			return err
		}

		dst, err := os.Create(fmt.Sprintf("./%v%v", time.Now().UnixNano(), file.Filename))
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			fmt.Println(err)
			return err
		}
	}

	response := web.Response{
		Status: "success",
		Data: UploadedResponse{
			FileName: "someshit",
		},
	}
	return c.JSON(http.StatusOK, response)
}
