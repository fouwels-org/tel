// SPDX-FileCopyrightText: 2015 (c) Quoc-Viet Nguyen
//
// SPDX-License-Identifier: BSD-3-Clause

// +build freebsd openbsd netbsd

package serial

import (
	"syscall"
)

func cfSetIspeed(termios *syscall.Termios, speed uint32) {
	termios.Ispeed = speed
}

func cfSetOspeed(termios *syscall.Termios, speed uint32) {
	termios.Ospeed = speed
}
