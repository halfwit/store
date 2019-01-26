# Store

## Overview

This utility is an extension to plan9's plumber to handle storage of remote resources, referenced by URL.

## Usage

`store <url>`

## Installation

Requires Go to build `store` binary:

```
go get github.com/halfwit/store
go install github.com/halfwit/store

```

## Setup

Add matches to your plumber's rulefile, with the `type matching` `type is` set to the content-type you wish to target. 
(Those wishing to find the content type of a resource out can inspect the header of the URL, or use a tool such as [content-type](https://github.com/halfwit/content-type) )

```

document='(pdf|PDF|ps|PS|djvu|epub)'
# 'store https://some.domain/path/to/file.pdf'
type matches application/$document
data matches 'https?://([^ ]/)+([^ ]+)' // validate url
dst is store
arg isdir /usr/halfwit/doc // make sure our document folder exists
data set $dir/$2 // set to 'file.pdf'
plumb start rc -c 'hget -o '$data' '$0'

type matches text/html
dst is store
data matches 'https://(github.com/[^ ]+)'
arg isdir /usr/halfwit/src/$1
plumb start rc -c 'cd '$dir' && git clone '$0

```

Alternatively, you can use [storage](https://github.com/halfwit/storage), which listens for plumb messages to `store` as a target.
The resulting rules file is somewhat simplified

```

type	is	application/pdf
data	matches	'https?//[^ ]+'
data	matches	'https?//[^ ]/([a-zA-Z0-9-_/\])'
attr	add	filename=/usr/halfwit/docs/$1
plumb	to	store
plumb	client	store	$0

type	matches	application/$document
data	matches	'https?://[^ ]+'
data	matches 'https://github.com/([^ ]+)'
attr	add	filepath=/usr/halfwit/src/$1
plumb	to	store
plumb	client	store	$0

```
