
SRCS = \
	main.go \
	parsescore.go \
	entry.go \
	root.go \
	dispatch.go \
	sync.go \
	walk.go \
	peer.go \
	conn.go

vtsync: $SRCS
	go build -o vtsync $SRCS
