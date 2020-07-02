package htsgetserver

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/biogo/hts/bam"
	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/genomics"
	"github.com/ga4gh/htsget-refserver/internal/htsgeterror"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgetformats"
)

// getReadsData serves the actual data from AWS back to client
func getReadsData(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	htsgetReq, err := htsgetrequest.ReadsDataEndpointSetAllParameters(request, writer, params)

	if err != nil {
		return
	}

	region := &genomics.Region{
		Name:  htsgetReq.ReferenceName(),
		Start: htsgetReq.Start(),
		End:   htsgetReq.End(),
	}

	args := getSamtoolsCmdArgs(region, htsgetReq)
	cmd := exec.Command("samtools", args...)
	pipe, err := cmd.StdoutPipe()

	if err != nil {
		msg := err.Error()
		htsgeterror.InternalServerError(writer, &msg)
		return
	}

	err = cmd.Start()
	if err != nil {
		msg := err.Error()
		htsgeterror.InternalServerError(writer, &msg)
		return
	}

	reader := bufio.NewReader(pipe)

	var eofLen int
	if htsgetReq.HtsgetBlockClass() == "header" {
		eofLen = config.BamHeaderEOFLen
	} else {
		eofLen = config.BamEOFLen
	}

	if (htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested()) || htsgetReq.HtsgetBlockClass() == "header" {
		if htsgetReq.HtsgetBlockClass() != "header" { // remove header
			headerLen, err := headerLen(htsgetReq.ID())
			writer.Header().Set("header-len", strconv.FormatInt(headerLen, 10))
			if err != nil {
				msg := err.Error()
				htsgeterror.InternalServerError(writer, &msg)
				return
			}
			headerBuf := make([]byte, headerLen)
			io.ReadFull(reader, headerBuf)
		}
		if htsgetReq.HtsgetBlockID() != htsgetReq.HtsgetNumBlocks() { // remove EOF if current block is not the last block
			bufSize := 65536
			buf := make([]byte, bufSize)
			n, err := io.ReadFull(reader, buf)
			if err != nil && err.Error() != "unexpected EOF" {
				msg := err.Error()
				htsgeterror.InternalServerError(writer, &msg)
				return
			}

			eofBuf := make([]byte, eofLen)
			for n == bufSize {
				copy(eofBuf, buf[n-eofLen:])
				writer.Write(buf[:n-eofLen])
				n, err = io.ReadFull(reader, buf)
				if err != nil && err.Error() != "unexpected EOF" {
					msg := err.Error()
					htsgeterror.InternalServerError(writer, &msg)
					return
				}
				if n == bufSize {
					writer.Write(eofBuf)
				}
			}

			if n >= eofLen {
				writer.Write(buf[:n-eofLen])
			} else {
				writer.Write(eofBuf[:eofLen-n])
			}
		} else {
			io.Copy(writer, reader)
		}
	} else {
		columns := make([]bool, 11)
		for _, field := range htsgetReq.Fields() {
			columns[config.BamFields[field]] = true
		}

		tmpDirPath, err := tmpDirPath()
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}

		tmpPath := tmpDirPath + htsgetReq.ID()
		tmp, err := os.Create(tmpPath)
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}

		/* Write the BAM Header to the temporary SAM file */
		headerCmd := exec.Command("samtools", "view", "-H", config.DataSourceURL+htsgetutils.FilePath(htsgetReq.ID()))
		headerPipe, err := headerCmd.StdoutPipe()
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
		}
		err = headerCmd.Start()
		headerReader := bufio.NewReader(headerPipe)
		hl, _, eof := headerReader.ReadLine()
		for ; eof == nil; hl, _, eof = headerReader.ReadLine() {
			_, err = tmp.Write([]byte(string(hl) + "\n"))
			if err != nil {
				msg := err.Error()
				htsgeterror.InternalServerError(writer, &msg)
				return
			}
		}

		/* Write the custom SAM Records to the temporary SAM file */
		l, _, eof := reader.ReadLine()
		for ; eof == nil; l, _, eof = reader.ReadLine() {
			if l[0] != 64 {
				samRecord := htsgetformats.NewSAMRecord(string(l))
				newSamRecord := samRecord.CustomEmit(htsgetReq)
				l = []byte(newSamRecord + "\n")
			} else {
				l = append(l, "\n"...)
			}
			_, err = tmp.Write(l)
			if err != nil {
				msg := err.Error()
				htsgeterror.InternalServerError(writer, &msg)
				return
			}
		}

		tmp.Close()
		bamCmd := exec.Command("samtools", "view", "-b", tmpPath)
		bamPipe, err := bamCmd.StdoutPipe()
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}

		err = bamCmd.Start()
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}

		// remove header bytes from 'body' class data streams
		headerByteCount, _ := headerLen(htsgetReq.ID())
		bamReader := bufio.NewReader(bamPipe)
		headerBuf := make([]byte, headerByteCount)
		io.ReadFull(bamReader, headerBuf)
		io.Copy(writer, bamReader)

		err = bamCmd.Wait()
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}

		err = os.Remove(tmpPath)
		if err != nil {
			msg := err.Error()
			htsgeterror.InternalServerError(writer, &msg)
			return
		}
	}
	cmd.Wait()
}

