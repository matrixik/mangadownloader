# Manga Downloader
A manga downloader written in Go (Golang).

## Features
- Download manga from websites (more to come)
    - http://www.mangareader.net
    - http://mangafox.me
    - http://www.mangahere.com
    - http://mangawall.com
    - http://www.tenmanga.com
- CBZ support
- Parallel download
- Command line

## Binaries
https://www.dropbox.com/sh/zc94728pgccei17/mskEqE4XuM

## Usage
```
Usage:
Pass urls (manga/chapter/page) as argument.

Options: (pass them BEFORE the arguments, Go's "flag" package is not really smart...)
  -cbz=false: CBZ
  -httpretry=5: Http retry
  -out="": Output directory
  -pagedigitcount=4: Page digit count
  -parallelchapter=4: Parallel chapter
  -parallelpage=8: Parallel page
```

Examples:

```
./mangadownloader http://www.mangareader.net/97/gantz.html
./mangadownloader http://mangafox.me/manga/berserk/c134/1.html
```

## Status
[![Build Status](https://travis-ci.org/pierrre/mangadownloader.png?branch=master)](https://travis-ci.org/pierrre/mangadownloader)

## Build
`go build mangadownloader/mangadownloader.go`

or

`go get github.com/pierrre/mangadownloader/mangadownloader` and the binary will be in your `GOPATH/bin`

## Help
- Twitter: @pierredurand87
- Github issue

## TODO
- More services
    - animea
    - unixmanga
    - goodmanga
    - mangacraze
    - mangago
    - anymanga
    - mangainn
    - mangaeden
    - mangable
    - deliciousmanga
    - mangastream
- Fix download freeze
- Tests
- Documentation
- Sync
- Input file
- Concurrency for input
- Progress
- User agent
- Improve error handling
- Improve http error codes
