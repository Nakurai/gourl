//go:build exclude

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func compressFile(src string, dest string) error {
	archive, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	fileToZip, err := os.Open(src)
	if err != nil {
		return err
	}

	fileName := filepath.Base(src)
	zipFileWriter, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFileWriter, fileToZip)
	if err != nil {
		return err
	}

	return nil

}

type CompileOption struct {
	Os          string
	Arch        string
	BinPath     string
	ArchivePath string
}

func main() {
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Error: You need to set the VERSION env variable before compiling")
		return
	}

	versionFlag := fmt.Sprintf("-X main.version=%s", version)

	compileOptions := []CompileOption{
		CompileOption{Os: "windows", Arch: "amd64", BinPath: fmt.Sprintf("./releases/%s/%s-gourl-windows.exe", version, version)},
		CompileOption{Os: "linux", Arch: "amd64", BinPath: fmt.Sprintf("./releases/%s/%s-gourl-linux-amd64", version, version)},
		CompileOption{Os: "linux", Arch: "arm64", BinPath: fmt.Sprintf("./releases/%s/%s-gourl-linux-arm64", version, version)},
		CompileOption{Os: "darwin", Arch: "amd64", BinPath: fmt.Sprintf("./releases/%s/%s-gourl-macos-amd64.app", version, version), ArchivePath: fmt.Sprintf("./releases/%s/%s-gourl-macos-amd64.zip", version, version)},
		CompileOption{Os: "darwin", Arch: "arm64", BinPath: fmt.Sprintf("./releases/%s/%s-gourl-macos-arm64.app", version, version), ArchivePath: fmt.Sprintf("./releases/%s/%s-gourl-macos-arm64.zip", version, version)},
	}

	for _, option := range compileOptions {
		os.Setenv("GOOS", option.Os)
		os.Setenv("GOARCH", option.Arch)
		cmd := exec.Command(
			"go", "build",
			"-ldflags", versionFlag,
			"-o", option.BinPath,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("ko %s/%s: %s\n", option.Os, option.Arch, output)
			break
		}
		if option.ArchivePath != "" {
			err := compressFile(option.BinPath, option.ArchivePath)
			if err != nil {
				fmt.Printf("ko %s/%s: %s\n", option.Os, option.Arch, err.Error())
				break
			}
		}
		fmt.Printf("ok %s/%s\n", option.Os, option.Arch)

	}

	fmt.Println("done")

}
