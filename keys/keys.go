// Copyright 2015 The sero.cash Authors
// This file is part of the sero.cash library.
//
// The libzero library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The libzero library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the libzero library. If not, see <http://www.gnu.org/licenses/>.

package keys

/*
#cgo CFLAGS: -I ../czero/include

#cgo LDFLAGS: -L ../czero/lib
#cgo LDFLAGS: -lczerod

#include "zero.h"
*/
import "C"

import (
	"crypto/rand"
	"errors"
	"unsafe"

	"github.com/sero-cash/go-sero/crypto/sha3"
)

func logBytes(bytes []byte) {
	C.zero_log_bytes(
		(*C.uchar)(unsafe.Pointer(&bytes[0])),
		(C.int)(len(bytes)),
	)
	return
}

func Seed2Tk(seed *Uint256) (tk Uint512) {
	C.zero_seed2tk(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		(*C.uchar)(unsafe.Pointer(&tk[0])),
	)
	return
}

func Seed2Addr(seed *Uint256) (addr Uint512) {
	C.zero_seed2pk(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		(*C.uchar)(unsafe.Pointer(&addr[0])),
	)
	return
}

func RandUint512() (hash Uint512) {
	rand.Read(hash[:])
	return
}

func RandUint256() (hash Uint256) {
	rand.Read(hash[:])
	return
}

func RandUint128() (hash Uint128) {
	rand.Read(hash[:])
	return
}

func Addr2PKr(addr *Uint512, r *Uint256) (pkr PKr) {
	if r == nil {
		t := RandUint256()
		r = &t
	} else {
		if (*r) == Empty_Uint256 {
			panic("gen pkr, but r is empty")
		}
	}

	C.zero_pk2pkr(
		(*C.uchar)(unsafe.Pointer(&addr[0])),
		(*C.uchar)(unsafe.Pointer(&r[0])),
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
	)
	return
}

func HashPKr(pkr *PKr) (ret Uint256) {
	d := sha3.NewKeccak256()
	d.Write(pkr[:])
	copy(ret[:], d.Sum(nil))
	return
}

const PROOF_WIDTH = 131

type LICr [PROOF_WIDTH]byte

func Addr2PKrAndLICr(addr *Uint512) (pkr PKr, licr LICr, ret bool) {
	r := C.zero_pk2pkr_and_licr(
		(*C.uchar)(unsafe.Pointer(&addr[0])),
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
		(*C.uchar)(unsafe.Pointer(&licr[0])),
	)
	if r == C.char(0) {
		ret = true
	} else {
		ret = false
	}
	return
}

func CheckLICr(pkr *PKr, licr *LICr) bool {
	r := C.zero_check_licr(
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
		(*C.uchar)(unsafe.Pointer(&licr[0])),
	)
	if r == C.char(0) {
		return true
	} else {
		return false
	}
}

func IsMyPKr(tk *Uint512, pkr *PKr) (succ bool) {
	ret := C.zero_ismy_pkr(
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
		(*C.uchar)(unsafe.Pointer(&tk[0])),
	)
	if ret == C.char(0) {
		succ = true
		return
	} else {
		succ = false
		return
	}
}

func SignPKr(seed *Uint256, data *Uint256, pkr *PKr) (sign Uint512, e error) {
	C.zero_sign_pkr(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
		(*C.uchar)(unsafe.Pointer(&sign[0])),
	)
	if sign == Empty_Uint512 {
		e = errors.New("SignOAddr: sign is empty")
		return
	} else {
		return
	}
}

func VerifyPKr(data *Uint256, sign *Uint512, pkr *PKr) bool {
	ret := C.zero_verify_pkr(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		(*C.uchar)(unsafe.Pointer(&sign[0])),
		(*C.uchar)(unsafe.Pointer(&pkr[0])),
	)
	if ret == C.char(0) {
		return true
	} else {
		return false
	}
}
