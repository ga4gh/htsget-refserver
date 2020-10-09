package htsserver

import "github.com/ga4gh/htsget-refserver/internal/htsrequest"

func noAfterSetup(handler *requestHandler) error {
	return nil
}

func addRegionFromQueryString(handler *requestHandler) error {
	htsReq := handler.HtsReq
	if htsReq.ReferenceNameRequested() {
		region := htsrequest.NewRegion()
		region.SetReferenceName(htsReq.GetReferenceName())
		region.SetStart(htsReq.GetStart())
		region.SetEnd(htsReq.GetEnd())
		htsReq.AddRegion(region)
	}
	return nil
}
