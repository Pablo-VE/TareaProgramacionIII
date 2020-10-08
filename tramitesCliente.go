package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Pablo-VE/TareaProgramacionIII/dto"
	"github.com/Pablo-VE/TareaProgramacionIII/util"
)

var usuarioLogeado dto.AuthenticationResponse

const url string = "http://localhost:8989/"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/tramitesRegistrados", login)
	http.HandleFunc("/buscarPorId", buscarPorID)
	http.HandleFunc("/buscarPorCedula", buscarPorCedula)
	http.HandleFunc("/buscarPorEstado", buscarPorEstado)
	http.HandleFunc("/buscarPorFecha", buscarPorFecha)
	http.HandleFunc("/TramitesRegistrados", limpiar)
	http.HandleFunc("/irTramite", irTramite)
	http.HandleFunc("/guardarEstado", guardarEstado)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

var tramiteRegistradoEnCuestion dto.TramiteRegistradoDTO

func guardarEstado(w http.ResponseWriter, r *http.Request) {
	festado := r.FormValue("cbxEstado")
	fdescripcion := r.FormValue("txtDescripcion")

	var tramiteEstado dto.TramiteEstadoDTO

	if festado == "Revisar" {
		tramiteEstado.Nombre = "En revisión"
	}
	if festado == "Anular" {
		tramiteEstado.Nombre = "Anulado"
	}
	if festado == "Finalizar" {
		tramiteEstado.Nombre = "Finalizado"
	}
	if festado == "Entregar" {
		tramiteEstado.Nombre = "Entregado"
	}
	tramiteEstado.Descripcion = fdescripcion

	j, err := json.Marshal(tramiteEstado)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url+"tramites_estados/", bytes.NewBuffer(j))
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
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 201 {
		//irTramite(w, r)
	} else {

		json.Unmarshal(body, &tramiteEstado)

		var tramiteCambioEstado dto.TramiteCambioEstadoDTO
		tramiteCambioEstado.TramiteEstadoID = tramiteEstado
		tramiteCambioEstado.TramitesRegistradosID = tramiteRegistradoEnCuestion
		tramiteCambioEstado.UsuarioID = usuarioLogeado.Usuario

		j, err := json.Marshal(tramiteCambioEstado)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		req, err := http.NewRequest("POST", url+"tramites_cambio_estado/", bytes.NewBuffer(j))
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
		body, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode != 201 {
			//irTramite(w, r)
		} else {
			json.Unmarshal(body, &tramiteCambioEstado)
			limpiar(w, r)
		}
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fid := r.FormValue("cedula")
	fpassword := r.FormValue("password")
	ar := dto.AuthenticationRequest{Cedula: fid, Password: fpassword}
	j, err := json.Marshal(ar)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url+"autenticacion/login", bytes.NewBuffer(j))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		tpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &usuarioLogeado)

		tramitesDTO := findAllTramitesRegistrados()
		tramitesTable := crearDatosTable(tramitesDTO)

		d := struct {
			Usuario  string
			Tramites []util.DatoTramitesTable
		}{
			Usuario:  usuarioLogeado.Usuario.NombreCompleto,
			Tramites: tramitesTable,
		}

		tpl.ExecuteTemplate(w, "tramites.html", d)
	}
}

func buscarPorID(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("txtID")

	nid, err := strconv.ParseInt(fid, 10, 64)
	if err != nil {

	}
	tramiteDTO := findTramitesRegistradosByID(nid)
	var tramitesDTO []dto.TramiteRegistradoDTO
	tramitesDTO = append(tramitesDTO, tramiteDTO)
	tramitesTable := crearDatosTable(tramitesDTO)

	d := struct {
		Usuario  string
		Tramites []util.DatoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramitesTable,
	}

	tpl.ExecuteTemplate(w, "tramites.html", d)

}

func buscarPorCedula(w http.ResponseWriter, r *http.Request) {
	fcedula := r.FormValue("txtCedula")

	tramitesDTO := findTramitesRegistradosByCedulaCliente(fcedula)
	tramitesTable := crearDatosTable(tramitesDTO)

	d := struct {
		Usuario  string
		Tramites []util.DatoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramitesTable,
	}

	tpl.ExecuteTemplate(w, "tramites.html", d)

}

