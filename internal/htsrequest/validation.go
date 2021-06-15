// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module validation defines functions for validating whether the value of
// HTTP request parameters are acceptable
package htsrequest

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
	"github.com/ga4gh/htsget-refserver/internal/awsutils"
)

// ParamValidator validates request parameters
type ParamValidator struct{}

// NewParamValidator instantiates a new ParamValidator object
func NewParamValidator() *ParamValidator {
	return new(ParamValidator)
}

// errorsByParam (map[string]func(http.ResponseWriter, *string)): the correct
// error to raise for each request parameter validation
var errorsByParam = map[string]func(http.ResponseWriter, *string){
	"id":               htserror.NotFound,
	"format":           htserror.UnsupportedFormat,
	"class":            htserror.InvalidInput,
	"referenceName":    htserror.InvalidRange,
	"start":            htserror.InvalidRange,
	"end":              htserror.InvalidRange,
	"fields":           htserror.InvalidInput,
	"tags":             htserror.InvalidInput,
	"notags":           htserror.InvalidInput,
	"regions":          htserror.InvalidRange,
	"HtsgetBlockClass": htserror.InvalidInput,
	"HtsgetBlockId":    htserror.InternalServerError,
	"HtsgetNumBlocks":  htserror.InternalServerError,
	"HtsgetFilePath":   htserror.InternalServerError,
	"Range":            htserror.InternalServerError,
}

// isGreaterThanEqualToZero determines if a string can be parsed as an integer,
// and is greater than or equal to zero
func isGreaterThanEqualToZero(num int) bool {
	if num < 0 {
		return false
	}
	return true
}

// NoValidation is an empty validation function for request parameters that do
// not need to be validated. always returns true
func (v *ParamValidator) NoValidation(htsgetReq *HtsgetRequest, value string) (bool, string) {
	return true, ""
}

// ValidateID validates the 'id' parameter. checks if an object matching
// the 'id' could be found from the data source
func (v *ParamValidator) ValidateID(htsgetReq *HtsgetRequest, id string) (bool, string) {
	objPath, err := htsconfig.GetObjectPath(htsgetReq.GetEndpoint(), id)
	if err != nil {
		return false, "The requested resource could not be associated with a registered data source"
	}

	// attempt to locate the object by http request (if url) or on local file
	// path
	if htsutils.IsValidURL(objPath) {
		if strings.HasPrefix(objPath, awsutils.S3Proto) {
			_, err := awsutils.HeadS3Object(awsutils.S3Dto{
				ObjPath: objPath,
			})
			if err != nil {
				return false, "Error accessing the requested S3 resource"
			}
		} else {
			res, err := http.Head(objPath)
			if err != nil {
				return false, "The requested resource was not found"
			}
			res.Body.Close()
			if res.Status == "404 Not Found" {
				return false, "The requested resource was not found"
			}
		}
	} else {
		_, err := os.Stat(objPath)
		if os.IsNotExist(err) {
			return false, "The requested resource was not found"
		}
	}
	return true, ""
}

// ValidateFormat validates the 'format' parameter. checks if the requested
// format is one of the allowed options based on endpoint
func (v *ParamValidator) ValidateFormat(htsgetReq *HtsgetRequest, format string) (bool, string) {
	allowedFormats := htsgetReq.GetEndpoint().AllowedFormats()
	if !htsutils.IsItemInArray(format, allowedFormats) {
		return false, "file format: '" + format + "' not supported"
	}
	return true, ""
}

// ValidateClass validates the 'class' parameter. checks if the requested class
// is one of the allowed options
func (v *ParamValidator) ValidateClass(htsgetReq *HtsgetRequest, class string) (bool, string) {
	switch strings.ToLower(class) {
	case htsconstants.ClassHeader:
		return true, ""
	case htsconstants.ClassBody:
		return false, "'body' only requests currently not supported" // currently not supported
	default:
		return false, "class: '" + class + "' not supported"
	}
}

