BCFTOOLS_VERSION="1.10.2"

wget https://github.com/samtools/bcftools/releases/download/${BCFTOOLS_VERSION}/bcftools-${BCFTOOLS_VERSION}.tar.bz2
tar -xjf bcftools-${BCFTOOLS_VERSION}.tar.bz2
cd bcftools-${BCFTOOLS_VERSION}
./configure --prefix=`pwd`
make
make install
export PATH="`pwd`/bin:${PATH}"
cd ..
