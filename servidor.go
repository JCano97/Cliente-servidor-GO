package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var listaProcesos []*Proceso

type Proceso struct {
	idProceso int
	contadorProceso uint64
	enviado bool
}

func servidor() {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		go handleClient(c)
	}
}
func handleClient(c net.Conn) {
	var proceso Proceso
	err := gob.NewDecoder(c).Decode(&proceso)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else {
		if proceso.enviado == true {
			enviarProceso(c)
		}else{ //cuando regrese el proceso del cliente, poner enviado en false
			listaProcesos = append(listaProcesos, &proceso)
			proceso.enviado = false 
			go ejecucionProceso(&proceso)
		}
	}
	c.Close()

}
func enviarProceso(c net.Conn){
	c, err := net.Dial("tcp", ":9998")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(listaProcesos) > 0 {
		procesoParaEnviar := listaProcesos[0]
		procesoParaEnviar.enviado = true

		err = gob.NewEncoder(c).Encode(procesoParaEnviar)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//time.Sleep(time.Millisecond * 250)

		eliminarDelista()
	} else {
		fmt.Println("No hay mas procesos")
	}
	c.Close()
}
func eliminarDelista(){
	listaProcesos[0] = listaProcesos[len(listaProcesos)-1]
	listaProcesos[len(listaProcesos)-1] = nil
	listaProcesos = listaProcesos[:len(listaProcesos)-1]
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
func main() {
	var input string
	for i := 1; i < 6; i++ {
		proceso := Proceso {
			idProceso: i,
			contadorProceso: uint64(0),
			enviado: false,
		}
		listaProcesos = append(listaProcesos, &proceso)
		go ejecucionProceso(&proceso)
		//fmt.Println("entra al for")
	}
	go servidor() 
	fmt.Scanln(&input)
}
