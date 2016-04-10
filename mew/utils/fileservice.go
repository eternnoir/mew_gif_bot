package utils

import ()

type FileService interface {
	PutFile(buffer []byte, filepath string) error
	GetFilesUrl(path string) ([]string, error)
}
