# Storage

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

A `storage` ruleset is included, which can be called in your lib/plumbing via `include storage`. 
It defines the following:

```

```

Adding additional rules is trivial; simply match `type` to your target mimetype, and continue on with a normal rule.

## Example rules file for plumber

```

# 'store https://some.domain/path/to/file.pdf'
type matches application/$document
dst is store
data matches 'https?://([^ ]/)+([^ ]+)' // validate url
arg isdir /usr/halfwit/doc // make sure our document folder exists
data set $dir/$2 // set to 'file.pdf'
plumb start rc -c 'hget -o '$data' '$0'

type matches application/$document
dst is store
data matches 'https://(github.com/[^ ]+)'
arg isdir /usr/halfwit/src/$1
plumb start rc -c 'cd '$dir' && git clone '$0

```
