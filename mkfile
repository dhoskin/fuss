# order matters
SRCS = \
	vtsync.go \
	parsescore.go \
	entry.go \
	root.go

vtsync: $SRCS
	go build $SRCS
