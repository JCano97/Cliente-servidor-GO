package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)
type Proceso struct {
	idProceso int
	contadorProceso uint64
	enviado bool
}
func ejecucionProceso(proceso *Proceso) {
	//fmt.Println("entra a ejecutar procesos")
	for {
		time.Sleep(time.Millisecond * 500)
		if proceso.enviado != true {
			fmt.Println(proceso.idProceso, " : ", proceso.contadorProceso)
		} else {
			return //se termina el proceso
		}
		proceso.contadorProceso++
	}
}
func main(){
	var proceso Proceso
	var input string
	proceso.enviado = true
	//CLIENTE
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = gob.NewEncoder(c).Encode(proceso)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c.Close()
	//RECIBE DE SERVIDOR
	s, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for {
		c, err = s.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		err = gob.NewDecoder(c).Decode(&proceso)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		s.Close()
		proceso.enviado = false
		go ejecucionProceso(&proceso)

		fmt.Scanln(&input)

		proceso.enviado = false
		//Regresamos proceso a servidor
		c, err := net.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = gob.NewEncoder(c).Encode(proceso)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println("Proceso devuelto, presione ENTER para terminar...")
		break
	}
}