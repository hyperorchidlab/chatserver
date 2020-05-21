package chatcrypt

import (
	"bytes"
	"crypto/ed25519"

	"errors"
)

func GenGroupAesKey(mainPriv ed25519.PrivateKey, pubkeys [][]byte) (aes []byte, groupKeys [][]byte, err error) {
	if len(pubkeys) <= 0 {
		return
	}

	derivePub := mainPriv.Public()

	for i := 0; i < len(pubkeys); i++ {
		if bytes.Compare(derivePub.(ed25519.PublicKey), pubkeys[i]) == 0 {
			length := len(pubkeys)
			pubkeys[i] = pubkeys[length-1]
			pubkeys = pubkeys[:length-1]
			break
		}
	}

	r := InsertionSortDArray(pubkeys)
	priv := mainPriv

	groupKeys = append(groupKeys, derivePub.(ed25519.PublicKey))

	var pub ed25519.PublicKey

	for i := 0; i < len(r); i++ {
		aes, err = GenerateAesKey(r[i], priv)
		if err != nil {
			return
		}

		if i == len(r)-1 {
			break
		}

		pub, priv = DeriveKey(aes)
		groupKeys = append(groupKeys, pub)
	}

	return
}

func DeriveGroupKey(priv ed25519.PrivateKey, groupPKs [][]byte, pubkeys [][]byte) (aes []byte, err error) {
	derivePub := priv.Public().(ed25519.PublicKey)

	for i := 0; i < len(pubkeys); i++ {
		if bytes.Compare(groupPKs[0], pubkeys[i]) == 0 {
			length := len(pubkeys)
			pubkeys[i] = pubkeys[length-1]
			pubkeys = pubkeys[:length-1]
			break
		}
	}

	if len(groupPKs) != len(pubkeys) {
		return nil, errors.New("pubkeys errors")
	}

	grpidx := -1

	r := InsertionSortDArray(pubkeys)

	for i := 0; i < len(r); i++ {
		if bytes.Compare(derivePub, r[i]) == 0 {
			grpidx = i
		}
	}

	if grpidx == -1 {
		return nil, errors.New("pubkey not found")
	}

	for i := grpidx; i < len(r); i++ {
		var pk ed25519.PublicKey
		if i == grpidx {
			pk = groupPKs[i]

		} else {
			pk = r[i]
		}
		aes, err = GenerateAesKey(pk, priv)
		if err != nil {
			return
		}

		if i == len(r)-1 {
			break
		}

		_, priv = DeriveKey(aes)

	}

	return
}

func InsertionSortDArray(arr [][]byte) [][]byte {

	r := make([][]byte, 0)
	r = append(r, arr[0])

	for i := 1; i < len(arr); i++ {
		flag := false
		for j := 0; j < len(r); j++ {
			if bytes.Compare(r[j], arr[i]) > 0 {
				if j == 0 {
					r1 := make([][]byte, 0)
					r1 = append(r1, arr[i])
					r1 = append(r1, r...)
					r = r1
				} else {
					r1 := make([][]byte, 0)
					r2 := r[:j]
					r3 := r[j:]
					r1 = append(r1, r2...)
					r1 = append(r1, arr[i])
					r1 = append(r1, r3...)
					r = r1
				}
				flag = true
				break
			}
		}
		if !flag {
			r = append(r, arr[i])
		}

	}

	return r
}
