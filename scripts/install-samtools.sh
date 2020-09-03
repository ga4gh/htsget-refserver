SAMTOOLS_VERSION="1.9"

wget https://github.com/samtools/samtools/releases/download/${SAMTOOLS_VERSION}/samtools-${SAMTOOLS_VERSION}.tar.bz2
tar -xjf samtools-${SAMTOOLS_VERSION}.tar.bz2
cd samtools-${SAMTOOLS_VERSION}
pwd
./configure --prefix=`pwd`
make
make install
export PATH="`pwd`/bin:${PATH}"