func buscarPorEstado(w http.ResponseWriter, r *http.Request) {
	festado := r.FormValue("cbxEstado")

	tramitesDTO := findAllTramitesRegistrados()
	tramitesTable := crearDatosTable(tramitesDTO)

	var tramites []util.DatoTramitesTable
	for i := 0; i < len(tramitesTable); i++ {
		if tramitesTable[i].Estado == festado {
			tramites = append(tramites, tramitesTable[i])
		}
	}

	d := struct {
		Usuario  string
		Tramites []util.DatoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramites,
	}

	tpl.ExecuteTemplate(w, "tramites.html", d)

}

func buscarPorFecha(w http.ResponseWriter, r *http.Request) {
	ffecha := r.FormValue("txtFecha")
	ffecha2 := strings.Split(ffecha, "-")
	fecha := strings.Join(ffecha2, "")
	now := time.Now()
	date := now.Format("20060102")
	fmt.Println(date)
	date = now.Format("2006-01-02")
	date2, err := time.Parse("20060102", fecha)
	if err == nil {

	}
	var fecha1 time.Time
	fecha1 = date2
	fechaString := fecha1.Format("Mon Jan _2 15:04:05 2006")
	tramitesDTO := findAllTramitesRegistrados()
	tramitesTable := crearDatosTable(tramitesDTO)
	var tramites []util.DatoTramitesTable
	for i := 0; i < len(tramitesTable); i++ {
		if getFecha(tramitesTable[i].FechaRegistro) == getFecha(fechaString) {
			tramites = append(tramites, tramitesTable[i])
		}
	}
	d := struct {
		Usuario  string
		Tramites []util.DatoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramites,
	}
	tpl.ExecuteTemplate(w, "tramites.html", d)
}
func getFecha(f string) (fechas string) {
	layout := "Mon Jan _2 15:04:05 2006"
	fecha, err := time.Parse(layout, f)

	if err != nil {
		fmt.Println(err)
	}
	fechas = fecha.Format("01-02-2006")
	return fechas

}

func limpiar(w http.ResponseWriter, r *http.Request) {

	tramitesDTO := findAllTramitesRegistrados()
	tramitesTable := crearDatosTable(tramitesDTO)

	d := struct {
		Usuario  string
		Tramites []datoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramitesTable,
	}

	tpl.ExecuteTemplate(w, "tramites.html", d)

}

func irTramite(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("txtVerID")
	tid, err := strconv.ParseInt(fid, 10, 64)
	if err != nil {

	}
	fmt.Println(tid)

	tramiteRegistradoEnCuestion = findTramitesRegistradosByID(tid)

	tramiteRegistradoView := getTramiteRegistradoView(tramiteRegistradoEnCuestion)

	tramiteRegistradoView.Usuario = usuarioLogeado.Usuario.NombreCompleto

	tpl.ExecuteTemplate(w, "TramiteRegistrado.html", tramiteRegistradoView)

}

func findAllTramitesRegistrados() (tramitesregistrados []dto.TramiteRegistradoDTO) {
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

func findTramitesRegistradosByID(idTR int64) (tramiteregistrado dto.TramiteRegistradoDTO) {
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

func findTramitesRegistradosByCedulaCliente(cedula string) (tramiteregistrados []dto.TramiteRegistradoDTO) {
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

func findTipoTramiteByID(idTT int64) (tramitetipo dto.TramiteTipoDTO) {
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

func findNotasByTramiteRegistradoID(idTR int64) (notas []dto.NotaDTO) {
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

func findRequisitosPresentadosByTramiteRegistradoID(idTR int64) (requisitosR []dto.RequisitoPresentadoDTO) {
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

func findTramiteCambioEstadoByTramiteRegistradoID(idTR int64) (tramitesCambio []dto.TramiteCambioEstadoDTO) {
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
