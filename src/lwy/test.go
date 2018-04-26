package main

import (
	"fmt"
)

type rename [4]byte

func main() {
	{ //test bytes copy
		buf := make([]byte, 100)
		src1 := []byte("abcefg")
		src2 := []byte("HJKLMN")
		//src3 := []int{1, 2, 3, 4, 5}
		src4 := rename{'1', '2', '3', '4'}
		fmt.Println(src4)
		copy(buf[1:3], src1)
		copy(buf[20:22], src2)
		copy(buf[30:40], src4[:])

		fmt.Println(buf)

	}
	/*
		{ //test struct compare

			var a, b [64]byte
			for i, v := range []byte("abcdefgHIJKLMN") {
				a[i] = v
				b[i] = v
			}
			fmt.Println(a == b)
			type Stru struct {
				A int
				B string
			}
			c := Stru{12, "abc"}
			d := Stru{12, "abc"}
			fmt.Println(c, d, c == d)
			e := &Stru{12, "abc"}
			fmt.Println(e.A, (*e).A)

		}
	*/
	{ //test file seek
		/*
			file, err := os.OpenFile("test2.dat", os.O_CREATE|os.O_RDWR, 0664)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			seekRet, err := file.Seek(1000000, os.SEEK_SET)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(seekRet)
			writeRet, err := file.Write([]byte("Hello Word"))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(writeRet)

			seekRet, err = file.Seek(-5, os.SEEK_END)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(seekRet)
			b := make([]byte, 10)
			readRet, err := file.Read(b)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(b)

			seekRet, err = file.Seek(2000000, os.SEEK_SET)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(seekRet)

			readRet, err = file.Read(b)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(readRet)
		*/
	}
}
