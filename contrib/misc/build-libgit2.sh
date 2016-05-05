#!/usr/bin/env sh
#
# Script used to build and install libgit2 on Travis CI
#

set -xe

LIBGIT2_VERSION="0.24.0"
wget -O libgit2-${LIBGIT2_VERSION}.tar.gz https://github.com/libgit2/libgit2/archive/v${LIBGIT2_VERSION}.tar.gz
tar -xzvf libgit2-${LIBGIT2_VERSION}.tar.gz
cd libgit2-${LIBGIT2_VERSION}
mkdir build && cd build
cmake -DBUILD_CLAR=OFF .. && make && sudo make install
sudo ldconfig