func getReferenceNamesInReadsObject(htsgetReq *HtsgetRequest) ([]string, error) {

	var referenceNames []string
	fileURL, err := htsconfig.GetObjectPath(htsgetReq.GetEndpoint(), htsgetReq.GetID())
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("samtools", "view", "-H", fileURL)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Start()
	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()

	for ; err == nil; l, _, err = reader.ReadLine() {
		pattern := regexp.MustCompile("^@SQ\tSN:(.+?)\t.+?$")
		submatches := pattern.FindStringSubmatch(string(l))
		if len(submatches) > 1 {
			referenceNames = append(referenceNames, submatches[1])
		}
	}

	// wait till job completes, if the exit code wasn't 0, then raise an error
	// likely the file was not found
	err = cmd.Wait()
	if err != nil {
		return nil, errors.New("Could not get referenceNames from requested alignment file")
	}

	return referenceNames, nil
}

func getReferenceNamesInVariantsObject(htsgetReq *HtsgetRequest) ([]string, error) {
	var referenceNames []string
	fileURL, err := htsconfig.GetObjectPath(htsgetReq.GetEndpoint(), htsgetReq.GetID())
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("bcftools", "view", "-h", fileURL)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Start()
	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()

	for ; err == nil; l, _, err = reader.ReadLine() {
		pattern := regexp.MustCompile("^##contig=<.*?ID=(.+?)[,>]")
		submatches := pattern.FindStringSubmatch(string(l))
		if len(submatches) > 1 {
			referenceNames = append(referenceNames, submatches[1])
		}
	}

	// wait till job completes, if the exit code wasn't 0, then raise an error
	// likely the file was not found
	err = cmd.Wait()
	if err != nil {
		return nil, errors.New("Could not get referenceNames from requested variant file")
	}
	return referenceNames, nil
}

// getAllowedReferenceNames
// for a given endpoint (BAM request / VCF request), return the allowable values
// for the 'referenceName' parameter for the requested object
func getReferenceNames(htsgetReq *HtsgetRequest) ([]string, error) {
	functions := map[htsconstants.APIEndpoint]func(htsgetReq *HtsgetRequest) ([]string, error){
		htsconstants.APIEndpointReadsTicket:    getReferenceNamesInReadsObject,
		htsconstants.APIEndpointReadsData:      getReferenceNamesInReadsObject,
		htsconstants.APIEndpointVariantsTicket: getReferenceNamesInVariantsObject,
		htsconstants.APIEndpointVariantsData:   getReferenceNamesInVariantsObject,
	}
	return functions[htsgetReq.endpoint](htsgetReq)
}

// ValidateReferenceName validates the 'referenceName' query string
// parameter. checks if the requested reference contig/chromosome is in the
// BAM/CRAM header sequence dictionary. if unplaced unmapped reads are requested
// (*), do not perform validation
func (v *ParamValidator) ValidateReferenceName(htsgetReq *HtsgetRequest, referenceName string) (bool, string) {

	// incompatible with header only request
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'referenceName' incompatible with header-only request"
	}
	// no validation if '*' was requested
	if referenceName == "*" {
		return true, ""
	}
	// otherwise, check that referenceName is in the header
	referenceNames, err := getReferenceNames(htsgetReq)
	if err != nil {
		return false, err.Error()
	}

	if htsutils.IsItemInArray(referenceName, referenceNames) {
		return true, ""
	}

	return false, "invalid 'referenceName': " + referenceName
}

// ValidateStart validates the 'start' query string parameter. checks that it is
// a valid, non-zero integer, and that it is being used correctly in conjunction
// with 'referenceName'
func (v *ParamValidator) ValidateStart(htsgetReq *HtsgetRequest, start int) (bool, string) {

	// incompatible with header only request
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'start' incompatible with header-only request"
	}

	// start requires referenceName to specify a true chromosome
	if htsgetReq.UnplacedUnmappedReadsRequested() {
		return false, "'start' cannot be requested with unplaced, unmapped reads"
	}

	// start requires referenceName to be specified as well
	if !htsgetReq.ReferenceNameRequested() {
		return false, "'start' cannot be set without 'referenceName'"
	}

	// start must be >= 0
	if !isGreaterThanEqualToZero(start) {
		return false, "'start' must be greater than or equal to zero"
	}

	return true, ""
}

