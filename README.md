# fuegoquasar

## Operación Fuego de Quasar

Han Solo ha sido recientemente nombrado General de la Alianza
Rebelde y busca dar un gran golpe contra el Imperio Galáctico para
reavivar la llama de la resistencia.
El servicio de inteligencia rebelde ha detectado un llamado de auxilio de
una nave portacarga imperial a la deriva en un campo de asteroides. El
manifiesto de la nave es ultra clasificado, pero se rumorea que
transporta raciones y armamento para una legión entera.

### Desafío

Como jefe de comunicaciones rebelde, tu misión es crear un programa en Golang que retorne
la fuente y contenido del mensaje de auxilio . Para esto, cuentas con tres satélites que te
permitirán triangular la posición, ¡pero cuidado! el mensaje puede no llegar completo a cada
satélite debido al campo de asteroides frente a la nave.
Posición de los satélites actualmente en servicio
- Kenobi: [-500, -200]

- Skywalker: [100, -100]

- Sato: [500, 100]

### Nivel 1

Crear un programa con las siguientes firmas:
```golang
// input: distancia al emisor tal cual se recibe en cada satélite
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
func GetLocation(distances ...float32) (x, y float32)
// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
func GetMessage(messages ...[]string) (msg string)
```
#### Consideraciones:

- La unidad de distancia en los parámetros de GetLocation es la misma que la que se
utiliza para indicar la posición de cada satélite.

- El mensaje recibido en cada satélite se recibe en forma de arreglo de strings.

- Cuando una palabra del mensaje no pueda ser determinada, se reemplaza por un string
en blanco en el array.

○ Ejemplo: [“este”, “es”, “”, “mensaje”]

- Considerar que existe un desfasaje (a determinar) en el mensaje que se recibe en cada
satélite.

○ Ejemplo:

■ Kenobi: [“”, “este”, “es”, “un”, “mensaje”]

■ Skywalker: [“este”, “”, “un”, “mensaje”]

■ Sato: [“”, ””, ”es”, ””, ”mensaje”]

### Nivel 2

Crear una API REST, hostear esa API en un cloud computing libre (Google App Engine,
Amazon AWS, etc), crear el servicio /topsecret/ en donde se pueda obtener la ubicación de
la nave y el mensaje que emite.
El servicio recibirá la información de la nave a través de un HTTP POST con un payload con el
siguiente formato:
POST → /topsecret/
```json
{
"satellites": [
{
“name”: "kenobi",
“distance”: 100.0,
“message”: ["este", "", "", "mensaje", ""]
},
{
“name”: "skywalker",
“distance”: 115.5
“message”: ["", "es", "", "", "secreto"]
},
{
“name”: "sato",
“distance”: 142.7
“message”: ["este", "", "un", "", ""]
}
]
}
```
La respuesta, por otro lado, deberá tener la siguiente forma:

RESPONSE CODE: 200
```json
{
"position": {
"x": -100.0,
"y": 75.5
},
"message": "este es un mensaje secreto"
}
```
Nota: la respuesta en este ejemplo es meramente ilustrativa y no debe ser considerada como
caso de prueba para validar la solución propuesta.
En caso que no se pueda determinar la posición o el mensaje, retorna:

RESPONSE CODE: 404

### Nivel 3
Considerar que el mensaje ahora debe poder recibirse en diferentes POST al nuevo servicio
/topsecret_split/ , respetando la misma firma que antes. Por ejemplo:
POST → /topsecret_split/{satellite_name}
```json
{
"distance": 100.0,
"message": ["este", "", "", "mensaje", ""]
}
```
Crear un nuevo servicio /topsecret_split/ que acepte POST y GET. En el GET la
respuesta deberá indicar la posición y el mensaje en caso que sea posible determinarlo y tener
la misma estructura del ejemplo del Nivel 2. Caso contrario, deberá responder un mensaje de
error indicando que no hay suficiente información.

### Entregables
● URL en donde este hosteado el servicio

