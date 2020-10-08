package files

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalFileManager struct {
	base string
}

func (m *LocalFileManager) FetchFile(fileName string) ([]byte, error) {
	filePath, err := GetFileAbsPath(m.base, fileName)
	if err != nil {
		return nil, err
	}
	if exist, err := m.IsFileExists(filePath); err != nil {
		return nil, err
	} else {
		if exist {
			dat, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			return dat, nil
		} else {
			return nil, errors.New("file or directory does not exist")
		}
	}
}

func (m *LocalFileManager) GetBase() string {
	return m.base
}

func (m *LocalFileManager) SetBase(base string) {
	m.base = base
}

func GetFileAbsPath(base, fileName string) (fileAbsPath string, err error) {
	fileAbsPath = ""

	spacePath := path.Join(base)
	spaceAbsPath, err := filepath.Abs(spacePath)
	if err != nil {
		return
	}

	filePath := path.Join(base, fileName)
	fileAbsPath, err = filepath.Abs(filePath)
	if err != nil {
		return
	}

	if !strings.HasPrefix(fileAbsPath, spaceAbsPath) {
		return "", errors.New("escape from base path")
	}
	return
}

func (m *LocalFileManager) CreateFile(fileName string, bytes []byte) (err error) {
	filePath, err := GetFileAbsPath(m.base, fileName)
	if err != nil {
		return
	}
	if exist, err := m.IsFileExists(filePath); err != nil {
		return err
	} else {
		if exist {
			return errors.New("file or directory exists")
		} else {
			err = ioutil.WriteFile(filePath, bytes, os.FileMode(0644))
		}
	}
	return
}

func (m *LocalFileManager) CreateDirectory(fileName string) (err error) {
	filePath, err := GetFileAbsPath(m.base, fileName)
	if err != nil {
		return
	}
	if exist, err := m.IsFileExists(filePath); err != nil {
		return err
	} else {
		if exist {
			return errors.New("file or directory exists")
		} else {
			err = os.MkdirAll(filePath, os.FileMode(0644))
		}
	}
	return
}

func (m *LocalFileManager) IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (m *LocalFileManager) IsDirectoryExists(fileName string) (bool, error) {
	panic("implement me")
}

func (m *LocalFileManager) GetFilesAndDirs(dsirname string) (files []string, dirs []string, err error) {
	tmpDirs, err := ioutil.ReadDir(dsirname)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range tmpDirs {
		if file.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, file.Name())
		} else {
			files = append(files, file.Name())
		}
	}

	return files, dirs, nil
}
