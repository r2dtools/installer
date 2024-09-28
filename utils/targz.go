package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/unknwon/com"
)

func ExtractTarGz(file, folder string) error {
	gzipStream, err := os.Open(file)

	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)

	if err != nil {
		return fmt.Errorf("could not create new targz reader: %v", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("could not read next entry: %v", err)
		}

		destPath := filepath.Join(folder, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if com.IsDir(destPath) {
				continue
			}

			if err := os.Mkdir(destPath, 0755); err != nil {
				return fmt.Errorf("could not create directory '%s': %v", destPath, err)
			}
		case tar.TypeReg:
			if com.IsFile(destPath) {
				continue
			}

			outFile, err := os.Create(destPath)

			if err != nil {
				return fmt.Errorf("could not create file '%s': %v", destPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("could not copy entry '%s': %v", destPath, err)
			}

			outFile.Close()
		default:
			return fmt.Errorf("uknown type: %s in %s", string(header.Typeflag), destPath)
		}
	}

	return nil
}
