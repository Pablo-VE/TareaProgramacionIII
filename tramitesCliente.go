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
	http.HandleFunc("/login", login)
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
	var respuesta string
	if res.StatusCode != 200 {
		respuesta = "Credenciales erroneas"
	} else {
		respuesta = "Login Exitoso"
	}
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &usuarioLogeado)
	fmt.Println("Authetication Response Struct: ", usuarioLogeado)
	d := struct {
		Respuesta string
	}{
		Respuesta: respuesta,
	}
	tpl.ExecuteTemplate(w, "menu.html", d)
}

func findAllTramitesRegistrados() {
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
		var tramitesregistrados []TramiteRegistradoDTO
		json.Unmarshal(bodyBytes, &tramitesregistrados)
	}
}

func findTramitesRegistradosByID(idTR int64) {
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
		var tramiteregistrado TramiteRegistradoDTO
		json.Unmarshal(bodyBytes, &tramiteregistrado)
	}
}

func findTipoTramiteByID(idTT int64) {
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
		var tramitetipo TramiteTipoDTO
		json.Unmarshal(bodyBytes, &tramitetipo)
	}

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

/*
func crearTramiteVista(int64 idTramiteRegistrado) {
	//buscarTramiteRegistradoPorId

}*/

//Estructura para el table view
type datoTramitesTable struct {
	id            int64
	nombreCliente string
	cedulaCliente string
	estado        string
	fechaRegistro string
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
	fechaRegistro            time.Time
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
	Variaciones   VariacionDTO `json:"variaciones"`
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
	UsuarioID             UsuarioDTO           `json:"usuarioId"`
	TramitesRegistradosID TramiteRegistradoDTO `json:"tramitesRegistradosId"`
	TramiteEstadoID       TramiteEstadoDTO     `json:"tramitesEstadoId"`
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
