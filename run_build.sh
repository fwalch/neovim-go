#!/bin/bash -e

if [ ! -d $TRAVIS_BUILD_DIR/_neovim ]
then
  mkdir -p $TRAVIS_BUILD_DIR/_neovim
  wget -q -O - https://github.com/fwalch/neovim/releases/download/nightly/neovim-x64-nightly.tar.gz \
    | tar xzf - --strip-components=1 -C $TRAVIS_BUILD_DIR/_neovim
fi

pushd $TRAVIS_BUILD_DIR/_cmd/gen_neovim_api/
go get -d -t -v ./... && go build -v ./...
x=`mktemp`
NEOVIM_BIN=$TRAVIS_BUILD_DIR/_neovim/bin/nvim ./gen_neovim_api -g | gofmt > $x
diff $x ../../gen_client_api.go
if [ $? -ne 0 ]
then
  echo "Neovim exposed API differs from committed generated API"
  exit 1
fi
rm $x
popd

NEOVIM_BIN=$TRAVIS_BUILD_DIR/_neovim/bin/nvim go test
NEOVIM_BIN=$TRAVIS_BUILD_DIR/_neovim/bin/nvim go test -race
