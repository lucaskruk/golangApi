# locator

## Api Localizadora

Locator es un programa en Golang que retorna
la fuente y contenido de un mensaje de auxilio de un barco perdido. Lo hace
triangulando la posición, de tres barcos contra el que esta pidiendo ayuda.Tambien completa el mensaje recibido entre los tres barcos de rescate.
Posición de los barcos de rescate
- BarcoUno: [-500, -200]

- BarcoDos: [100, -100]

- BarcoTres: [500, 100]


El programa tiene las siguientes firmas:
```golang
// input: distancia al emisor tal cual se recibe en cada barco
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
func GetLocation(distances ...float32) (x, y float32)
// input: el mensaje tal cual es recibido en cada barco
// output: el mensaje tal cual lo genera el emisor del mensaje
func GetMessage(messages ...[]string) (msg string)
```
#### Consideraciones:

- La unidad de distancia en los parámetros de GetLocation es la misma que la que se
utiliza para indicar la posición de cada barco rescatista.

- El mensaje recibido en cada barco se recibe en forma de arreglo de strings.

- Cuando una palabra del mensaje no pueda ser determinada, se reemplaza por un string
en blanco en el array.

○ Ejemplo: [“este”, “es”, “”, “mensaje”]

- Considerar que existe un desfasaje (a determinar) en el mensaje que se recibe en cada
satélite.

○ Ejemplo:

■ BarcoUno: [“”, “este”, “es”, “un”, “mensaje”]

■ BarcoDos: [“este”, “”, “un”, “mensaje”]

■ BarcoTres: [“”, ””, ”es”, ””, ”mensaje”]

Se expone el endpoint /helpme/ en donde se puede obtener la ubicación de
la nave y el mensaje que emite. En el caso de no poder determinar la posicion o el mensaje, devuelve 404

Tambien se expone el endpoint /helpme_split/ , respetando la misma firma que antes. Por ejemplo:
Este endpoint acepta POST y GET. En el GET la
respuesta indica la posición y el mensaje en caso que sea posible determinarlo y tiene
la misma estructura del response del endpoint /helpme/ 
De lo contrario, responde un mensaje de
error indicando que no hay suficiente información.


Para poder ejecutar la aplicacion, hay que tener el entorno de GO instalado, para informacion de como obtenerlo, ve a 
https://golang.org/doc/install

En tu entorno local de go, se necesitan los siguientes paquetes:
"github.com/gorilla/mux"
"gopkg.in/yaml.v2"
Cuando se ejecuta go build se descargan automaticamente, pero tambien puedes obtenerlos con los comandos 
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2


Una vez instalado go, abre una consola de comandos en el directorio donde hayas clonado el repositorio, y ejecuta el comando
go run locator.go
para poder probar la aplicacion. 

En el caso de que necesites modificar el puerto, puedes modificar el archivo /config/config.yml
En ese archivo tambien se encuentra la ubicación de los satelites que usan para obtener la posicion.

Para invocar los diferentes metodos post y get, una posibilidad es utilizar POSTMAN https://www.postman.com/downloads/ donde se puede colocar la url, el metodo de la invocacion, y el body con los datos de prueba en formato json. 

#### Como lo hice:

Para poder encontrar la ubicación, lo primero que realicé fue calcular la intersección entre los tres círculos formados por la posicion de cada satelite mas la distancia a la nave que pide auxilio como radio. La informacíon del comportamiento algebráico la encontre en http://paulbourke.net/geometry/circlesphere/

Por cada interseccion entre círculos puedo obtener uno, dos puntos, o ninguno. Luego de calcular las 3 intersecciones, comparo todos los puntos obtenidos para ver cual es el que se repite en todas ellas. Ese punto es la ubicación del objeto en el plano.
A tener en cuenta que los resultados de la posicion estan fijados con una precision maxima de 1 cifra decimal. Aprovechando esta precision definida, para no agregar un campo bool que modifique la firma del metodo GetLocation a la que se encuentra en la consigna, defini el valor -0.09 en las coordenadas de X y de Y del resultado como "Posicion no encontrada" para tomar luego una decisión al ser llamado en la api.

