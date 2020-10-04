package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/TramitesRegistrados", limpiar)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
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
	req, err := http.NewRequest("GET", url+"tramites_registrados/"+cedula, nil)
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

func findNotasByTramiteRegistradoID(idTR int64) {
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
	var notas []NotaDTO
	if res.StatusCode == 200 {

		json.Unmarshal(bodyBytes, &notas)
	}
}

func findRequisitosPresentadosByTramiteRegistradoID(idTR int64) {
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
	var requisitosR []RequisitoPresentadoDTO
	if res.StatusCode == 200 {

		json.Unmarshal(bodyBytes, &requisitosR)
	}
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
		tramite := datoTramitesTable{ID: tramitesRegistrados[i].ID, NombreCliente: tramitesRegistrados[i].ClienteID.NombreCompleto, CedulaCliente: tramitesRegistrados[i].ClienteID.Cedula, TipoTramite: findTipoTramiteByID(int64(tramitesRegistrados[i].TramitesTiposID)).Descripcion, FechaRegistro: obtenerUltimoEstado(tramitesRegistrados[i].ID).fechaRegistro, Estado: obtenerUltimoEstado(tramitesRegistrados[i].ID).nombreTramiteEstado}
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
		tramiteCE.nombreTramiteEstado = tramiteCEDTO.TramiteEstadoID.Nombre
		tramiteCE.descripcionTramiteEstado = tramiteCEDTO.TramiteEstadoID.Descripcion
		tramiteCE.nombreUsuario = tramiteCEDTO.UsuarioID.NombreCompleto
		tramiteCE.fechaRegistro = tramiteCEDTO.FechaRegistro.String()
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

//requisitosPresentados es una estructura para la lista de requisitosPresentados con la informacion de ellos que queremos mostrar
type requisitosPresentados struct {
	fechaRegistro        time.Time
	nombreRequisito      string
	descripcionVariacion string
}

//Notas es una estructura para la lista de notas con la informacion de ellas que queremos mostrar
type Notas struct {
	titulo    string
	contenido string
}

//TramitesCambioEstados es una estructura para la lista de los cambios de estado de los tramites con la informacion de ellos que queremos mostrar
type TramitesCambioEstados struct {
	nombreTramiteEstado      string
	descripcionTramiteEstado string
	nombreUsuario            string
	fechaRegistro            string
}

//TramiteRegistrado es la estructura de los tramites registrados que queremos mostrar en el html
type TramiteRegistrado struct {
	nombreCliente          string
	cedulaCliente          string
	descripcionTipoTramite string
	nombreDepartamento     string
	requisitos             []requisitosPresentados
	notas                  []Notas
	tramitesCambioEstados  []TramitesCambioEstados
	estadoActualNombre     string
	descripcionEstado      string
}

//AuthenticationRequest is...
type AuthenticationRequest struct {
	Cedula   string `json:"cedula"`
	Password string `json:"password"`
}

//AuthenticationResponse is...
type AuthenticationResponse struct {
	Jwt      string               `json:"jwt"`
	Usuario  UsuarioDTO           `json:"usuario"`
	Permisos []PermisoOtorgadoDTO `json:"permisosOtorgados"`
}

//UsuarioDTO is...
type UsuarioDTO struct {
	ID                 int64           `json:"id"`
	NombreCompleto     string          `json:"nombreCompleto"`
	Cedula             string          `json:"cedula"`
	PasswordEncriptado string          `json:"passwordEncriptado"`
	Estado             bool            `json:"estado"`
	FechaRegistro      time.Time       `json:"fechaRegistro"`
	FechaModificacion  time.Time       `json:"fechaModificacion"`
	EsJefe             bool            `json:"esJefe"`
	Departamento       DepartamentoDTO `json:"departamento"`
}

//DepartamentoDTO is...
type DepartamentoDTO struct {
	ID                int64     `json:"id"`
	Nombre            string    `json:"nombre"`
	Estado            bool      `json:"estado"`
	FechaRegistro     time.Time `json:"fechaRegistro"`
	FechaModificacion time.Time `json:"fechaModificacion"`
}

//TramiteTipoDTO is...
type TramiteTipoDTO struct {
	ID                int64           `json:"id"`
	Descripcion       string          `json:"descripcion"`
	Estado            bool            `json:"estado"`
	Departamento      DepartamentoDTO `json:"departamento"`
	FechaRegistro     time.Time       `json:"fechaRegistro"`
	FechaModificacion time.Time       `json:"fechaModificacion"`
	Variaciones       []VariacionDTO  `json:"variaciones"`
}

//VariacionDTO is...
type VariacionDTO struct {
	ID            int64          `json:"id"`
	Grupo         int32          `json:"grupo"`
	Descripcion   string         `json:"descripcion"`
	Estado        bool           `json:"estado"`
	FechaRegistro time.Time      `json:"fechaRegistro"`
	TramitesTipos TramiteTipoDTO `json:"tramitesTipos"`
}

//RequisitoDTO is...
type RequisitoDTO struct {
	ID            int64        `json:"id"`
	Descripcion   string       `json:"descripcion"`
	Estado        bool         `json:"estado"`
	FechaRegistro time.Time    `json:"fechaRegistro"`
	Variacion     VariacionDTO `json:"variaciones"`
}

//ClienteDTO is...
type ClienteDTO struct {
	ID                 int64     `json:"id"`
	NombreCompleto     string    `json:"nombreCompleto"`
	Cedula             string    `json:"cedula"`
	Telefono           string    `json:"telefono"`
	Direccion          string    `json:"direccion"`
	Estado             bool      `json:"estado"`
	FechaRegistro      time.Time `json:"fechaRegistro"`
	FechaModificacion  time.Time `json:"fechaModificacion"`
	PasswordEncriptado string    `json:"passwordEncriptado"`
}

//TramiteRegistradoDTO is...
type TramiteRegistradoDTO struct {
	ID              int64      `json:"id"`
	TramitesTiposID int32      `json:"tramitesTiposId"`
	ClienteID       ClienteDTO `json:"cliente"`
}

//RequisitoPresentadoDTO is...
type RequisitoPresentadoDTO struct {
	ID                  int64                `json:"id"`
	FechaRegistro       time.Time            `json:"fechaRegistro"`
	TramiteRegistradoID TramiteRegistradoDTO `json:"tramiteRegistradoId"`
	RequisitoID         RequisitoDTO         `json:"requisitoId"`
}

//TramiteCambioEstadoDTO is...
type TramiteCambioEstadoDTO struct {
	ID                    int64                `json:"id"`
	UsuarioID             UsuarioDTO           `json:"usuario"`
	TramitesRegistradosID TramiteRegistradoDTO `json:"tramitesRegistrados"`
	TramiteEstadoID       TramiteEstadoDTO     `json:"tramitesEstados"`
	FechaRegistro         time.Time            `json:"fechaRegistro"`
}

//TramiteEstadoDTO is...
type TramiteEstadoDTO struct {
	ID               int64  `json:"id"`
	Nombre           string `json:"nombre"`
	Descripcion      string `json:"descripcion"`
	EstadosSucesores string `json:"estadosSucesores"`
}

//PermisoDTO is...
type PermisoDTO struct {
	ID                int64     `json:"id"`
	Codigo            string    `json:"codigo"`
	Descripcion       string    `json:"descripcion"`
	FechaRegistro     time.Time `json:"fechaRegistro"`
	FechaModificacion time.Time `json:"fechaModificacion"`
	Estado            bool      `json:"estado"`
}

//PermisoOtorgadoDTO is...
type PermisoOtorgadoDTO struct {
	ID            int64      `json:"id"`
	Usuario       UsuarioDTO `json:"usuario"`
	Permiso       PermisoDTO `json:"permiso"`
	FechaRegistro time.Time  `json:"fechaRegistro"`
	Estado        bool       `json:"estado"`
}

//NotaDTO is...
type NotaDTO struct {
	ID                  int64                `json:"id"`
	Estado              bool                 `json:"estado"`
	Tipo                bool                 `json:"tipo"`
	Titulo              string               `json:"titulo"`
	Contenido           string               `json:"contenido"`
	FechaRegistro       time.Time            `json:"fechaRegistro"`
	FechaModificacion   time.Time            `json:"fechaModificacion"`
	TramitesRegistrados TramiteRegistradoDTO `json:"tramitesRegistrados"`
}
