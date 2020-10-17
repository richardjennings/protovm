module github.com/richardjennings/proto

go 1.14

require (
	github.com/richardjennings/proto/vm v0.0.0
	github.com/richardjennings/proto/asm v0.0.0
)
replace (
   github.com/richardjennings/proto/vm => ./vm
   github.com/richardjennings/proto/asm => ./asm
)
