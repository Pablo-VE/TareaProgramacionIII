package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Pablo-VE/TareaProgramacionIII/conexionservidor"
	"github.com/Pablo-VE/TareaProgramacionIII/dto"
	"github.com/Pablo-VE/TareaProgramacionIII/util"
)

var usuarioLogeado dto.AuthenticationResponse

var solicitudDenegada string

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/tramitesRegistrados", login)
	http.HandleFunc("/alerta", alerta)
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

	tramiteEstadoGuardado := conexionservidor.CreateTramiteEstado(tramiteEstado)
	if tramiteEstadoGuardado.Nombre == "" {
		//irTramite(w, r)
	} else {
		var tramiteCambioEstado dto.TramiteCambioEstadoDTO
		tramiteCambioEstado.TramiteEstadoID = tramiteEstadoGuardado
		tramiteCambioEstado.TramitesRegistradosID = tramiteRegistradoEnCuestion
		tramiteCambioEstado.UsuarioID = usuarioLogeado.Usuario
		tramiteCambioEstado = conexionservidor.CreateTramiteCambioEstado(tramiteCambioEstado)
		limpiar(w, r)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fid := r.FormValue("cedula")
	fpassword := r.FormValue("password")

	usuarioLogeado = conexionservidor.Logear(fid, fpassword)

	if usuarioLogeado.Jwt != "" {

		if verificarPermiso(usuarioLogeado.Permisos, "TRA06") == true {
			tramitesDTO := conexionservidor.FindAllTramitesRegistrados()
			tramitesTable := util.CrearDatosTable(tramitesDTO)
			d := struct {
				Usuario  string
				Tramites []util.DatoTramitesTable
			}{
				Usuario:  usuarioLogeado.Usuario.NombreCompleto,
				Tramites: tramitesTable,
			}
			tpl.ExecuteTemplate(w, "tramites.html", d)
		} else {
			solicitudDenegada = "Observar todos los tramites registrados"
			d := struct {
				Usuario   string
				Solicitud string
			}{
				Usuario:   usuarioLogeado.Usuario.NombreCompleto,
				Solicitud: solicitudDenegada,
			}
			tpl.ExecuteTemplate(w, "alerts.html", d)
		}
	} else {
		tpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func buscarPorID(w http.ResponseWriter, r *http.Request) {

	if verificarPermiso(usuarioLogeado.Permisos, "TRA05") == true {
		fid := r.FormValue("txtID")

		nid, err := strconv.ParseInt(fid, 10, 64)
		if err != nil {

		}
		tramiteDTO := conexionservidor.FindTramitesRegistradosByID(nid)
		var tramitesDTO []dto.TramiteRegistradoDTO
		tramitesDTO = append(tramitesDTO, tramiteDTO)
		tramitesTable := util.CrearDatosTable(tramitesDTO)

		d := struct {
			Usuario  string
			Tramites []util.DatoTramitesTable
		}{
			Usuario:  usuarioLogeado.Usuario.NombreCompleto,
			Tramites: tramitesTable,
		}

		tpl.ExecuteTemplate(w, "tramites.html", d)
	} else {
		solicitudDenegada = "Buscar tramite por id"
		d := struct {
			Usuario   string
			Solicitud string
		}{
			Usuario:   usuarioLogeado.Usuario.NombreCompleto,
			Solicitud: solicitudDenegada,
		}
		tpl.ExecuteTemplate(w, "alerts.html", d)
	}

}

func buscarPorCedula(w http.ResponseWriter, r *http.Request) {
	if verificarPermiso(usuarioLogeado.Permisos, "TRA05") == true {
		fcedula := r.FormValue("txtCedula")

		tramitesDTO := conexionservidor.FindTramitesRegistradosByCedulaCliente(fcedula)
		tramitesTable := util.CrearDatosTable(tramitesDTO)

		d := struct {
			Usuario  string
			Tramites []util.DatoTramitesTable
		}{
			Usuario:  usuarioLogeado.Usuario.NombreCompleto,
			Tramites: tramitesTable,
		}

		tpl.ExecuteTemplate(w, "tramites.html", d)
	} else {
		solicitudDenegada = "Buscar tramites por cedula"
		d := struct {
			Usuario   string
			Solicitud string
		}{
			Usuario:   usuarioLogeado.Usuario.NombreCompleto,
			Solicitud: solicitudDenegada,
		}
		tpl.ExecuteTemplate(w, "alerts.html", d)
	}

}

func buscarPorEstado(w http.ResponseWriter, r *http.Request) {
	if verificarPermiso(usuarioLogeado.Permisos, "TRA05") == true {
		festado := r.FormValue("cbxEstado")

		tramitesDTO := conexionservidor.FindAllTramitesRegistrados()
		tramitesTable := util.CrearDatosTable(tramitesDTO)

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
	} else {
		solicitudDenegada = "Buscar tramites por estado"
		d := struct {
			Usuario   string
			Solicitud string
		}{
			Usuario:   usuarioLogeado.Usuario.NombreCompleto,
			Solicitud: solicitudDenegada,
		}
		tpl.ExecuteTemplate(w, "alerts.html", d)
	}

}

func buscarPorFecha(w http.ResponseWriter, r *http.Request) {
	if verificarPermiso(usuarioLogeado.Permisos, "TRA05") == true {
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
		tramitesDTO := conexionservidor.FindAllTramitesRegistrados()
		tramitesTable := util.CrearDatosTable(tramitesDTO)
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
	} else {
		solicitudDenegada = "Buscar tramites por fecha"
		d := struct {
			Usuario   string
			Solicitud string
		}{
			Usuario:   usuarioLogeado.Usuario.NombreCompleto,
			Solicitud: solicitudDenegada,
		}
		tpl.ExecuteTemplate(w, "alerts.html", d)
	}
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

	tramitesDTO := conexionservidor.FindAllTramitesRegistrados()
	tramitesTable := util.CrearDatosTable(tramitesDTO)

	d := struct {
		Usuario  string
		Tramites []util.DatoTramitesTable
	}{
		Usuario:  usuarioLogeado.Usuario.NombreCompleto,
		Tramites: tramitesTable,
	}

	tpl.ExecuteTemplate(w, "tramites.html", d)

}

func irTramite(w http.ResponseWriter, r *http.Request) {

	if verificarPermiso(usuarioLogeado.Permisos, "TRA02") == true {
		fid := r.FormValue("txtVerID")
		tid, err := strconv.ParseInt(fid, 10, 64)
		if err != nil {

		}
		fmt.Println(tid)

		tramiteRegistradoEnCuestion = conexionservidor.FindTramitesRegistradosByID(tid)

		tramiteRegistradoView := util.GetTramiteRegistradoView(tramiteRegistradoEnCuestion)

		tramiteRegistradoView.Usuario = usuarioLogeado.Usuario.NombreCompleto

		tpl.ExecuteTemplate(w, "TramiteRegistrado.html", tramiteRegistradoView)
	} else {
		solicitudDenegada = "Modificar estado de tramite"
		d := struct {
			Usuario   string
			Solicitud string
		}{
			Usuario:   usuarioLogeado.Usuario.NombreCompleto,
			Solicitud: solicitudDenegada,
		}
		tpl.ExecuteTemplate(w, "alerts.html", d)
	}

}

func verificarPermiso(permisos []dto.PermisoOtorgadoDTO, codigo string) (credencial bool) {

	if len(permisos) > 0 {
		for i := 0; i < len(permisos); i++ {
			if permisos[i].Permiso.Codigo == codigo {
				return true
			}
		}
	}
	return false
}

func alerta(w http.ResponseWriter, r *http.Request) {
	if solicitudDenegada == "Observar todos los tramites registrados" {
		tpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		if solicitudDenegada == "Modificar estado de tramite" || solicitudDenegada == "Buscar tramite por id" || solicitudDenegada == "Buscar tramites por cedula" || solicitudDenegada == "Buscar tramites por estado" || solicitudDenegada == "Buscar tramites por fecha" {
			limpiar(w, r)
		}
	}

}
