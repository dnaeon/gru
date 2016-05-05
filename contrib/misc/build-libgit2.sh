#!/usr/bin/env sh

git clone https://github.com/libgit2/libgit2.git
cd libgit2/
mkdir build && cd build
cmake ..
cmake --build . --target install
