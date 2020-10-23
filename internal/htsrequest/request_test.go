// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module request_test tests module request
package htsrequest

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

// requestIDTC test cases for Set/Get ID
var requestIDTC = []struct {
	id string
}{
	{"object1"},
	{"tabulamuris.00001"},
	{"1000genomes.99999"},
}

// requestFormatTC test cases for Set/Get Format
var requestFormatTC = []struct {
	format string
}{
	{"BAM"},
	{"CRAM"},
	{"VCF"},
	{"BCF"},
}

// requestClassTC test cases for Set/Get Class
var requestClassTC = []struct {
	class string
}{
	{"header"},
}

// requestReferenceNameTC test cases for Set/Get ReferenceName
var requestReferenceNameTC = []struct {
	referenceName string
}{
	{"chr1"},
	{"chr22"},
	{"chrMT"},
}

// requestStartTC test cases for Set/Get Start
var requestStartTC = []struct {
	start int
}{
	{10000},
	{20000000},
	{55000000},
}

// requestEndTC test cases for Set/Get End
var requestEndTC = []struct {
	end int
}{
	{999999},
	{90000000},
	{79000000},
}

// requestHtsgetBlockClassTC test cases for Set/Get HtsgetBlockClass
var requestHtsgetBlockClassTC = []struct {
	blockClass string
}{
	{"header"},
	{"body"},
}

// requestHtsgetCurrentBlockTC test cases for Set/Get HtsgetCurrentBlock
var requestHtsgetCurrentBlockTC = []struct {
	currentBlock string
}{
	{"0"},
	{"1"},
	{"100"},
}

// requestHtsgetTotalBlocksTC test cases for Set/Get HtsgetTotalBlocks
var requestHtsgetTotalBlocksTC = []struct {
	totalBlocks string
}{
	{"1"},
	{"10"},
	{"1000"},
}

// requestHtsgetFilePathTC test cases for Set/Get HtsgetFilePath
var requestHtsgetFilePathTC = []struct {
	filePath string
}{
	{"/path/to/the/file.bam"},
	{"./object1.vcf"},
	{"https://example.com/files/file99.cram"},
}

// requestRangeTC test cases for Set/Get Range
var requestRangeTC = []struct {
	Range string
}{
	{"bytes=10-20"},
	{"bytes=9999-9999999"},
	{"bytes=600-900"},
}

// requestFieldsTC test cases for Set/Get Fields
var requestFieldsTC = []struct {
	fields []string
}{
	{[]string{"SEQ", "QUAL"}},
	{[]string{"TLEN"}},
	{[]string{"QNAME", "FLAG", "TLEN", "SEQ", "QUAL"}},
}

// requestTagsTC test cases for Set/Get Tags
var requestTagsTC = []struct {
	tags []string
}{
	{[]string{"NM", "MD"}},
	{[]string{"NZ"}},
	{[]string{"NM", "NZ", "MD", "QL"}},
}

// requestNoTagsTC test cases for Set/Get NoTags
var requestNoTagsTC = []struct {
	notags []string
}{
	{[]string{"NM", "MD"}},
	{[]string{"NZ"}},
	{[]string{"NM", "NZ", "MD", "QL"}},
}

// requestHeaderOnlyRequestedTC test cases for HeaderOnlyRequested
var requestHeaderOnlyRequestedTC = []struct {
	class string
	exp   bool
}{
	{"", false},
	{"header", true},
	{"body", false},
}

// requestUnplacedUnmappedReadsRequestedTC test cases for UnplacedUnmappedReadsRequested
var requestUnplacedUnmappedReadsRequestedTC = []struct {
	referenceName string
	exp           bool
}{
	{"", false},
	{"chr1", false},
	{"chr22", false},
	{"*", true},
}

// requestReferenceNameRequestedTC test cases for ReferenceNameRequested
var requestReferenceNameRequestedTC = []struct {
	referenceName string
	exp           bool
}{
	{"", false},
	{"chr1", true},
	{"chr22", true},
}

// requestStartRequestedTC test cases for StartRequested
var requestStartRequestedTC = []struct {
	start int
	exp   bool
}{
	{-1, false},
	{100, true},
	{20000000, true},
}

// requestEndRequestedTC test cases for EndRequested
var requestEndRequestedTC = []struct {
	end int
	exp bool
}{
	{-1, false},
	{100, true},
	{20000000, true},
}

// requestAllRegionsRequestedTC test cases for AllRegionsRequested
var requestAllRegionsRequestedTC = []struct {
	regions []*Region
	exp     bool
}{
	{[]*Region{}, true},
	{[]*Region{&Region{ReferenceName: "chr1"}}, false},
	{[]*Region{&Region{ReferenceName: "chr22"}}, false},
}