// ValidateEnd validates the 'end' query string parameter. checks that it is a
// valid, non-zero integer, that it's being used correctly in conjunction with
// 'referenceName', and that the end coordinate is greater than the start
// coordinate
func (v *ParamValidator) ValidateEnd(htsgetReq *HtsgetRequest, end int) (bool, string) {

	// incompatible with header only request
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'end' incompatible with header-only request"
	}

	start := htsgetReq.GetStart()

	// end requires referenceName to specify a true chromosome
	if htsgetReq.UnplacedUnmappedReadsRequested() {
		return false, "'end' cannot be requested with unplaced, unmapped reads"
	}

	// end requires referenceName to be specified as well
	if !htsgetReq.ReferenceNameRequested() {
		return false, "'end' cannot be set without 'referenceName'"
	}

	// end must be >= 0
	if !isGreaterThanEqualToZero(end) {
		return false, "'end' must be greater than or equal to zero"
	}

	// if start is specified, end must be greater than start
	if start != -1 {
		if start >= end {
			return false, "'end' MUST be higher than 'start'"
		}
	}
	return true, ""
}

// ValidateFields validates 'fields' parameter. every requested field must
// have an acceptable BAM/CRAM column name
func (v *ParamValidator) ValidateFields(htsgetReq *HtsgetRequest, fields []string) (bool, string) {

	// incompatible with header only request
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'fields' incompatible with header-only request"
	}

	for _, fieldItem := range fields {
		if _, ok := htsconstants.BamFields[fieldItem]; !ok {
			return false, "'" + fieldItem + "' not an acceptable field"
		}
	}
	return true, ""
}

// ValidateTags only validates that tags hasn't been requested alongside
// a header only request
func (v *ParamValidator) ValidateTags(htsgetReq *HtsgetRequest, tags []string) (bool, string) {
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'tags' incompatible with header-only request"
	}
	return true, ""
}

// ValidateNoTags validates the 'notags' query string parameter. checks that
// there is no overlap between tags included by 'tags' and tags excluded by
// 'notags'
func (v *ParamValidator) ValidateNoTags(htsgetReq *HtsgetRequest, notags []string) (bool, string) {

	// incompatible with header only request
	if htsgetReq.HeaderOnlyRequested() {
		return false, "'notags' incompatible with header-only request"
	}

	tags := htsgetReq.GetTags()
	for _, tagItem := range tags {
		for _, notagItem := range notags {
			if tagItem == notagItem {
				return false, "'" + tagItem + "' cannot be in both 'tags' and 'notags'"
			}
		}
	}
	return true, ""
}

// ValidateRegions validates whether every region within an array of regions is
// valid, that is, contains acceptable referenceName, start, and end values
func (v *ParamValidator) ValidateRegions(htsgetReq *HtsgetRequest, regions []*Region) (bool, string) {

	allowedReferenceNames, err := getReferenceNames(htsgetReq)
	if err != nil {
		return false, err.Error()
	}

	for _, region := range regions {

		if region.ReferenceNameRequested() {
			if !htsutils.IsItemInArray(region.GetReferenceName(), allowedReferenceNames) {
				return false, "Invalid referenceName in regions list: '" + region.GetReferenceName() + "'"
			}
		}

		if region.StartRequested() {
			if !region.ReferenceNameRequested() {
				return false, "Invalid region(s): 'start' cannot be set without 'referenceName'"
			}
			if !isGreaterThanEqualToZero(region.GetStart()) {
				return false, "Invalid region(s): 'start' MUST be greater than or equal to zero"
			}
		}

		if region.EndRequested() {
			if !region.ReferenceNameRequested() {
				return false, "Invalid region(s): 'end' cannot be set without 'referenceName'"
			}
			if !isGreaterThanEqualToZero(region.GetEnd()) {
				return false, "Invalid region(s): 'end' MUST be greater than or equal to zero"
			}
			if region.StartRequested() {
				if region.GetStart() >= region.GetEnd() {
					return false, "Invalid region(s): 'end' MUST be greater than 'start'"
				}
			}
		}
	}
	return true, ""
}
