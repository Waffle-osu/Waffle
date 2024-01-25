package zip_utils

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ZipDirectory(zipfilename string, path string) error {
	outFile, err := os.Create(zipfilename)
	if err != nil {
		return err
	}

	w := zip.NewWriter(outFile)

	if err := AddFilesToZip(w, path, ""); err != nil {
		_ = outFile.Close()
		return err
	}

	if err := w.Close(); err != nil {
		_ = outFile.Close()
		return errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	if err := outFile.Close(); err != nil {
		return errors.New("Warning: closing zipfile failed: " + err.Error())
	}

	return nil
}

func AddFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())
		if _, err := os.Stat(fullfilepath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Mode()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if err := AddFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else if file.Mode().IsRegular() {
			dat, err := ioutil.ReadFile(fullfilepath)
			if err != nil {
				return err
			}
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else {
			// we ignore non-regular files because they are scary
		}
	}
	return nil
}

func UnzipFile(filename, dest string, skipOsus bool) error {
	archive, archiveOpenErr := zip.OpenReader(filename)

	if archiveOpenErr != nil {
		return archiveOpenErr
	}

	_, err := os.Stat(dest)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dirCreateErr := os.Mkdir(dest, os.ModePerm)

			//We just wanna make sure the directory exists
			//Weird if it did occur but you never know
			if dirCreateErr != nil && !errors.Is(dirCreateErr, os.ErrExist) {
				archive.Close()
				return dirCreateErr
			}
		} else {
			archive.Close()
			return err
		}
	}

	for _, f := range archive.File {
		filePath := filepath.Join(dest, f.Name)

		//For BSS purposes, They all get reuploaded, we don't need them twice
		if skipOsus && strings.HasSuffix(filePath, ".osu") {
			continue
		}

		if strings.ContainsAny(filePath, "?:*/<>|") {
			archive.Close()
			return errors.ErrUnsupported
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			archive.Close()
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			dstFile.Close()
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			fileInArchive.Close()
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	archive.Close()

	return nil
}
