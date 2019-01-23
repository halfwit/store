# Storage

## Overview

Slightly modified Plan9 plumber, used to store URL endpoints, using plan9 plumber rule files.

## Usage

`store <url>`

## Installation

Requires Plan9, hasn't been tested with plan9port yet
Requires Go to build `store` binary:

```
go get github.com/halfwit/storage/store
go install github.com/halfwit/storage/store

```

## Caveats

This differs from Plumber's rules handling in that it introduces a new verb and a new action.

`output` - this sets the file that a subsequent `fetch`. 

```
# this will attempt to create $home/doc if it doesn't exist
type matches application/$document
data matches 'https?://([^ ]/)+([^ ]+)'
arg isdir $home/doc
output is $home/doc/$2
store fetch $0

type matches application/$document
data matches 'https?://[^ ]/+[^ ]+'
data matches 'https://somesite/([^ ]+)'
data set 'https://somemirrorsite/$1'
output is $data
store fetch $data

```
