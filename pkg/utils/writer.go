package utils

import (
    "bytes"
    "log"
)

type DefaultGinWriter struct{}

func(d *DefaultGinWriter) Write(p []byte) (n int, err error){
    b := bytes.NewBuffer(p)
    log.Print(b.String())
    return b.Len(), nil
}

//func demo(){
//    a := []int{}
//    for i:=0; i< 10; i++ {
//        a = append(a, i)
//    }
//    fmt.Println(a)
//}