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
)

var usuarioLogeado AuthenticationResponse

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

var tramiteRegistradoEnCuestion TramiteRegistradoDTO

func guardarEstado(w http.ResponseWriter, r *http.Request) {
	festado := r.FormValue("cbxEstado")
	fdescripcion := r.FormValue("txtDescripcion")

	var tramiteEstado TramiteEstadoDTO

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

		var tramiteCambioEstado TramiteCambioEstadoDTO
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
	ar := AuthenticationRequest{fid, fpassword}
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
			Tramites []datoTramitesTable
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
	var tramitesDTO []TramiteRegistradoDTO
	tramitesDTO = append(tramitesDTO, tramiteDTO)
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

func buscarPorCedula(w http.ResponseWriter, r *http.Request) {
	fcedula := r.FormValue("txtCedula")

	tramitesDTO := findTramitesRegistradosByCedulaCliente(fcedula)
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

func buscarPorEstado(w http.ResponseWriter, r *http.Request) {
	festado := r.FormValue("cbxEstado")

	tramitesDTO := findAllTramitesRegistrados()
	tramitesTable := crearDatosTable(tramitesDTO)

	var tramites []datoTramitesTable
	for i := 0; i < len(tramitesTable); i++ {
		if tramitesTable[i].Estado == festado {
			tramites = append(tramites, tramitesTable[i])
		}
	}

	d := struct {
		Usuario  string
		Tramites []datoTramitesTable
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
	var tramites []datoTramitesTable
	for i := 0; i < len(tramitesTable); i++ {
		if getFecha(tramitesTable[i].FechaRegistro) == getFecha(fechaString) {
			tramites = append(tramites, tramitesTable[i])
		}
	}
	d := struct {
		Usuario  string
		Tramites []datoTramitesTable
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

//funcion para crear el tramite Registrado
func getTramiteRegistradoView(tramiteRegistradoDTO TramiteRegistradoDTO) (tramiteRegistradoView TramiteRegistrado) {
	tramiteRegistradoView.NombreCliente = tramiteRegistradoDTO.ClienteID.NombreCompleto
	tramiteRegistradoView.CedulaCliente = tramiteRegistradoDTO.ClienteID.Cedula
	tipoTramite := findTipoTramiteByID(tramiteRegistradoDTO.ID)
	tramiteRegistradoView.DescripcionTipoTramite = tipoTramite.Descripcion
	tramiteRegistradoView.NombreDepartamento = tipoTramite.Departamento.Nombre
	tramiteRegistradoView.EstadoActualNombre = obtenerUltimoEstado(tramiteRegistradoDTO.ID).NombreTramiteEstado
	tramiteRegistradoView.DescripcionEstado = obtenerUltimoEstado(tramiteRegistradoDTO.ID).DescripcionTramiteEstado

	notasDTO := findNotasByTramiteRegistradoID(tramiteRegistradoDTO.ID)
	var notas []Notas
	if len(notasDTO) > 0 {
		for i := 0; i < len(notasDTO); i++ {
			nota := Notas{Titulo: notasDTO[i].Titulo, Contenido: notasDTO[i].Contenido}
			notas = append(notas, nota)
		}
	}
	tramiteRegistradoView.Notas = notas

	requisitosR := findRequisitosPresentadosByTramiteRegistradoID(tramiteRegistradoDTO.ID)
	var requisitosPresentados []RequisitosPresentados
	if len(requisitosR) > 0 {
		fmt.Println("entro if")
		for i := 0; i < len(requisitosR); i++ {
			requisitoPresentado := RequisitosPresentados{FechaRegistro: requisitosR[i].FechaRegistro.Format("Mon Jan _2 15:04:05 2006"), NombreRequisito: requisitosR[i].RequisitoID.Descripcion, DescripcionVariacion: requisitosR[i].RequisitoID.Variacion.Descripcion}
			requisitosPresentados = append(requisitosPresentados, requisitoPresentado)
		}
		fmt.Printf("%v", requisitosPresentados)
	}
	tramiteRegistradoView.Requisitos = requisitosPresentados

	tramitesCambio := findTramiteCambioEstadoByTramiteRegistradoID(tramiteRegistradoDTO.ID)
	var tramitesCambioEstados []TramitesCambioEstados
	if len(tramitesCambio) > 0 {
		for i := 0; i < len(tramitesCambio); i++ {
			tramiteCE := TramitesCambioEstados{NombreTramiteEstado: tramitesCambio[i].TramiteEstadoID.Nombre, DescripcionTramiteEstado: tramitesCambio[i].TramiteEstadoID.Descripcion, NombreUsuario: tramitesCambio[i].UsuarioID.NombreCompleto, FechaRegistro: tramitesCambio[i].FechaRegistro.Format("Mon Jan _2 15:04:05 2006")}
			tramitesCambioEstados = append(tramitesCambioEstados, tramiteCE)
		}
	}
	fmt.Printf("%v", tramitesCambioEstados)
	tramiteRegistradoView.TramitesCambioEstados = tramitesCambioEstados

	return tramiteRegistradoView
}

func findAllTramitesRegistrados() (tramitesregistrados []TramiteRegistradoDTO) {
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

func findTramitesRegistradosByID(idTR int64) (tramiteregistrado TramiteRegistradoDTO) {
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

func findTramitesRegistradosByCedulaCliente(cedula string) (tramiteregistrados []TramiteRegistradoDTO) {
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

func findTipoTramiteByID(idTT int64) (tramitetipo TramiteTipoDTO) {
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

func findNotasByTramiteRegistradoID(idTR int64) (notas []NotaDTO) {
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

func findRequisitosPresentadosByTramiteRegistradoID(idTR int64) (requisitosR []RequisitoPresentadoDTO) {
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

func findTramiteCambioEstadoByTramiteRegistradoID(idTR int64) (tramitesCambio []TramiteCambioEstadoDTO) {
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

func crearDatosTable(tramitesRegistrados []TramiteRegistradoDTO) (tramitesTable []datoTramitesTable) {
	for i := 0; i < len(tramitesRegistrados); i++ {
		tramite := datoTramitesTable{ID: tramitesRegistrados[i].ID, NombreCliente: tramitesRegistrados[i].ClienteID.NombreCompleto, CedulaCliente: tramitesRegistrados[i].ClienteID.Cedula, TipoTramite: findTipoTramiteByID(int64(tramitesRegistrados[i].TramitesTiposID)).Descripcion, FechaRegistro: obtenerUltimoEstado(tramitesRegistrados[i].ID).FechaRegistro, Estado: obtenerUltimoEstado(tramitesRegistrados[i].ID).NombreTramiteEstado}
		tramitesTable = append(tramitesTable, tramite)
	}
	return tramitesTable

}

func obtenerUltimoEstado(idTR int64) (tramiteCE TramitesCambioEstados) {
	var tramiteCEDTO TramiteCambioEstadoDTO
	tramitesCambio := findTramiteCambioEstadoByTramiteRegistradoID(idTR)
	if len(tramitesCambio) > 0 {
		idMayor := tramitesCambio[0].ID
		tramiteCEDTO = tramitesCambio[0]
		for i := 0; i < len(tramitesCambio); i++ {
			if idMayor < tramitesCambio[i].ID {
				tramiteCEDTO = tramitesCambio[i]
			}
		}
		tramiteCE.NombreTramiteEstado = tramiteCEDTO.TramiteEstadoID.Nombre
		tramiteCE.DescripcionTramiteEstado = tramiteCEDTO.TramiteEstadoID.Descripcion
		tramiteCE.NombreUsuario = tramiteCEDTO.UsuarioID.NombreCompleto
		tramiteCE.FechaRegistro = tramiteCEDTO.FechaRegistro.Format("Mon Jan _2 15:04:05 2006")
	}
	return tramiteCE
}

//Estructura para el table view
type datoTramitesTable struct {
	ID            int64
	NombreCliente string
	CedulaCliente string
	TipoTramite   string
	Estado        string
	FechaRegistro string
}

//RequisitosPresentados es una estructura para la lista de requisitosPresentados con la informacion de ellos que queremos mostrar
type RequisitosPresentados struct {
	FechaRegistro        string
	NombreRequisito      string
	DescripcionVariacion string
}

//Notas es una estructura para la lista de notas con la informacion de ellas que queremos mostrar
type Notas struct {
	Titulo    string
	Contenido string
}

//TramitesCambioEstados es una estructura para la lista de los cambios de estado de los tramites con la informacion de ellos que queremos mostrar
type TramitesCambioEstados struct {
	NombreTramiteEstado      string
	DescripcionTramiteEstado string
	NombreUsuario            string
	FechaRegistro            string
}

//TramiteRegistrado es la estructura de los tramites registrados que queremos mostrar en el html
type TramiteRegistrado struct {
	Usuario                string
	NombreCliente          string
	CedulaCliente          string
	DescripcionTipoTramite string
	NombreDepartamento     string
	Requisitos             []RequisitosPresentados
	Notas                  []Notas
	TramitesCambioEstados  []TramitesCambioEstados
	EstadoActualNombre     string
	DescripcionEstado      string
}