// requestAllFieldsRequestedTC test cases for AllFieldsRequested
var requestAllFieldsRequestedTC = []struct {
	fields []string
	exp    bool
}{
	{[]string{"ALL"}, true},
	{[]string{"FLAG"}, false},
	{[]string{"QNAME", "SEQ", "QUAL"}, false},
}

// requestAllTagsRequestedTC test cases for AllTagsRequested
var requestAllTagsRequestedTC = []struct {
	tags   []string
	notags []string
	exp    bool
}{
	{[]string{"ALL"}, []string{"NONE"}, true},
	{[]string{"ALL"}, []string{"NM", "MD"}, false},
	{[]string{"NM", "MD"}, []string{"NONE"}, false},
}

// requestIsHeaderBlockTC test cases for IsHeaderBlock
var requestIsHeaderBlockTC = []struct {
	currentBlock string
	exp          bool
}{
	{"0", true},
	{"10", false},
	{"100", false},
}

// requestIsFinalBlockTC test cases for IsFinalBlock
var requestIsFinalBlockTC = []struct {
	currentBlock string
	totalBlocks  string
	exp          bool
}{
	{"0", "10", false},
	{"10", "100", false},
	{"9", "10", true},
}

// requestConstructDataEndpointURLTC test cases for ConstructDataEndpointURL
var requestConstructDataEndpointURLTC = []struct {
	endpoint                 htsconstants.APIEndpoint
	id, class, referenceName string
	start, end               int
	fields, tags, notags     []string
	exp                      string
	useRegion                bool
	regionI                  int
	expErr                   bool
	useBadConfig             bool
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"object0052",
		"",
		"chr1",
		65000,
		420000,
		defaultFields,
		defaultTags,
		defaultNoTags,
		"http://localhost:3000/reads/data/object0052?end=420000&referenceName=chr1&start=65000",
		true,
		0,
		false,
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		"",
		"chr22",
		11000000,
		45000000,
		[]string{"SEQ", "QUAL"},
		[]string{"NM", "HI"},
		defaultNoTags,
		"http://localhost:3000/reads/data/tabulamuris.00001?end=45000000&fields=SEQ%2CQUAL&referenceName=chr22&start=11000000&tags=NM%2CHI",
		true,
		0,
		false,
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		"header",
		defaultReferenceName,
		defaultStart,
		defaultEnd,
		defaultFields,
		defaultTags,
		defaultNoTags,
		"http://localhost:3000/reads/data/tabulamuris.00001?class=header",
		false,
		0,
		false,
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		defaultClass,
		defaultReferenceName,
		defaultStart,
		defaultEnd,
		defaultFields,
		defaultTags,
		[]string{"NM", "HI"},
		"http://localhost:3000/reads/data/tabulamuris.00001?notags=NM%2CHI",
		false,
		0,
		false,
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		defaultClass,
		defaultReferenceName,
		defaultStart,
		defaultEnd,
		defaultFields,
		defaultTags,
		[]string{"NM", "HI"},
		"http://localhost:3000/reads/data/tabulamuris.00001?notags=NM%2CHI",
		false,
		0,
		true,
		true,
	},
}

// requestDataSourceRegistryTC test cases for DataSourceRegistry
var requestDataSourceRegistryTC = []struct {
	endpoint          htsconstants.APIEndpoint
	expSource0Pattern string
	expSource0Path    string
}{
	{
		htsconstants.APIEndpointReadsTicket,
		htsconstants.DfltReadsDataSourceTabulaMuris10XPattern,
		htsconstants.DfltReadsDataSourceTabulaMuris10XPath,
	},
	{
		htsconstants.APIEndpointVariantsTicket,
		htsconstants.DfltVariantsDataSource1000GPattern,
		htsconstants.DfltVariantsDataSource1000GPath,
	},
}

var requestServiceInfoTC = []struct {
	endpoint    htsconstants.APIEndpoint
	expDatatype string
}{
	{htsconstants.APIEndpointReadsTicket, htsconstants.HtsgetExtensionDatatypeReads},
	{htsconstants.APIEndpointVariantsTicket, htsconstants.HtsgetExtensionDatatypeVariants},
}

// TestRequestID tests Set/Get ID functions
func TestRequestID(t *testing.T) {
	for _, tc := range requestIDTC {
		r := NewHtsgetRequest()
		r.SetID(tc.id)
		assert.Equal(t, tc.id, r.GetID())
	}
}

// TestRequestFormat tests Set/Get Format functions
func TestRequestFormat(t *testing.T) {
	for _, tc := range requestFormatTC {
		r := NewHtsgetRequest()
		r.SetFormat(tc.format)
		assert.Equal(t, tc.format, r.GetFormat())
	}
}

// TestRequestClass tests Set/Get Class functions
func TestRequestClass(t *testing.T) {
	for _, tc := range requestClassTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		assert.Equal(t, tc.class, r.GetClass())
	}
}

