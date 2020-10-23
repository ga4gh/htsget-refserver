CWD=`pwd`
go mod download github.com/ga4gh/htsget-refserver-utils@v1.0.0
cd $HOME/gopath/pkg/mod/github.com/ga4gh/htsget-refserver-utils@v1.0.0
go install
cd ${CWD}
