package conexionservidor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Pablo-VE/TareaProgramacionIII/dto"
)

const url string = "http://localhost:8989/"

var usuarioLogeado dto.AuthenticationResponse

//POST funcion para hacer request tipo Post
func POST(direccion string, estructura interface{}) (response []byte) {
	j, err := json.Marshal(estructura)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url+direccion, bytes.NewBuffer(j))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", "bearer "+usuarioLogeado.Jwt)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 {

	} else {
		if res.StatusCode == 200 || res.StatusCode == 201 {
			body, _ := ioutil.ReadAll(res.Body)
			log.Printf("Response: %s", body)
			response = body
		}
	}
	return response
}

//GET funcion para hacer request tipo get
func GET(direccion string) (response []byte) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+direccion, nil)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", "bearer "+usuarioLogeado.Jwt)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		response = body
	}
	return response
}

//Logear es la funcion para hacer el request del logeo
func Logear(fid string, fpassword string) (usuario dto.AuthenticationResponse) {
	ar := dto.AuthenticationRequest{Cedula: fid, Password: fpassword}
	json.Unmarshal(POST("autenticacion/login", ar), &usuarioLogeado)
	return usuarioLogeado
}

//FindAllTramitesRegistrados is ...
func FindAllTramitesRegistrados() (tramitesregistrados []dto.TramiteRegistradoDTO) {
	json.Unmarshal(GET("tramites_registrados/"), &tramitesregistrados)
	return tramitesregistrados
}

//FindTramitesRegistradosByID is ...
func FindTramitesRegistradosByID(idTR int64) (tramiteregistrado dto.TramiteRegistradoDTO) {
	id := strconv.FormatInt(idTR, 10)
	json.Unmarshal(GET("tramites_registrados/"+id), &tramiteregistrado)
	return tramiteregistrado
}

//FindTramitesRegistradosByCedulaCliente is ..
func FindTramitesRegistradosByCedulaCliente(cedula string) (tramiteregistrados []dto.TramiteRegistradoDTO) {
	json.Unmarshal(GET("tramites_registrados/cedula/"+cedula), &tramiteregistrados)
	return tramiteregistrados
}

//FindTipoTramiteByID is ..
func FindTipoTramiteByID(idTT int64) (tramitetipo dto.TramiteTipoDTO) {
	id := strconv.FormatInt(idTT, 10)
	json.Unmarshal(GET("tramites_tipos/"+id), &tramitetipo)
	return tramitetipo
}

//FindNotasByTramiteRegistradoID is ..
func FindNotasByTramiteRegistradoID(idTR int64) (notas []dto.NotaDTO) {
	id := strconv.FormatInt(idTR, 10)
	json.Unmarshal(GET("notas/tramitesRegistrados/"+id), &notas)
	return notas
}

//FindRequisitosPresentadosByTramiteRegistradoID is..
func FindRequisitosPresentadosByTramiteRegistradoID(idTR int64) (requisitosR []dto.RequisitoPresentadoDTO) {
	id := strconv.FormatInt(idTR, 10)
	json.Unmarshal(GET("requisitos_presentados/tramite_registrado/"+id), &requisitosR)
	return requisitosR
}

//FindTramiteCambioEstadoByTramiteRegistradoID is ..
func FindTramiteCambioEstadoByTramiteRegistradoID(idTR int64) (tramitesCambio []dto.TramiteCambioEstadoDTO) {
	id := strconv.FormatInt(idTR, 10)
	json.Unmarshal(GET("tramites_cambio_estado/tramitesRegistrados/"+id), &tramitesCambio)
	return tramitesCambio
}

//CreateTramiteEstado is ...
func CreateTramiteEstado(tramiteEstadoAGuardar dto.TramiteEstadoDTO) (tramiteEstado dto.TramiteEstadoDTO) {
	json.Unmarshal(POST("tramites_estados/", tramiteEstadoAGuardar), &tramiteEstado)
	return tramiteEstado
}

//CreateTramiteCambioEstado is ...
func CreateTramiteCambioEstado(tramiteCambioEstadoAGuardar dto.TramiteCambioEstadoDTO) (tramiteCambioEstado dto.TramiteCambioEstadoDTO) {
	json.Unmarshal(POST("tramites_cambio_estado/", tramiteCambioEstadoAGuardar), &tramiteCambioEstado)
	return tramiteCambioEstado
}
