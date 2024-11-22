package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archiver/v4"
)

func GetNowTimeStamp() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

// MakeVarLogFile
//
//	@param tag
//	@return string
func MakeVarLogFile(tag string) (string, error) {
	logFile := filepath.Join("/var/log", tag, GetNowTimeStamp()+".log")
	dir, _ := filepath.Split(logFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return logFile, nil
}

// MakeVarLogDir MakeVarLog
//
//	@param tag
//	@return string
func MakeVarLogDir(tag string) string {
	workDir := filepath.Join("/var/log", tag, GetNowTimeStamp())
	os.MkdirAll(workDir, os.ModePerm)
	return workDir
}

// MakeTmpWorkDir
//
//	@param tag
//	@return string
func MakeTmpWorkDir(tag string) string {
	workDir := filepath.Join("/tmp", tag, GetNowTimeStamp())
	os.MkdirAll(workDir, os.ModePerm)
	return workDir
}

// GetFiles
//
//	@param folder
//	@return []string
func GetFiles(folder string) []string {
	retData := []string{}
	err := filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		retData = append(retData, filepath.Join(folder, path))

		return nil
	})
	if err != nil {
		return retData
	}
	return retData
}

// GetFilesMap
//
//	@param fileDir
//	@return map
func GetFilesMap(folder string) map[string]string {
	retData := map[string]string{}
	filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		retData[path] = filepath.Join(folder, path)

		return nil
	})
	return retData
}

// GetDeviceSN
//
//	@return string
func GetDeviceSN() string {
	if bytes, err := os.ReadFile("/realsee/etc/device_sn"); err == nil {
		return strings.TrimSpace(string(bytes))
	}
	return "UnknownDevice"
}

// ReadConfig
//
//	@param dest
//	@return error
func ReadConfig(dest interface{}) error {
	if len(os.Args) < 2 {
		return errors.New("no config file")
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

// PxeServiceVersion
type PxeServiceVersion struct {
	Current map[string]struct {
		Version string `json:"version"`
	} `json:"Current"`
}

// GetPxeServiceVersion
//
//	@return string
func GetPxeServiceVersion() (*PxeServiceVersion, error) {
	bytes, err := os.ReadFile("/realsee/var/lib/ota/info.json")
	if err != nil {
		return nil, err
	}
	retData := &PxeServiceVersion{}
	err = json.Unmarshal(bytes, retData)
	return retData, err
}

// Zip
//
//	@param srcDir
//	@param dest
//	@return error
func Zip(srcDir, dest string, print bool) error {
	file, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	fileMap := map[string]string{}
	if fileStat.IsDir() {
		fileMap = GetFilesMap(srcDir)
	} else {
		// srcDir是一个文件
		_, name := filepath.Split(srcDir)
		fileMap[name] = srcDir
	}

	if print {
		for key := range fileMap {
			fmt.Println("Adding::", key)
		}
	}

	// map files on disk to their paths in the archive
	files, err := archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		return err
	}

	// create the output file we'll write to
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)
	format := archiver.Archive{
		Compression: nil,
		Archival:    archiver.Zip{},
	}

	// create the archive
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}
	return nil
}

// Unzip ...
//
//	@param srcFile
//	@param destSrc
//	@return error
func Unzip(srcFile, destSrc string, print bool) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	format := &archiver.Zip{}
	return format.Extract(context.Background(), file, func(ctx context.Context, f archiver.FileInfo) error {
		if f.IsDir() {
			return os.MkdirAll(filepath.Join(destSrc, f.NameInArchive), os.ModePerm)
		}
		reader, err := f.Open()
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("read %s error:%s", f.NameInArchive, err.Error())
		}

		tmpDest := filepath.Join(destSrc, f.NameInArchive)
		tmpDir := filepath.Dir(tmpDest)
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
				return err
			}
		}
		if print {
			fmt.Println("Extracting:", filepath.Join(destSrc, f.NameInArchive))
		}
		return os.WriteFile(filepath.Join(destSrc, f.NameInArchive), bytes, os.ModePerm)
	})
}

func JsonDecodeX(inputJSON string) (interface{}, error) {
	if len(inputJSON) < 1 {
		return nil, errors.New("cant decode empty string")
	}
	var data interface{}
	err := json.Unmarshal([]byte(inputJSON), &data)
	return data, err
}

// ExtractURL
//
//	@param value
//	@param matchFunc
//	@return []string
func ExtractURL(value interface{}, matchFunc func(string) bool) []string {
	if strValue, ok := value.(string); ok {
		if data, err := JsonDecodeX(strValue); err == nil {
			return ExtractURL(data, matchFunc)
		}
		if matchFunc(strValue) {
			return []string{strValue}
		}
		return []string{}
	}

	var retData []string
	if mapValue, ok := value.(map[string]interface{}); ok {
		for _, value := range mapValue {
			retData = append(retData, ExtractURL(value, matchFunc)...)
		}
		return retData
	}

	if listValue, ok := value.([]interface{}); ok {
		for _, value := range listValue {
			retData = append(retData, ExtractURL(value, matchFunc)...)
		}
		return retData
	}
	return nil
}

// ReadFileAs
//
//	@param path
//	@param v
//	@return error
func ReadFileAs(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

// FileExists
//
//	@param filename
//	@return bool
func FileExists(filename string) bool {
	if value, err := os.Stat(filename); err == nil && !value.IsDir() {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}
