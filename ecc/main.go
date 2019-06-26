package main

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	runMode  string
	randKey  string
	randSign string
	prk      *ecdsa.PrivateKey
	puk      ecdsa.PublicKey
	curve    elliptic.Curve
)

func main() {
	var err error
	if err != nil {
		return
	}
	randSign = "aaa"
	if len(randSign) == 0 {
		return
	}
	randKey = "aaaa"
	if len(randKey) == 0 {
		return
	}
	//根据rand长度，使用相应的加密椭圆参数
	length := len([]byte(randKey))
	if length < 224/8 {
		fmt.Println("The length of Rand Key is too small, Crypt init failed, Please reset it again !")
		return
	}
	if length >= 521/8+8 {
		fmt.Println("Rand length =", length, "Using 521 level !")
		curve = elliptic.P521()
	} else if length >= 384/8+8 {
		fmt.Println("Rand length =", length, "Using 384 level !")
		curve = elliptic.P384()
	} else if length >= 256/8+8 {
		fmt.Println("Rand length =", length, "Using 256 level !")
		curve = elliptic.P256()
	} else if length >= 224/8+8 {
		fmt.Println("Rand length =", length, "Using 244 level !")
		curve = elliptic.P224()
	}
	//创建密匙对
	prk, err = ecdsa.GenerateKey(curve, strings.NewReader(randKey))
	if err != nil {
		fmt.Println("Crypt init fail,", err, " need = ", curve.Params().BitSize)
		return
	}
	puk = prk.PublicKey
}

//Encrypt 对Text进行加密，返回加密后的字节流
func Sign(text string) (string, error) {
	r, s, err := ecdsa.Sign(strings.NewReader(randSign), prk, []byte(text))
	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return hex.EncodeToString(b.Bytes()), nil
}

//解密
func getSign(text, byterun []byte) (rint, sint big.Int, err error) {
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode =", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

//Verify 对密文和明文进行匹配校验
func Verify(text, passwd string) (bool, error) {
	byterun, err := hex.DecodeString(passwd)
	if err != nil {
		return false, err
	}
	rint, sint, err := getSign([]byte(text), byterun)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(&puk, []byte(text), &rint, &sint)
	return result, nil
}
