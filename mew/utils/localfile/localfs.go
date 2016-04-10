package localfile

import (
	"fmt"
	"io/ioutil"
)

type LocalFS struct {
	Path      string
	ServerUrl string
}

func (lfs *LocalFS) PutFile(buffer []byte, filepath string) error {
	return nil
}
func (lfs *LocalFS) GetFilesUrl(path string) ([]string, error) {
	files, _ := ioutil.ReadDir(lfs.Path + "/" + path)
	ret := []string{}
	for _, f := range files {
		url := lfs.ServerUrl + "/" + path + "/" + f.Name()
		fmt.Println(url)
		ret = append(ret, url)
	}
	return ret, nil
}
