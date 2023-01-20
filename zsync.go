// 单线程版本

package main

import (
	"time"
	"math"
	"fmt"
	"flag"
	"os"
	"io"
	"io/ioutil"
	"path/filepath"
	"crypto/md5"
	"encoding/hex"
)

const TEMP_DIR_NAME string = "zsync_temp"

type fileInfo struct {
	fileName string
	fullPath string
	fileSize string
	md5      string
}

// 字节的单位转换 保留一位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
	   return fmt.Sprintf("%.1fB", float64(fileSize)/float64(1))
	} else if fileSize < int64(math.Pow(1024, 2)) {
	   return fmt.Sprintf("%.1fKB", float64(fileSize)/float64(1024))
	} else if fileSize < int64(math.Pow(1024, 3)) {
	   return fmt.Sprintf("%.1fMB", float64(fileSize)/float64(math.Pow(1024, 2)))
	} else if fileSize < int64(math.Pow(1024, 4)) {
	   return fmt.Sprintf("%.1fGB", float64(fileSize)/float64(math.Pow(1024, 3)))
	} else if fileSize < int64(math.Pow(1024, 5)) {
	   return fmt.Sprintf("%.1fTB", float64(fileSize)/float64(math.Pow(1024, 4)))
	} else {
	   return fmt.Sprintf("%.1fPB", float64(fileSize)/float64(math.Pow(1024, 5)))
	}
 }

func CalcFileMd5(filePath string) string {
	pFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open file failed, file=%v, err=%v\n", filePath, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	return hex.EncodeToString(md5h.Sum(nil))
}

func WalkDir(dirPath string, fileMap map[string]fileInfo) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range files {
		fullPath := filepath.Join(dirPath, v.Name())
		if v.IsDir() {
			WalkDir(fullPath, fileMap)
		} else {
			f, _ := os.Stat(fullPath)
			_md5 := CalcFileMd5(fullPath)
			if _md5 == "" {
				continue
			}
			_fileInfo := fileInfo{
				fileName: v.Name(),
				fullPath: fullPath,
				fileSize: FormatFileSize(f.Size()),
				md5: _md5,
			}
			fileMap[_md5] = _fileInfo
		}
	}
}

func CalcDirDiff(srcDirMap map[string]fileInfo, dstDirMap map[string]fileInfo) []fileInfo {
	var diff []fileInfo
	for _md5, _fileInfo := range srcDirMap {
		if _, ok := dstDirMap[_md5]; !ok {
			diff = append(diff, _fileInfo)
		}
	}
	return diff
}

func PrintDirDiff(dirDiff []fileInfo) {
	for idx, _fileInfo := range dirDiff {
		fmt.Printf("%4d %9s %s\n", idx+1, _fileInfo.fileSize, _fileInfo.fullPath)
	}
}

func CopyFile(srcFile, destFile string)(int64, error){
	fmt.Printf("copy %s ===> %s\n", srcFile, destFile)
    file1,err := os.Open(srcFile)
    if err != nil {
        return 0, err
    }
    file2, err := os.OpenFile(destFile,os.O_WRONLY|os.O_CREATE,os.ModePerm)
    if err != nil {
        return 0, err
    }
    defer file1.Close()
    defer file2.Close()
    return io.Copy(file2,file1)
}

func CopyToTempDir(dirDiff []fileInfo, dstDir string) {
	tempDir := filepath.Join(dstDir, TEMP_DIR_NAME)
	_, err := os.Stat(tempDir)
	if err != nil || os.IsNotExist(err) {
		err := os.Mkdir(tempDir, os.ModePerm)
		if err != nil {
			fmt.Printf("create temppoary directory %s failed.", tempDir)
			return
		}
	}

	for _, v := range dirDiff {
		dstPath := filepath.Join(tempDir, v.fileName)
		_, err := CopyFile(v.fullPath, dstPath)
		if err != nil {
			fmt.Printf("copy file %s failed, err=%e", v.fullPath, err)
		}
	}
}


func main() {
	t0 := time.Now()

	var srcDir = flag.String("src", "", "absolute path of source directory")
	var dstDir = flag.String("dst", "", "absolute path of destination directory")
	var copy = flag.Bool("c", false, "copy to temp dir in destination directory")
	flag.Parse()

	if *srcDir == "" {
		fmt.Println("srcDir is empty.")
		return
	}
	_, err := os.Stat(*srcDir)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("srcDir not exist.")
		return
	}

	if *dstDir == "" {
		fmt.Println("dstDir is empty.")
		return
	}
	_, err = os.Stat(*dstDir)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("dstDir not exist.")
		return
	}

	// key = md5(file)
	srcDirMap := make(map[string]fileInfo)
	dstDirMap := make(map[string]fileInfo)
	WalkDir(*srcDir, srcDirMap)
	WalkDir(*dstDir, dstDirMap)
	fmt.Printf("srcDirMap = %d\n", len(srcDirMap))
	fmt.Printf("dstDirMap = %d\n", len(dstDirMap))
	dirDiff := CalcDirDiff(srcDirMap, dstDirMap)
	PrintDirDiff(dirDiff)

	if *copy == true {
		CopyToTempDir(dirDiff, *dstDir)
	}

	fmt.Printf("=========== run time: %v\n", time.Since(t0))
}


