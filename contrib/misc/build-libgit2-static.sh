#!/usr/bin/env sh
#
# Script used to build a static libgit2 library
#

LIBGIT2_VERSION="0.24.0"
wget -O libgit2-${LIBGIT2_VERSION}.tar.gz https://github.com/libgit2/libgit2/archive/v${LIBGIT2_VERSION}.tar.gz
tar -xzvf libgit2-${LIBGIT2_VERSION}.tar.gz
cd libgit2-${LIBGIT2_VERSION}
mkdir {build,install} && cd build
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_C_FLAGS=-fPIC \
      -DCMAKE_INSTALL_PREFIX=../install \
      ..

cmake --build . --target install

