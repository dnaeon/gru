#!/usr/bin/env sh

set -xe

git clone https://github.com/libgit2/libgit2.git
cd libgit2/
mkdir build && cd build
cmake -DBUILD_CLAR=OFF .. && make && sudo make install
sudo ldconfig