El servicio esta hosteado en la url https://fluent-burner-308712.rj.r.appspot.com/

● Código fuente en repositorio privado de GitHub

El Codigo esta disponible en el repositorio de GITHUB https://github.com/lucaskruk/fuegoquasar

● Documentación que indique cómo ejecutar el programa

Para poder ejecutar la aplicacion, hay que tener el entorno de GO instalado, para informacion de como obtenerlo, ve a 
https://golang.org/doc/install

En tu entorno local de go, se necesitan los siguientes paquetes:
"github.com/gorilla/mux"
"gopkg.in/yaml.v2"
Cuando se ejecuta go build se descargan automaticamente, pero tambien puedes obtenerlos con los comandos 
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2


Una vez instalado go, abre una consola de comandos en el directorio donde hayas clonado el repositorio, y ejecuta el comando
go run fuegoquasar.go
para poder probar la aplicacion. 

En el caso de que necesites modificar el puerto, puedes modificar el archivo /config/config.yml
En ese archivo tambien se encuentra la ubicación de los satelites que usan para obtener la posicion

#### Documentación:

Para poder encontrar la ubicación, lo primero que realicé fue calcular la intersección entre los tres círculos formados por la posicion de cada satelite mas la distancia a la nave que pide auxilio como radio. La informacíon del comportamiento algebráico la encontre en http://paulbourke.net/geometry/circlesphere/

Por cada interseccion entre círculos puedo obtener uno, dos puntos, o ninguno. Luego de calcular las 3 intersecciones, comparo todos los puntos obtenidos para ver cual es el que se repite en todas ellas. Ese punto es la ubicación del objeto en el plano.
A tener en cuenta que los resultados de la posicion estan fijados con una precision maxima de 1 cifra decimal. Aprovechando esta precision definida, para no agregar un campo bool que modifique la firma del metodo GetLocation a la que se encuentra en la consigna, defini el valor -0.09 en las coordenadas de X y de Y del resultado como "Posicion no encontrada" para tomar luego una desición al ser llamado en la api.

Para asegurarme de que el calculo sea correcto, valido que los nombres recibidos en el request coincidan con los nombres que se configuraron, de lo contrario no podría relacionar la ubicación de cada satelite con la distancia recibida en el request.

El método que calcula los mensajes, asume que los tres satélites envian el mismo mensaje, con la posibilidad de que uno o varios elementos del array puedan estar vacios. En el caso de que para una de las posiciones del array de mensajes, no haya valor en ninguno de los tres satelites, se considera como mensaje no generado, y se devuelve el response code 404. Tambien determina el desfasaje contando la cantidad de elementos de cada array y seleccionando el mas corto como el mensaje sin "desfasajes" y trunca los primeros elementos de los otros mensajes. 

El método Get, que devuelve la informacion de topsecret_split, siempre limpia el array de las naves cargadas, para poder volver a realizar la carga en el caso de que hayan sido cargados incorrectamente los datos.

Es crucial pasar los nombres correctos en la URL del metodo POST /topsecret_split/(nombre) asi como tambien en el json del post del metodo /topsecret/ 

Los metodos de la URL /topsecret_split tanto post como get tienen algunos mensajes orientativos en los response, segun los diferentes escenarios.

El formato de response es el mismo en todos los metodos para simplificar un poco las cosas.

##### Detalle de los metodos expuestos y sus parámetros:
Para todos los response, en el caso de que al querer hacer el parseo del json del request ocurra un error, por formato de json no valido, se envia un mensaje mencionando el error.


POST → /topsecret/ 

REQUEST 
 
Satelites: Es un Array del tipo Satelite
Satelite: Es un struct con los siguientes campos: 
	name: String - Lleva el nombre del satelite, es obligatorio para saber su posicion configurada
	distance: float32 - es la distancia del satelite hacia la nave que pide auxilio.
	message: []string - es un array de strings que forman el mensaje de auxilio.

