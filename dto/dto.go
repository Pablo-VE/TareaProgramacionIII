package dto

import "time"

//AuthenticationRequest is...
type AuthenticationRequest struct {
	Cedula   string `json:"cedula"`
	Password string `json:"password"`
}

//AuthenticationResponse is...
type AuthenticationResponse struct {
	Jwt      string               `json:"jwt"`
	Usuario  UsuarioDTO           `json:"usuario"`
	Permisos []PermisoOtorgadoDTO `json:"permisos"`
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
	TramiteRegistradoID TramiteRegistradoDTO `json:"tramitesRegistrados"`
	RequisitoID         RequisitoDTO         `json:"requisito"`
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
