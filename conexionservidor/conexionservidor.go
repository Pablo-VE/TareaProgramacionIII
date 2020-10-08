package conexionservidor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Pablo-VE/TareaProgramacionIII/dto"
)

const url string = "http://localhost:8989/"

var usuarioLogeado dto.AuthenticationResponse

//FindAllTramitesRegistrados is ...
func FindAllTramitesRegistrados() (tramitesregistrados []dto.TramiteRegistradoDTO) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"tramites_registrados/", nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &tramitesregistrados)
		return tramitesregistrados
	}
	return nil
}

//FindTramitesRegistradosByID is ...
func FindTramitesRegistradosByID(idTR int64) (tramiteregistrado dto.TramiteRegistradoDTO) {
	id := strconv.FormatInt(idTR, 10)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"tramites_registrados/"+id, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &tramiteregistrado)
	}
	return tramiteregistrado
}

//FindTramitesRegistradosByCedulaCliente is ..
func FindTramitesRegistradosByCedulaCliente(cedula string) (tramiteregistrados []dto.TramiteRegistradoDTO) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"tramites_registrados/cedula/"+cedula, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &tramiteregistrados)
	}
	return tramiteregistrados
}

//FindTipoTramiteByID is ..
func FindTipoTramiteByID(idTT int64) (tramitetipo dto.TramiteTipoDTO) {
	id := strconv.FormatInt(idTT, 10)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"tramites_tipos/"+id, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &tramitetipo)
	}
	return tramitetipo

}

//FindNotasByTramiteRegistradoID is ..
func FindNotasByTramiteRegistradoID(idTR int64) (notas []dto.NotaDTO) {
	id := strconv.FormatInt(idTR, 10)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"notas/tramitesRegistrados/"+id, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {

		json.Unmarshal(bodyBytes, &notas)
	}
	return notas
}

//FindRequisitosPresentadosByTramiteRegistradoID is..
func FindRequisitosPresentadosByTramiteRegistradoID(idTR int64) (requisitosR []dto.RequisitoPresentadoDTO) {
	id := strconv.FormatInt(idTR, 10)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"requisitos_presentados/tramite_registrado/"+id, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &requisitosR)
		return requisitosR
	}

	return nil
}

//FindTramiteCambioEstadoByTramiteRegistradoID is ..
func FindTramiteCambioEstadoByTramiteRegistradoID(idTR int64) (tramitesCambio []dto.TramiteCambioEstadoDTO) {
	id := strconv.FormatInt(idTR, 10)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"tramites_cambio_estado/tramitesRegistrados/"+id, nil)
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
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		json.Unmarshal(bodyBytes, &tramitesCambio)
		return tramitesCambio
	}

	return nil

}
