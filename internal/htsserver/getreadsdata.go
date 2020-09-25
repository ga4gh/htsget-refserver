package htsserver

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
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsformats"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getReadsData(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointReadsData,
		noAfterSetup,
		getReadsDataHandler,
	).handleRequest(writer, request)
}

// getReadsData serves the actual data from AWS back to client
func getReadsDataHandler(handler *requestHandler) {
	fileURL, err := htsconfig.GetObjectPath(handler.HtsReq.GetEndpoint(), handler.HtsReq.GetID())
	if err != nil {
		return
	}

	start := handler.HtsReq.GetStart()
	end := handler.HtsReq.GetEnd()
	region := &htsrequest.Region{
		ReferenceName: handler.HtsReq.GetReferenceName(),
		Start:         &start,
		End:           &end,
	}

	args := getSamtoolsCmdArgs(region, handler.HtsReq, fileURL)
	cmd := exec.Command("samtools", args...)
	pipe, err := cmd.StdoutPipe()

	if err != nil {
		msg := err.Error()
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	err = cmd.Start()
	if err != nil {
		msg := err.Error()
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	reader := bufio.NewReader(pipe)

	var eofLen int
	if handler.HtsReq.GetHtsgetBlockClass() == "header" {
		eofLen = htsconstants.BamHeaderEOFLen
	} else {
		eofLen = htsconstants.BamEOFLen
	}

	if (handler.HtsReq.AllFieldsRequested() && handler.HtsReq.AllTagsRequested()) || handler.HtsReq.GetHtsgetBlockClass() == "header" {
		if handler.HtsReq.GetHtsgetBlockClass() != "header" { // remove header
			headerLen, err := headerLen(handler.HtsReq.GetID(), fileURL)
			handler.Writer.Header().Set("header-len", strconv.FormatInt(headerLen, 10))
			if err != nil {
				msg := err.Error()
				htserror.InternalServerError(handler.Writer, &msg)
				return
			}
			headerBuf := make([]byte, headerLen)
			io.ReadFull(reader, headerBuf)
		}
		if !handler.HtsReq.IsFinalBlock() { // remove EOF if current block is not the last block
			bufSize := 65536
			buf := make([]byte, bufSize)
			n, err := io.ReadFull(reader, buf)
			if err != nil && err.Error() != "unexpected EOF" {
				msg := err.Error()
				htserror.InternalServerError(handler.Writer, &msg)
				return
			}

			eofBuf := make([]byte, eofLen)
			for n == bufSize {
				copy(eofBuf, buf[n-eofLen:])
				handler.Writer.Write(buf[:n-eofLen])
				n, err = io.ReadFull(reader, buf)
				if err != nil && err.Error() != "unexpected EOF" {
					msg := err.Error()
					htserror.InternalServerError(handler.Writer, &msg)
					return
				}
				if n == bufSize {
					handler.Writer.Write(eofBuf)
				}
			}

			if n >= eofLen {
				handler.Writer.Write(buf[:n-eofLen])
			} else {
				handler.Writer.Write(eofBuf[:eofLen-n])
			}
		} else {
			io.Copy(handler.Writer, reader)
		}
	} else {
		columns := make([]bool, 11)
		for _, field := range handler.HtsReq.GetFields() {
			columns[htsconstants.BamFields[field]] = true
		}

		tmpPath := htsconfig.GetTempfilePath(handler.HtsReq.GetID())
		tmp, err := htsconfig.CreateTempfile(handler.HtsReq.GetID())
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
			return
		}

		/* Write the BAM Header to the temporary SAM file */
		tmpHeaderPath := htsconfig.GetTempfilePath(handler.HtsReq.GetID() + ".header.bam")
		headerCmd := exec.Command("samtools", "view", "-H", "-O", "SAM", "-o", tmpHeaderPath, fileURL)
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
		}
		err = headerCmd.Start()
		headerCmd.Wait()
		f, err := os.Open(tmpHeaderPath)
		headerReader := bufio.NewReader(f)
		hl, _, eof := headerReader.ReadLine()
		for ; eof == nil; hl, _, eof = headerReader.ReadLine() {
			_, err = tmp.Write([]byte(string(hl) + "\n"))
			if err != nil {
				msg := err.Error()
				htserror.InternalServerError(handler.Writer, &msg)
				return
			}
		}

		/* Write the custom SAM Records to the temporary SAM file */
		l, _, eof := reader.ReadLine()
		for ; eof == nil; l, _, eof = reader.ReadLine() {
			if l[0] != 64 {
				samRecord := htsformats.NewSAMRecord(string(l))
				newSamRecord := samRecord.CustomEmit(handler.HtsReq)
				l = []byte(newSamRecord + "\n")
			} else {
				l = append(l, "\n"...)
			}
			_, err = tmp.Write(l)
			if err != nil {
				msg := err.Error()
				htserror.InternalServerError(handler.Writer, &msg)
				return
			}
		}

		tmp.Close()
		bamCmd := exec.Command("samtools", "view", "-b", tmpPath)
		bamPipe, err := bamCmd.StdoutPipe()
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
			return
		}

		err = bamCmd.Start()
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
			return
		}

		// remove header bytes from 'body' class data streams
		headerByteCount, _ := headerLen(handler.HtsReq.GetID(), fileURL)
		bamReader := bufio.NewReader(bamPipe)
		headerBuf := make([]byte, headerByteCount)
		io.ReadFull(bamReader, headerBuf)
		io.Copy(handler.Writer, bamReader)

		err = bamCmd.Wait()
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
			return
		}

		err = htsconfig.RemoveTempfile(tmp)
		if err != nil {
			msg := err.Error()
			htserror.InternalServerError(handler.Writer, &msg)
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

func getSamtoolsCmdArgs(region *htsrequest.Region, htsgetReq *htsrequest.HtsgetRequest, fileURL string) []string {
	args := []string{"view", fileURL}
	if htsgetReq.GetHtsgetBlockClass() == "header" {
		args = append(args, "-H")
		args = append(args, "-b")
	} else {
		if htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested() {
			args = append(args, "-b")
		}
		if region.ExportSamtools() != "" {
			args = append(args, region.ExportSamtools())
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

func headerLen(id string, fileURL string) (int64, error) {
	cmd := exec.Command("samtools", "view", "-H", "-b", fileURL)
	tmpHeader, err := htsconfig.CreateTempfile(id + "_header")
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
	htsconfig.RemoveTempfile(tmpHeader)
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