```json 
{
"satelites":[
{
"name": "Kenobi",
"distance": 761.577,
"message": ["","", "este", "es", "o", ""]
},
{
"name": "SkyWalker",
"distance": 223.607,
"message": ["", "es", "", "mensaje"]
},
{
"name": "Sato",
"distance": 300,
"message": ["", "", "es", "o", ""]
}
]
}
```

RESPONSE

Position: Es un struct con las coordenadas de la nave que pide auxilio
	x: float32 coordenada del eje X en el plano
	y: float32 coordenada del eje Y en el plano
Message: es un string con el mensaje de auxilio. En el caso de que ocurra algun problema, los mensajes de error o descripciones se informan en este campo. 

```json
{
    "position": {
        "x": 200,
        "y": 100
    },
    "message": "este es o mensaje"
}
```
En el caso de no encontrar el satelite o no poder armar el mensaje, devuelve status code 404 sin datos en el response

POST → /topsecret_split/{satellite_name}
REQUEST:

	distance: float32 - es la distancia del satelite hacia la nave que pide auxilio.
	message: []string - es un array de strings que forman el mensaje de auxilio.
```json
{
"distance": 300,
"message": ["","este", "es", "un", "", "especial"]
}
```
RESPONSE: identico al método anterior, con la salvedad de que en caso de cargar con éxito, devuelve en el campo message "Satelite <name> cargado". En el caso de haber intentado cargar una cantidad de satelites superior a la cantidad configurada, dara un mensaje de error.
Caso OK
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "Satelite sato cargado."
} ```
en caso de superar el maximo
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "Todas las naves fueron cargadas. Ejecute un get para obtener la posicion y limpiar el array."
}
```

GET → /topsecret_split/
Este metodo trabaja con los datos que se fueron cargando en las sucesivas invocaciones del metodo POST sobre esta ruta. 

RESPONSE
casos OK
```json
{
    "position": {
        "x": 200,
        "y": 100
    },
    "message": "este es o mensaje"
}
```
En caso de no encontrar la ubicacion, o no colocar los nombres configurados, devuelve status 404 y el siguiente mensaje:
```json
{
    "position": {
        "x": -0.09,
        "y": -0.09
    },
    "message": "No se pudo encontrar la ubicacion, o no se pudo descifrar el mensaje. Revise sus parametros e intente nuevamente"
}
```
En el caso que no se haya cargado previamente en el metodo post, los satelites configurados, devuelve:
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "No se cargaron los tres satelites, o los nombres no coinciden. Vuelva a cargarlos correctamente"
}```

##### Configuraciones

La aplicacion guarda la informacion de los parametros configurables en el archivo /config/config.yml con el formato YAML.

```yaml
# config.yml

rebelships: #array con los satelites configurados
  -  name: "kenobi" #nombre del satelite, la verificacion luego es sin distinguir mayusculas
     x: -500 #posicion en x del satelite
     y: -200 #posicion en y del satelite

  -  name: "skywalker"
     x: 100
     y: -100

  -  name: "sato"
     x: 500
     y: 100

#connection settings  
server: #datos del webserver
  host: #vacio implica localhost
  port: 8080 #puerto para ejecucion local
  timeout: #limites de tiempo de espera para dar error de timeout
    server: 30
    read: 15
    write: 10
    idle: 5
```
Editando estos parametros y reiniciando la aplicacion se pueden realizar diferentes pruebas. 
El array de satelites permite mas de 3 naves, en el caso de ser menos de 3 no se podria calcular la posicion. En el caso de usar mas de 3 solo se consideran los tres primeros dado que hacer el calculo de todas las permutaciones se volvia pesado y excede los requisitos de este desafio.
Si se consideran los mensajes de todas las naves configuradas.

##### Tests automatizados

Se generaron algunas pruebas simples de los metodos en el archivo /api/api_test.go
para poder ejecutar las pruebas, basta con abrir la consola en el directorio donde se encuentra el archivo de test, y ejecutar el comando: GO TEST -v
con el parametro v para obtener el verbose de logs