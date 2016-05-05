#!/usr/bin/env sh

git clone https://github.com/libgit2/libgit2.git
cd libgit2/
mkdir build && cd build
cmake .. -DCMAKE_INSTALL_PREFIX=../_install -DBUILD_CLAR=OFF
cmake --build . --target install
