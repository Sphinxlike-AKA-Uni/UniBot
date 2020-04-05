#!/bin/bash
function RunUni
{
	./Uni -config ../UniConfig.inf
}

#go build -gccgoflags "-L /lib64 -l pthread" Uni.go
go build Uni.go
if [ $? == 0 ]; then
	RunUni
	if [ $? == 1 ]; then
		bash RunBash.sh
	fi
fi