func getTempPath(id string, blockID int) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(cwd)
	tempPath := parent + "/temp/" + id + "_" + strconv.Itoa(blockID)

	exists, err := fileExists(parent + "/temp")
	if err != nil {
		return "", err
	}
	if !exists {
		os.Mkdir(parent+"/temp/", 0755)
	} else {
		isDir, err := isDir(parent + "/temp")
		if err != nil {
			return "", err
		}
		if !isDir {
			os.Mkdir(parent+"/temp/", 0755)
		}
	}
	return tempPath, nil
}

func getSamtoolsCmdArgs(region *genomics.Region, htsgetReq *htsgetrequest.HtsgetRequest) []string {
	args := []string{"view", config.DataSourceURL + htsgetutils.FilePath(htsgetReq.ID())}
	if htsgetReq.HtsgetBlockClass() == "header" {
		args = append(args, "-H")
		args = append(args, "-b")
	} else {
		if htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested() {
			args = append(args, "-b")
		}
		if region.String() != "" {
			args = append(args, region.String())
		}
	}
	return args
}

func samToBam(tempPath string) string {
	bamPath := tempPath + "_bam"
	cmd := exec.Command("samtools", "view", "-h", "-b", tempPath, "-o", bamPath)
	cmd.Run()
	return bamPath
}

func headerLen(id string) (int64, error) {
	cmd := exec.Command("samtools", "view", "-H", "-b", config.DataSourceURL+htsgetutils.FilePath(id))
	tmpDirPath, err := tmpDirPath()
	if err != nil {
		return 0, err
	}
	path := tmpDirPath + id + "_header"
	tmpHeader, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	cmd.Stdout = tmpHeader
	cmd.Run()

	fi, err := tmpHeader.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size() - 12
	tmpHeader.Close()
	os.Remove(path)
	return size, nil
}

func removeHeader(path string) error {
	fin, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fin.Close()
	fin.Seek(0, 0)
	b, err := bam.NewReader(fin, 0)
	if err != nil {
		return err
	}

	b.Header()
	b.Read()
	lastChunk := b.LastChunk()
	hLen := lastChunk.Begin.File

	fDest, err := os.Create(path + "_copy")
	if err != nil {
		return err
	}

	fin.Seek(hLen, 0)
	io.Copy(fin, fDest)

	os.Remove(path)
	os.Rename(path+"_copy", path)
	return nil
}

func removeEOF(tempPath string) error {
	fi, err := os.Stat(tempPath)
	if err != nil {
		return err
	}
	return os.Truncate(tempPath, fi.Size()-28)
}

// exists returns whether the given file or directory exists.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}

// isDir returns whether the given path is directory
func isDir(path string) (bool, error) {
	src, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if src.Mode().IsRegular() {
		fmt.Println(path + " already exist as a file!")
		return false, nil
	}
	return true, nil
}

func tmpDirPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd + "/temp/", nil
}
