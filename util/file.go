package util

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

const AesIV string = "1234567890123456"

// DecryptFile 使用AES算法解密文件
//
//	@param inputFile
//	@param outputFile
//	@param key
//	@return error
func AesDecryptFile(inputFile, outputFile string, key []byte, inputIv string) error {
	// 打开加密文件
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	// 读取加密文件内容
	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// 初始化解密模式
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// 为了解密，我们需要使用相同的初始化向量IV
	// 在这个例子中，我们假设IV是硬编码的或者以某种方式传递给了解密函数
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// 使用CBC模式进行解密
	mode := cipher.NewCBCDecrypter(block, []byte(inputIv))

	// 初始化解密后的内容切片
	plaintext := make([]byte, len(ciphertext))

	// 去除填充
	mode.CryptBlocks(plaintext, ciphertext)

	// AES加密是块加密，所以我们需要去除填充
	padding := int(plaintext[len(plaintext)-1])
	plaintext = plaintext[:len(plaintext)-padding]

	// 写入解密后的内容到新文件
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(plaintext)
	if err != nil {
		return err
	}

	return nil
}

// AesEncryptFile
//
//	@param inputFile
//	@param outputFile
//	@param key
//	@return error
func AesEncryptFile(inputFile, outputFile string, key []byte, inputIv string) error {
	// 打开源文件
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	// 读取源文件内容
	input, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// 填充数据，使其长度为16的倍数
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	blockSize := block.BlockSize()
	padding := blockSize - len(input)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	input = append(input, padtext...)

	// 初始化加密模式
	mode := cipher.NewCBCEncrypter(block, []byte(inputIv))

	// 加密
	ciphertext := make([]byte, len(input))
	mode.CryptBlocks(ciphertext, input)

	// 写入加密内容到新文件
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(ciphertext)
	if err != nil {
		return err
	}

	return nil
}
