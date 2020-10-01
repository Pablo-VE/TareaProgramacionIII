package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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
	req, err := http.NewRequest("POST", "http://localhost:8989/usuarios/login", bytes.NewBuffer(j))
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
		respuesta = "Credenciales erroneas"
	} else {
		respuesta = "Login Exitoso"
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)
	fmt.Println("Authetication Response String: ", bodyString)
	var data AuthenticationResponse
	json.Unmarshal(body, &data)

	d := struct {
		Respuesta string
	}{
		Respuesta: respuesta,
	}

	tpl.ExecuteTemplate(w, "menu.html", d)
}

var respuesta string

//AuthenticationRequest es el dto para hacer el login
type AuthenticationRequest struct {
	Cedula   string `json:"cedula"`
	Password string `json:"password"`
}

//AuthenticationResponse es el dto para la respuesta del login
type AuthenticationResponse struct {
	Jwt      string               `json:"jwt"`
	Usuario  UsuarioDTO           `json:"usuario"`
	Permisos []permisoOtorgadoDTO `json:"permisos"`
}

func request() {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Funciono xD")
		fmt.Println(string(data))
	}
}

//UsuarioDTO es
type UsuarioDTO struct {
	ID                 int64           `json:"id"`
	NombreCompleto     string          `json:"nombreCompleto"`
	Cedula             string          `json:"cedula"`
	PasswordEncriptado string          `json:"passwordEncriptado"`
	Estado             bool            `json:"estado"`
	FechaRegistro      time.Time       `json:"fechaRegistro"`
	FechaModificacion  time.Time       `json:"fechaModificacion"`
	EsJefe             bool            `json:"esJefe"`
	Departamento       departamentoDTO `json:"departamento"`
}

type departamentoDTO struct {
	ID                int64     `json:"id"`
	Nombre            string    `json:"nombre"`
	Estado            bool      `json:"estado"`
	FechaRegistro     time.Time `json:"fechaRegistro"`
	FechaModificacion time.Time `json:"fechaModificacion"`
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
	grupo         int32
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
	tramitesTiposID int32
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
	usuarioID             UsuarioDTO
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

type permisoDTO struct {
	id                int64
	codigo            string
	descripcion       string
	fechaRegistro     time.Time
	fechaModificacion time.Time
	estado            bool
}

type permisoOtorgadoDTO struct {
	id            int64
	usuario       UsuarioDTO
	permiso       permisoDTO
	fechaRegistro time.Time
	estado        bool
}
