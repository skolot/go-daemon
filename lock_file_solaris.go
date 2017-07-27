// +build solaris

package daemon

import (
	"io"
	"fmt"
	"syscall"
)

func lockFile(fd uintptr) error {
	var fl syscall.Flock_t

	fl.Start = 0;
	fl.Len = 0;
        fl.Whence = io.SeekStart
	fl.Type = syscall.F_WRLCK
	cmd := syscall.F_SETLK

	err := syscall.FcntlFlock(fd, cmd, &fl) 

	if err == syscall.EWOULDBLOCK {
		err = ErrWouldBlock
	}

	return err
}

func unlockFile(fd uintptr) error {
        var fl syscall.Flock_t
        var cmd int

	fl.Type = syscall.F_UNLCK
	cmd = syscall.F_SETLK

	err := syscall.FcntlFlock(fd, cmd, &fl)

	if err == syscall.EWOULDBLOCK {
		err = ErrWouldBlock
	}

	return err
}

const pathMax = 0x1000

func getFdName(fd uintptr) (name string, err error) {
	path := fmt.Sprintf("/proc/self/path/%d", int(fd))
	// We use predefined pathMax const because /proc directory contains special files
	// so that unable to get correct size of pseudo-symlink through lstat.
	// please see notes and example for readlink syscall:
	// http://man7.org/linux/man-pages/man2/readlink.2.html#NOTES
	buf := make([]byte, pathMax)
	var n int
	if n, err = syscall.Readlink(path, buf); err == nil {
		name = string(buf[:n])
	}
	return
}
