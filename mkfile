# order matters
SRCS = \
	vtsync.go \
	parsescore.go \
	entry.go \
	root.go \
	dispatch.go \
	syncmsg.go \
	coord.go

vtsync: $SRCS
	go build $SRCS
