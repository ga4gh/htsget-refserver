package htsrequest

/*

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

var requestIDTC = []struct {
	id string
}{
	{"object1"},
	{"tabulamuris.00001"},
	{"1000genomes.99999"},
}

var requestFormatTC = []struct {
	format string
}{
	{"BAM"},
	{"CRAM"},
	{"VCF"},
	{"BCF"},
}

var requestClassTC = []struct {
	class string
}{
	{"header"},
}

var requestReferenceNameTC = []struct {
	referenceName string
}{
	{"chr1"},
	{"chr22"},
	{"chrMT"},
}

var requestStartTC = []struct {
	start string
}{
	{"10000"},
	{"20000000"},
	{"55000000"},
}

var requestEndTC = []struct {
	end string
}{
	{"999999"},
	{"90000000"},
	{"79000000"},
}

var requestHtsgetBlockClassTC = []struct {
	blockClass string
}{
	{"header"},
	{"body"},
}

var requestHtsgetBlockIDTC = []struct {
	blockID string
}{
	{"0"},
	{"1"},
	{"100"},
}

var requestHtsgetNumBlocksTC = []struct {
	numBlocks string
}{
	{"1"},
	{"10"},
	{"1000"},
}

var requestHtsgetFilePathTC = []struct {
	filePath string
}{
	{"/path/to/the/file.bam"},
	{"./object1.vcf"},
	{"https://example.com/files/file99.cram"},
}

var requestRangeTC = []struct {
	Range string
}{
	{"bytes=10-20"},
	{"bytes=9999-9999999"},
	{"bytes=600-900"},
}

var requestFieldsTC = []struct {
	fields []string
}{
	{[]string{"SEQ", "QUAL"}},
	{[]string{"TLEN"}},
	{[]string{"QNAME", "FLAG", "TLEN", "SEQ", "QUAL"}},
}

var requestTagsTC = []struct {
	tags []string
}{
	{[]string{"NM", "MD"}},
	{[]string{"NZ"}},
	{[]string{"NM", "NZ", "MD", "QL"}},
}

var requestNoTagsTC = []struct {
	notags []string
}{
	{[]string{"NM", "MD"}},
	{[]string{"NZ"}},
	{[]string{"NM", "NZ", "MD", "QL"}},
}

var requestHeaderOnlyRequestedTC = []struct {
	class string
	exp   bool
}{
	{"", false},
	{"header", true},
	{"body", false},
}

var requestUnplacedUnmappedReadsRequestedTC = []struct {
	referenceName string
	exp           bool
}{
	{"", false},
	{"chr1", false},
	{"chr22", false},
	{"*", true},
}

var requestReferenceNameRequestedTC = []struct {
	referenceName string
	exp           bool
}{
	{"", false},
	{"chr1", true},
	{"chr22", true},
}

var requestStartRequestedTC = []struct {
	start string
	exp   bool
}{
	{"-1", false},
	{"100", true},
	{"20000000", true},
}

var requestEndRequestedTC = []struct {
	end string
	exp bool
}{
	{"-1", false},
	{"100", true},
	{"20000000", true},
}

var requestAllRegionsRequestedTC = []struct {
	referenceName string
	exp           bool
}{
	{"", true},
	{"chr1", false},
	{"chr22", false},
}

var requestAllFieldsRequestedTC = []struct {
	fields []string
	exp    bool
}{
	{[]string{"ALL"}, true},
	{[]string{"FLAG"}, false},
	{[]string{"QNAME", "SEQ", "QUAL"}, false},
}

var requestAllTagsRequestedTC = []struct {
	tags   []string
	notags []string
	exp    bool
}{
	{[]string{"ALL"}, []string{"NONE"}, true},
	{[]string{"ALL"}, []string{"NM", "MD"}, false},
	{[]string{"NM", "MD"}, []string{"NONE"}, false},
}

var requestConstructDataEndpointURLTC = []struct {
	endpoint                             htsconstants.APIEndpoint
	id, class, referenceName, start, end string
	fields, tags, notags                 []string
	exp                                  string
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"object0052",
		"",
		"chr1",
		"65000",
		"420000",
		defaultListParameterValues["fields"],
		defaultListParameterValues["tags"],
		defaultListParameterValues["notags"],
		"http://localhost:3000/reads/data/object0052?end=420000&referenceName=chr1&start=65000",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		"",
		"chr22",
		"11000000",
		"45000000",
		[]string{"SEQ", "QUAL"},
		[]string{"NM", "HI"},
		defaultListParameterValues["notags"],
		"http://localhost:3000/reads/data/tabulamuris.00001?end=45000000&fields=SEQ%2CQUAL&referenceName=chr22&start=11000000&tags=NM%2CHI",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		"header",
		defaultScalarParameterValues["referenceName"],
		defaultScalarParameterValues["start"],
		defaultScalarParameterValues["end"],
		defaultListParameterValues["fields"],
		defaultListParameterValues["tags"],
		defaultListParameterValues["notags"],
		"http://localhost:3000/reads/data/tabulamuris.00001?class=header",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		defaultScalarParameterValues["class"],
		defaultScalarParameterValues["referenceName"],
		defaultScalarParameterValues["start"],
		defaultScalarParameterValues["end"],
		defaultListParameterValues["fields"],
		defaultListParameterValues["tags"],
		[]string{"NM", "HI"},
		"http://localhost:3000/reads/data/tabulamuris.00001?notags=NM%2CHI",
	},
}

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

func TestRequestID(t *testing.T) {
	for _, tc := range requestIDTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("id", tc.id)
		assert.Equal(t, tc.id, r.ID())
	}
}

func TestRequestFormat(t *testing.T) {
	for _, tc := range requestFormatTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("format", tc.format)
		assert.Equal(t, tc.format, r.Format())
	}
}

func TestRequestClass(t *testing.T) {
	for _, tc := range requestClassTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		assert.Equal(t, tc.class, r.Class())
	}
}

func TestRequestReferenceName(t *testing.T) {
	for _, tc := range requestReferenceNameTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("referenceName", tc.referenceName)
		assert.Equal(t, tc.referenceName, r.ReferenceName())
	}
}

func TestRequestStart(t *testing.T) {
	for _, tc := range requestStartTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("start", tc.start)
		assert.Equal(t, tc.start, r.Start())
	}
}

func TestRequestEnd(t *testing.T) {
	for _, tc := range requestEndTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("end", tc.end)
		assert.Equal(t, tc.end, r.End())
	}
}

func TestRequestHtsgetBlockClass(t *testing.T) {
	for _, tc := range requestHtsgetBlockClassTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("HtsgetBlockClass", tc.blockClass)
		assert.Equal(t, tc.blockClass, r.HtsgetBlockClass())
	}
}

func TestRequestHtsgetBlockID(t *testing.T) {
	for _, tc := range requestHtsgetBlockIDTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("HtsgetBlockId", tc.blockID)
		assert.Equal(t, tc.blockID, r.HtsgetBlockID())
	}
}

func TestRequestHtsgetNumBlocks(t *testing.T) {
	for _, tc := range requestHtsgetNumBlocksTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("HtsgetNumBlocks", tc.numBlocks)
		assert.Equal(t, tc.numBlocks, r.HtsgetNumBlocks())
	}
}

func TestRequestFilePath(t *testing.T) {
	for _, tc := range requestHtsgetFilePathTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("HtsgetFilePath", tc.filePath)
		assert.Equal(t, tc.filePath, r.HtsgetFilePath())
	}
}

func TestRequestRange(t *testing.T) {
	for _, tc := range requestRangeTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("Range", tc.Range)
		assert.Equal(t, tc.Range, r.Range())
	}
}

func TestRequestFields(t *testing.T) {
	for _, tc := range requestFieldsTC {
		r := NewHtsgetRequest()
		r.AddListParam("fields", tc.fields)
		assert.Equal(t, tc.fields, r.Fields())
	}
}

func TestRequestTags(t *testing.T) {
	for _, tc := range requestTagsTC {
		r := NewHtsgetRequest()
		r.AddListParam("tags", tc.tags)
		assert.Equal(t, tc.tags, r.Tags())
	}
}

func TestRequestNoTags(t *testing.T) {
	for _, tc := range requestNoTagsTC {
		r := NewHtsgetRequest()
		r.AddListParam("notags", tc.notags)
		assert.Equal(t, tc.notags, r.NoTags())
	}
}

func TestRequestHeaderOnlyRequested(t *testing.T) {
	for _, tc := range requestHeaderOnlyRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		assert.Equal(t, tc.exp, r.HeaderOnlyRequested())
	}
}

func TestRequestUnplacedUnmappedReadsRequested(t *testing.T) {
	for _, tc := range requestUnplacedUnmappedReadsRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("referenceName", tc.referenceName)
		assert.Equal(t, tc.exp, r.UnplacedUnmappedReadsRequested())
	}
}

func TestRequestReferenceNameRequested(t *testing.T) {
	for _, tc := range requestReferenceNameRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("referenceName", tc.referenceName)
		assert.Equal(t, tc.exp, r.ReferenceNameRequested())
	}
}

func TestRequestStartRequested(t *testing.T) {
	for _, tc := range requestStartRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("start", tc.start)
		assert.Equal(t, tc.exp, r.StartRequested())
	}
}

func TestRequestEndRequested(t *testing.T) {
	for _, tc := range requestEndRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("end", tc.end)
		assert.Equal(t, tc.exp, r.EndRequested())
	}
}

func TestRequestAllRegionsRequested(t *testing.T) {
	for _, tc := range requestAllRegionsRequestedTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("referenceName", tc.referenceName)
		assert.Equal(t, tc.exp, r.AllRegionsRequested())
	}
}

func TestRequestAllFieldsRequested(t *testing.T) {
	for _, tc := range requestAllFieldsRequestedTC {
		r := NewHtsgetRequest()
		r.AddListParam("fields", tc.fields)
		assert.Equal(t, tc.exp, r.AllFieldsRequested())
	}
}

func TestRequestAllTagsRequested(t *testing.T) {
	for _, tc := range requestAllTagsRequestedTC {
		r := NewHtsgetRequest()
		r.AddListParam("tags", tc.tags)
		r.AddListParam("notags", tc.notags)
		assert.Equal(t, tc.exp, r.AllTagsRequested())
	}
}

func TestRequestConstructDataEndpointURL(t *testing.T) {

	for _, tc := range requestConstructDataEndpointURLTC {
		request := NewHtsgetRequest()
		request.SetEndpoint(tc.endpoint)
		request.AddScalarParam("id", tc.id)
		request.AddScalarParam("class", tc.class)
		request.AddScalarParam("referenceName", tc.referenceName)
		request.AddScalarParam("start", tc.start)
		request.AddScalarParam("end", tc.end)
		request.AddListParam("fields", tc.fields)
		request.AddListParam("tags", tc.tags)
		request.AddListParam("notags", tc.notags)
		ep, _ := request.ConstructDataEndpointURL()
		assert.Equal(t, tc.exp, ep.String())
	}
}

func TestRequestGetDataSourceRegistry(t *testing.T) {
	for _, tc := range requestDataSourceRegistryTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		registry := r.GetDataSourceRegistry()
		assert.Equal(t, tc.expSource0Pattern, registry.Sources[0].Pattern)
		assert.Equal(t, tc.expSource0Path, registry.Sources[0].Path)
	}
}

func TestRequestGetServiceInfo(t *testing.T) {
	for _, tc := range requestServiceInfoTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		si := r.GetServiceInfo()
		assert.Equal(t, tc.expDatatype, si.HtsgetExtension.Datatype)
	}
}

*/
