package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Estructura
type Empleados struct {
	Id     int
	Nombre string
	Correo string
}

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+contrasenia+"@tcp(127.0.0.1)/"+Nombre) // Conexion a la db

	// Manejador de errores
	if err != nil {
		panic(err.Error())
	}

	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*")) // Obtener informacion dentro de una carpeta para las plantillas

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/crear", Create)
	http.HandleFunc("/insertar", Insert) // Para insertar los datos
	http.HandleFunc("/borrar", Delete)   // Para borrar los datos

	fmt.Println("Server Started ....") // hace lo mismo que log
	//log.Println("Server Started ....")
	http.ListenAndServe(":8080", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hola!!!") // Imprimir en el navegador

	// INSERTAR DATOS
	// conexionEstablecida := conexionDB()
	// insertarRegistro, err := conexionEstablecida.Prepare("INSERT INTO `empleados` (`id`, `nombre`, `correo`) VALUES (NULL, 'mauriTest', 'mTest@gmail.com')")

	// if err != nil {
	// 	panic(err.Error()) // Mostrar Error
	// }

	// insertarRegistro.Exec()

	// OBTENER REGISTROS
	conexionEstablecida := conexionDB()
	resgistros, err := conexionEstablecida.Query("SELECT * FROM empleados") // ejecuta sentencia sql y la devuelve

	if err != nil {
		panic(err.Error()) // Mostrar Error
	}

	empleado := Empleados{}
	arregloEmpleado := []Empleados{}

	for resgistros.Next() { // recorrer registros
		var id int
		var nombre, correo string
		err = resgistros.Scan(&id, &nombre, &correo) // saca la informacion y la coloca en las variables

		if err != nil {
			panic(err.Error()) // Mostrar Error
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	// fmt.Println(arregloEmpleado)

	// plantillas.ExecuteTemplate(w, "inicio", nil)             // accede a la plantilla inicio
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado) // accede a la plantilla inicio y envia arregloEmpleado
}

func Create(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil) // accede a la plantilla crear
}

func Insert(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" { // si el metodo es post
		// obtener datos del formulario
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		// INSERTAR DATOS
		conexionEstablecida := conexionDB()
		insertarRegistro, err := conexionEstablecida.Prepare("INSERT INTO empleados (nombre, correo) VALUES (?,?)")

		if err != nil {
			panic(err.Error()) // Mostrar Error
		}

		insertarRegistro.Exec(nombre, correo) // envia los datos parametros para q los cambie por ?
	}

	http.Redirect(w, r, "/", 301) // redireccionar
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") // obtener datos desde la url
	// fmt.Println(idEmpleado)               // pare verificar la devolucion

	// INSERTAR DATOS
	conexionEstablecida := conexionDB()
	deleteRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados where id = ?")

	if err != nil {
		panic(err.Error()) // Mostrar Error
	}

	deleteRegistro.Exec(idEmpleado) // envia los datos parametros para q los cambie por ?
	http.Redirect(w, r, "/", 301)   // redireccionar
}
