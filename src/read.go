package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const SYS_PROCESS_VM_READV = 310

type iovec struct {
	iov_base unsafe.Pointer
	iov_len uint64
}

var _zero uintptr
func ProcessVMReadv(pid int, local []iovec, remote []iovec, flags uint) (n int, err error) {
	var p0 unsafe.Pointer

	if len(local) > 0 {
		//local.iov_base
		p0 = local[0].iov_base
	} else {
		p0 = unsafe.Pointer(&_zero)
	}

	var p1 unsafe.Pointer
	
	if len(remote) > 0 {
		//remote.iov_base
		p1 = unsafe.Pointer(&remote[0])
	} else {
		p1 = unsafe.Pointer(&_zero)
	}

	r0, _, err := syscall.Syscall6(
		SYS_PROCESS_VM_READV,
		uintptr(pid),
		uintptr(p0),
		uintptr( len(local) ),
		uintptr(p1),
		uintptr( len(remote) ),
		uintptr(flags),
	)

	n = int(r0)
	return
}

func remote_process_read(pid int, address unsafe.Pointer, buffer []byte, len uint64) {

	var local iovec = iovec{}
	local.iov_base = unsafe.Pointer(&buffer)
	local.iov_len = len

	var remote iovec = iovec{}
	remote.iov_base = address
	remote.iov_len = local.iov_len

	iovarr1 := []iovec{local}
	iovarr2 := []iovec{remote}

	
	n, _ := ProcessVMReadv(pid, iovarr1, iovarr2, 0)
	fmt.Println("Read", n, "bytes from", address)
	
	var i uint64
	for i = 0; i < len; i++ {
		fmt.Printf("%02X\n", buffer[i])
	}
}

//enum
type arg int
const (
	APP_NAME arg = iota
	ARG_PID
	ARG_ADDR
	ARG_BYTES
	ARG_COUNT
)

func main() {
	
	if len(os.Args) != int(ARG_COUNT) {
		log.Fatal("Invalid Syntax")
	}

	pid, err := strconv.Atoi(os.Args[ARG_PID])
	if err != nil {
		log.Fatal(err)
	}
	
	address, err := strconv.ParseInt(os.Args[ARG_ADDR], 16, 0)
	if err != nil {
		log.Fatal(err)
	}

	len, err := strconv.ParseUint(os.Args[ARG_BYTES], 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, len)
	
	remote_process_read(pid, unsafe.Pointer(uintptr(address)), buffer, len)
}