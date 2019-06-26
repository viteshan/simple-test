package basic

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os/user"
	"strconv"
	"testing"
)

func TestUint(t *testing.T) {
	var u uint32 = 17
	var s = strconv.FormatUint(uint64(u), 10)
	println(s)

}

func TestByteArr(t *testing.T) {
	x := uint32(40)
	buf := new(bytes.Buffer)
	// for int32, the resulting size of buf will be 4 bytes
	// for int64, the resulting size of buf will be 8 bytes
	err := binary.Write(buf, binary.BigEndian, x)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x\n", buf.Bytes())
	fmt.Printf("%x\n", hex.EncodeToString(buf.Bytes()))
	//fmt.Printf("%s \n", string(buf.Bytes()))
}

func TestHome(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)
}
func TestUint64(t *testing.T) {
	u := uint64(12)
	u = u - 13
	println(u)
}

func TestBinary(t *testing.T) {
	u := uint64(111111111)
	var hashbuf [8]byte
	binary.BigEndian.PutUint64(hashbuf[:], u)

	i := binary.BigEndian.Uint64(hashbuf[:])
	println(i)
	if i != u {
		t.Error("not equals")
	}
}

func TestBigInt(t *testing.T) {
	diff := big.NewInt(20)
	diff.SetUint64(20)
	zero := big.NewInt(0)
	println(diff.Cmp(zero) > 0)
	if diff.Cmp(zero) > 0 {
		println("---------")
	}
	println()
}
