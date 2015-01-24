package fio // custom io/ioutil and os file mirroring and tweaking for convenience

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var (
	Truncate = os.Truncate
	Rename   = os.Rename
	R2T      = Rename2Time

	AF  = AppendFile
	DD  = DeleteDirectory
	CFC = CreateFileClose //Touch
	FE  = FileExists
	LDF = ListDirFiles
	MD  = Makedir
	MDA = MakedirAll
	RF  = ioutil.ReadFile
	RWF = ReadWriteFile
	WF  = WriteFile //WFT for Truncate?
	WFA = WriteFileAt

	DFP  = os.FileMode(0775)
	RWFP = os.FileMode(0666)
	WFP  = os.FileMode(0644)
)

/*
type File struct {
	os.File
}
*/

func ListDirFiles(path string) (a []string) {
	f,_ := ioutil.ReadDir(path)
	if f == nil { return }
	for _,v := range f {
		if !v.IsDir() {
			a = append(a,v.Name())
		}
	}
	return
}

func Rename2Time(oldpath string) error {
	return Rename(oldpath, oldpath+"_"+strconv.Itoa(int(time.Now().UnixNano())))
}

func WriteFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, WFP)
}

func WriteFileAt(filename string, data []byte, off int64) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_EXCL, WFP)
	if err != nil {
		return err
	}
	n, err := f.WriteAt(data, off)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func AppendFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, WFP)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func Mkdir(filename string, perm os.FileMode) error {
	return os.Mkdir(filename, perm)
}

func Makedir(filename string) error {
	return os.Mkdir(filename, DFP)
}

func ReadWriteFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, RWFP)
}

func ReadWriteFileAt(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR, RWFP)
}

func MkdirAll(filename string, perm os.FileMode) error {
	return os.MkdirAll(filename, perm)
}

func MakedirAll(filename string) error {
	return os.MkdirAll(filename, DFP)
}

func DeleteDirectory(filename string) error {
	return os.RemoveAll(filename)
}

func CreateFileClose(filename string) error {
	if f, err := os.OpenFile(filename, os.O_CREATE, RWFP); err != nil {
		return err
	} else {
		f.Close()
	}
	return nil
}

func FileExists(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_EXCL, WFP); err == nil {
		f.Close()
		return true
	}
	return false //More 'can use' than exists
}