// TestRequestReferenceName tests Set/Get ReferenceName functions
func TestRequestReferenceName(t *testing.T) {
	for _, tc := range requestReferenceNameTC {
		r := NewHtsgetRequest()
		r.SetReferenceName(tc.referenceName)
		assert.Equal(t, tc.referenceName, r.GetReferenceName())
	}
}

// TestRequestStart tests Set/Get Start functions
func TestRequestStart(t *testing.T) {
	for _, tc := range requestStartTC {
		r := NewHtsgetRequest()
		r.SetStart(tc.start)
		assert.Equal(t, tc.start, r.GetStart())
	}
}

// TestRequestEnd tests Set/Get End functions
func TestRequestEnd(t *testing.T) {
	for _, tc := range requestEndTC {
		r := NewHtsgetRequest()
		r.SetEnd(tc.end)
		assert.Equal(t, tc.end, r.GetEnd())
	}
}

// TestRequestHtsgetBlockClass tests Set/Get HtsgetBlockClass functions
func TestRequestHtsgetBlockClass(t *testing.T) {
	for _, tc := range requestHtsgetBlockClassTC {
		r := NewHtsgetRequest()
		r.SetHtsgetBlockClass(tc.blockClass)
		assert.Equal(t, tc.blockClass, r.GetHtsgetBlockClass())
	}
}

// TestRequestHtsgetCurrentBlock tests Set/Get HtsgetCurrentBlock functions
func TestRequestHtsgetCurrentBlock(t *testing.T) {
	for _, tc := range requestHtsgetCurrentBlockTC {
		r := NewHtsgetRequest()
		r.SetHtsgetCurrentBlock(tc.currentBlock)
		assert.Equal(t, tc.currentBlock, r.GetHtsgetCurrentBlock())
	}
}

// TestRequestHtsgetTotalBlocks tests Set/Get HtsgetTotalBlocks functions
func TestRequestHtsgetTotalBlocks(t *testing.T) {
	for _, tc := range requestHtsgetTotalBlocksTC {
		r := NewHtsgetRequest()
		r.SetHtsgetTotalBlocks(tc.totalBlocks)
		assert.Equal(t, tc.totalBlocks, r.GetHtsgetTotalBlocks())
	}
}

// TestRequestFilePath tests Set/Get HtsgetFilePath functions
func TestRequestFilePath(t *testing.T) {
	for _, tc := range requestHtsgetFilePathTC {
		r := NewHtsgetRequest()
		r.SetHtsgetFilePath(tc.filePath)
		assert.Equal(t, tc.filePath, r.GetHtsgetFilePath())
	}
}

// TestRequestRange tests Set/Get HtsgetRange functions
func TestRequestRange(t *testing.T) {
	for _, tc := range requestRangeTC {
		r := NewHtsgetRequest()
		r.SetHtsgetRange(tc.Range)
		assert.Equal(t, tc.Range, r.GetHtsgetRange())
	}
}

// TestRequestFields tests Set/Get Fields functions
func TestRequestFields(t *testing.T) {
	for _, tc := range requestFieldsTC {
		r := NewHtsgetRequest()
		r.SetFields(tc.fields)
		assert.Equal(t, tc.fields, r.GetFields())
	}
}

// TestRequestTags tests Set/Get Tags functions
func TestRequestTags(t *testing.T) {
	for _, tc := range requestTagsTC {
		r := NewHtsgetRequest()
		r.SetTags(tc.tags)
		assert.Equal(t, tc.tags, r.GetTags())
	}
}

// TestRequestNoTags tests Set/Get NoTags functions
func TestRequestNoTags(t *testing.T) {
	for _, tc := range requestNoTagsTC {
		r := NewHtsgetRequest()
		r.SetNoTags(tc.notags)
		assert.Equal(t, tc.notags, r.GetNoTags())
	}
}

// TestRequestHeaderOnlyRequested tests HeaderOnlyRequested function
func TestRequestHeaderOnlyRequested(t *testing.T) {
	for _, tc := range requestHeaderOnlyRequestedTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		assert.Equal(t, tc.exp, r.HeaderOnlyRequested())
	}
}

// TestRequestUnplacedUnmappedReadsRequested tests UnplacedUnmappedReadsRequested function
func TestRequestUnplacedUnmappedReadsRequested(t *testing.T) {
	for _, tc := range requestUnplacedUnmappedReadsRequestedTC {
		r := NewHtsgetRequest()
		r.SetReferenceName(tc.referenceName)
		assert.Equal(t, tc.exp, r.UnplacedUnmappedReadsRequested())
	}
}

