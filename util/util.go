package util

import (
	"fmt"

	"github.com/Pablo-VE/TareaProgramacionIII/dto"
)

//DatoTramitesTable es una estructura para el table view
type DatoTramitesTable struct {
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

//GetTramiteRegistradoView es una funcion para crear el tramite Registrado
func GetTramiteRegistradoView(tramiteRegistradoDTO dto.TramiteRegistradoDTO) (tramiteRegistradoView TramiteRegistrado) {
	tramiteRegistradoView.NombreCliente = tramiteRegistradoDTO.ClienteID.NombreCompleto
	tramiteRegistradoView.CedulaCliente = tramiteRegistradoDTO.ClienteID.Cedula
	tipoTramite := findTipoTramiteByID(tramiteRegistradoDTO.ID)
	tramiteRegistradoView.DescripcionTipoTramite = tipoTramite.Descripcion
	tramiteRegistradoView.NombreDepartamento = tipoTramite.Departamento.Nombre
	tramiteRegistradoView.EstadoActualNombre = ObtenerUltimoEstado(tramiteRegistradoDTO.ID).NombreTramiteEstado
	tramiteRegistradoView.DescripcionEstado = ObtenerUltimoEstado(tramiteRegistradoDTO.ID).DescripcionTramiteEstado

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

//CrearDatosTable is...
func CrearDatosTable(tramitesRegistrados []dto.TramiteRegistradoDTO) (tramitesTable []datoTramitesTable) {
	for i := 0; i < len(tramitesRegistrados); i++ {
		tramite := datoTramitesTable{ID: tramitesRegistrados[i].ID, NombreCliente: tramitesRegistrados[i].ClienteID.NombreCompleto, CedulaCliente: tramitesRegistrados[i].ClienteID.Cedula, TipoTramite: findTipoTramiteByID(int64(tramitesRegistrados[i].TramitesTiposID)).Descripcion, FechaRegistro: obtenerUltimoEstado(tramitesRegistrados[i].ID).FechaRegistro, Estado: obtenerUltimoEstado(tramitesRegistrados[i].ID).NombreTramiteEstado}
		tramitesTable = append(tramitesTable, tramite)
	}
	return tramitesTable

}

//ObtenerUltimoEstado is ..
func ObtenerUltimoEstado(idTR int64) (tramiteCE TramitesCambioEstados) {
	var tramiteCEDTO dto.TramiteCambioEstadoDTO
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
