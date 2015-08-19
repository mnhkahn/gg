package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

func Backup() {
	deleteFile(NewGGConfig().RunDirectory)
	// println("backup ")
}

func Deploy() {
	if err := CopyFile(NewGGConfig().AppPath, NewGGConfig().RunDirectory+NewGGConfig().AppName+".tar.gz"); err != nil {
		log.Println("Copy file error: ", err)
		return
	}

	if err := unPackFile(NewGGConfig().RunDirectory + NewGGConfig().AppName + ".tar.gz"); err != nil {
		log.Println("Tar package error: ", err)
	} else {
		pss, _ := ps.Processes()
		for _, p := range pss {
			if strings.Index(p.Executable(), NewGGConfig().AppName) != -1 {
				log.Printf("[pgrep %s] got pid: %d.\n", NewGGConfig().AppName, p.Pid())
				killProcess(p.Pid())
			}
		}
	}
}

func CopyFile(src, dst string) (err error) {
	log.Println("Copy file", src, "to", dst)

	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return err
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func killProcess(pid int) {
	log.Println("kill process", pid)
	cmd := fmt.Sprintf("kill %v", pid)
	runCommand(cmd)
}

func runCommand(cmd string) (string, error) {
	res, err := exec.Command("/bin/sh", "-c", cmd).Output()
	return string(res), err
}

func deleteFile(walkDir string) error {
	fileNames := make([]string, 0)
	dirNames := make([]string, 0)
	//遍历文件夹并把文件或文件夹名称加入相应的slice
	err := filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirNames = append(dirNames, path)
		} else {
			fileNames = append(fileNames, path)
		}
		return err
	})
	if err != nil {
		return err
	}
	//把所有文件名称连接成一个字符串
	fileNamesAll := strings.Join(fileNames, "")
	for i := len(dirNames) - 1; i >= 0; i-- {
		//文件夹名称不存在文件名称字符串内说明是个空文件夹
		if !strings.Contains(fileNamesAll, dirNames[i]) {
			log.Printf("%s is empty\n", dirNames[i])
			err := os.Remove(dirNames[i])
			if err != nil {
				return err
			} else {
				log.Println("Delete file", dirNames[i])
			}
		}
	}
	return nil
}