// TestRequestReferenceNameRequested tests ReferenceNameRequested function
func TestRequestReferenceNameRequested(t *testing.T) {
	for _, tc := range requestReferenceNameRequestedTC {
		r := NewHtsgetRequest()
		r.SetReferenceName(tc.referenceName)
		assert.Equal(t, tc.exp, r.ReferenceNameRequested())
	}
}

// TestRequestStartRequested tests StartRequested function
func TestRequestStartRequested(t *testing.T) {
	for _, tc := range requestStartRequestedTC {
		r := NewHtsgetRequest()
		r.SetStart(tc.start)
		assert.Equal(t, tc.exp, r.StartRequested())
	}
}

// TestRequestEndRequested tests EndRequested function
func TestRequestEndRequested(t *testing.T) {
	for _, tc := range requestEndRequestedTC {
		r := NewHtsgetRequest()
		r.SetEnd(tc.end)
		assert.Equal(t, tc.exp, r.EndRequested())
	}
}

// TestRequestAllRegionsRequested tests AllRegionsRequested function
func TestRequestAllRegionsRequested(t *testing.T) {
	for _, tc := range requestAllRegionsRequestedTC {
		r := NewHtsgetRequest()
		r.SetRegions(tc.regions)
		assert.Equal(t, tc.exp, r.AllRegionsRequested())
	}
}

// TestRequestAllFieldsRequested tests AllFieldsRequested function
func TestRequestAllFieldsRequested(t *testing.T) {
	for _, tc := range requestAllFieldsRequestedTC {
		r := NewHtsgetRequest()
		r.SetFields(tc.fields)
		assert.Equal(t, tc.exp, r.AllFieldsRequested())
	}
}

// TestRequestAllTagsRequested tests AllTagsRequested function
func TestRequestAllTagsRequested(t *testing.T) {
	for _, tc := range requestAllTagsRequestedTC {
		r := NewHtsgetRequest()
		r.SetTags(tc.tags)
		r.SetNoTags(tc.notags)
		assert.Equal(t, tc.exp, r.AllTagsRequested())
	}
}

// TestRequestIsHeaderBlock tests IsHeaderBlock function
func TestRequestIsHeaderBlock(t *testing.T) {
	for _, tc := range requestIsHeaderBlockTC {
		r := NewHtsgetRequest()
		r.SetHtsgetCurrentBlock(tc.currentBlock)
		assert.Equal(t, tc.exp, r.IsHeaderBlock())
	}
}

// TestRequestIsFinalBlock tests IsFinalBlock function
func TestRequestIsFinalBlock(t *testing.T) {
	for _, tc := range requestIsFinalBlockTC {
		r := NewHtsgetRequest()
		r.SetHtsgetCurrentBlock(tc.currentBlock)
		r.SetHtsgetTotalBlocks(tc.totalBlocks)
		assert.Equal(t, tc.exp, r.IsFinalBlock())
	}
}

// TestRequestConstructDataEndpointURL tests ConstructDataEndpointURL function
func TestRequestConstructDataEndpointURL(t *testing.T) {

	for _, tc := range requestConstructDataEndpointURLTC {

		// assign request properties
		request := NewHtsgetRequest()
		request.SetEndpoint(tc.endpoint)
		request.SetID(tc.id)
		request.SetClass(tc.class)
		request.SetReferenceName(tc.referenceName)
		request.SetStart(tc.start)
		request.SetEnd(tc.end)
		request.SetFields(tc.fields)
		request.SetTags(tc.tags)
		request.SetNoTags(tc.notags)

		region := &Region{ReferenceName: tc.referenceName, Start: &tc.start, End: &tc.end}
		request.AddRegion(region)

		// set the host to a badhost to trigger a url parse error
		if tc.useBadConfig {
			htsconfig.SetHost(":badhost")
		}

		// execute function, if error expected assert that it is not nil,
		// otherwise assert it is nil and assert value of url
		url, err := request.ConstructDataEndpointURL(tc.useRegion, tc.regionI)
		if tc.expErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.exp, url)
		}

		// reset the host to default
		htsconfig.SetHost(htsconstants.DfltServerPropsHost)
	}
}

// TestRequestGetDataSourceRegistry tests GetDataSourceRegistry function
func TestRequestGetDataSourceRegistry(t *testing.T) {
	for _, tc := range requestDataSourceRegistryTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		registry := r.GetDataSourceRegistry()
		assert.Equal(t, tc.expSource0Pattern, registry.Sources[0].Pattern)
		assert.Equal(t, tc.expSource0Path, registry.Sources[0].Path)
	}
}

// TestRequestGetServiceInfo tests GetServiceInfo function
func TestRequestGetServiceInfo(t *testing.T) {
	for _, tc := range requestServiceInfoTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		si := r.GetServiceInfo()
		assert.Equal(t, tc.expDatatype, si.HtsgetExtension.Datatype)
	}
}