Para asegurarme de que el calculo sea correcto, valido que los nombres recibidos en el request coincidan con los nombres que se configuraron, de lo contrario no podría relacionar la ubicación de cada satelite con la distancia recibida en el request.

El método que calcula los mensajes, asume que los tres barcos reciben el mismo mensaje, con la posibilidad de que uno o varios elementos del array puedan estar vacios. En el caso de que para una de las posiciones del array de mensajes, no haya valor en ninguno de los tres satelites, se considera como mensaje no generado, y se devuelve el response code 404. Tambien determina el desfasaje contando la cantidad de elementos de cada array y seleccionando el mas corto como el mensaje sin "desfasajes" y trunca los primeros elementos de los otros mensajes. A tener en cuenta, que en cada una las posiciones combinadas de los arrays recortados, debe existir un array con valor distinto de vacio.  

El método Get, que devuelve la informacion de helpme_split, siempre limpia el array de las naves cargadas, para poder volver a realizar la carga en el caso de que hayan sido cargados incorrectamente los datos.

Es crucial pasar los nombres correctos en la URL del metodo POST /helpme_split/(nombre) asi como tambien en el json del post del metodo /topsecret/ 

Los metodos de la URL /helpme_split tanto post como get tienen algunos mensajes orientativos en los response, segun los diferentes escenarios.

El formato de response es el mismo en todos los metodos para simplificar un poco las cosas.

##### Detalle de los metodos expuestos y sus parámetros:
Para todos los response, en el caso de que al querer hacer el parseo del json del request ocurra un error, por formato de json no valido, se envia un mensaje mencionando el error.


POST → /helpme/ 

REQUEST 
 
Ships: Es un Array del tipo Satelite
Ship: Es un struct con los siguientes campos: 
	name: String - Lleva el nombre del satelite, es obligatorio para saber su posicion configurada
	distance: float32 - es la distancia del satelite hacia la nave que pide auxilio.
	message: []string - es un array de strings que forman el mensaje de auxilio.

```json 
{
"ships":[
{
"name": "BarcoUno",
"distance": 761.577,
"message": ["","", "este", "es", "o", ""]
},
{
"name": "BarcoDos",
"distance": 223.607,
"message": ["", "es", "", "mensaje"]
},
{
"name": "BarcoTres",
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

POST → /helpme_split/{ship_name}
REQUEST:

	distance: float32 - es la distancia del satelite hacia la nave que pide auxilio.
	message: []string - es un array de strings que forman el mensaje de auxilio.
```json
{
"distance": 300,
"message": ["","este", "es", "un", "", "especial"]
}
```
RESPONSE: identico al método anterior, con la salvedad de que en caso de cargar con éxito, devuelve en el campo message "Barco <name> cargado". En el caso de haber intentado cargar una cantidad de barcos superior a la cantidad configurada, dara un mensaje de error.
Caso OK
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "Ship BarcoUno loaded."
} 
```
en caso de superar el maximo
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "All ships loaded. Run get methot to obtain position and clean array."
}
```

GET → /helpme_split/
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
    "message": "Message can't be recovered. Location can't be identified. Check your parameters and try again."
}
```
En el caso que no se haya cargado previamente en el metodo post, los satelites configurados, devuelve:
```json
{
    "position": {
        "x": 0,
        "y": 0
    },
    "message": "No se cargo la informacion de los tres barcos, o los nombres no coinciden. Vuelva a cargarlos correctamente"
}
```

##### Configuraciones

La aplicacion guarda la informacion de los parametros configurables en el archivo /config/config.yml con el formato YAML.

```yaml
# config.yml

ships: #array con los barcos configurados
  -  name: "barcouno" #nombre del barco, la verificacion luego es sin distinguir mayusculas
     x: -500 #posicion en x del satelite
     y: -200 #posicion en y del satelite

  -  name: "barcodos"
     x: 100
     y: -100

  -  name: "barcotres"
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
