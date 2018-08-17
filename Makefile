.DEFAULT_GOAL := build-all

export GO15VENDOREXPERIMENT=1

build-all: lvm-scheduler

lvm-scheduler:
	go build -i -o ./cmd/lvm-scheduler/lvm-scheduler ./cmd/lvm-scheduler/
