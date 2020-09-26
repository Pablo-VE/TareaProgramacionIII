package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	response, err := http.Get("http://localhost:8989/")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Funciono xD")
		fmt.Println(string(data))
	}
}

type usuarioDTO struct {
	id                 int64
	nombreCompleto     string
	cedula             string
	passwordEncriptado string
	estado             bool
	fechaRegistro      time.Time
	fechaModificacion  time.Time
	esJefe             bool
	departamento       departamentoDTO
}

type departamentoDTO struct {
	id                int64
	nombre            string
	estado            bool
	fechaRegistro     time.Time
	fechaModificacion time.Time
}

type tramiteTipoDTO struct {
	id                int64
	descripcion       string
	estado            bool
	departamento      departamentoDTO
	fechaRegistro     time.Time
	fechaModificacion time.Time
}

type variacionDTO struct {
	id            int64
	grupo         int
	descripcion   string
	estado        bool
	fechaRegistro time.Time
	tramitesTipos tramiteTipoDTO
}

type requisitoDTO struct {
	id            int64
	descripcion   string
	estado        bool
	fechaRegistro time.Time
	variaciones   variacionDTO
}

type clienteDTO struct {
	id                 int64
	nombreCompleto     string
	cedula             string
	telefono           string
	direccion          string
	estado             bool
	fechaRegistro      time.Time
	fechaModificacion  time.Time
	passwordEncriptado string
}

type tramiteRegistradoDTO struct {
	id              int64
	tramitesTiposID int
	clienteID       clienteDTO
}

type requisitoPresentadoDTO struct {
	id                  int64
	fechaRegistro       time.Time
	tramiteRegistradoID tramiteRegistradoDTO
	requisitoID         requisitoDTO
}

type tramiteCambioEstadoDTO struct {
	id                    int64
	usuarioID             usuarioDTO
	tramitesRegistradosID tramiteRegistradoDTO
	tramiteEstadoID       tramiteEstadoDTO
	fechaRegistro         time.Time
}

type tramiteEstadoDTO struct {
	id               int64
	nombre           string
	descripcion      string
	estadosSucesores string
}
