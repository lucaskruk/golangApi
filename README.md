# locator

## Locator API

Locator is a Golang API that returns 
the position and content from the source of a SOS message in a two dimensional surface. it does it by triangling its position, from three listening ships who receive the message. It also completes the content of the message between all the listeners.
The position of the listener ships is defined below:

- BarcoUno: [-500, -200]

- BarcoDos: [100, -100]

- BarcoTres: [500, 100]


The application has the following functions:
```golang
// input: distance from listening ship to missing ship
// output: X and Y coordinates of missing ship
func GetLocation(distances ...float32) (x, y float32)
// input: message as received by each listener ship
// output: complete SOS message
func GetMessage(messages ...[]string) (msg string)
```
#### Considerations:

- All parameters from GetLocation function have the same measure unit, and that is the same used at X and Y axis to determine position.

- The SOS message is received in an array of strings.

- When a word in the SOS message cannot be determined is replaced by an empty string in message array.

○ Example: [“este”, “es”, “”, “mensaje”]

- Consider that a time displacement could occur in some of the messages, so the first elements could appear as empty, and the length of the array be longer than the original message.

○ Example:

■ BarcoUno: [“”, “este”, “es”, “un”, “mensaje”]

■ BarcoDos: [“este”, “”, “un”, “mensaje”]

■ BarcoTres: [“”, ””, ”es”, ””, ”mensaje”]

A /helpme/ endpoint is available to obtain missing ship location and message. If it is not able to determine position or message, it returns 404

Also /helpme_split/ endpoint is available, This endpoint can receive POST and GET requests. 
With POST request we can send each listener ship message and distance.
With GET request we will get the position and message with the same format as /helpme/ endpoint, and the same 404 behavior. If not all three listener ships where sent before issuing a GET request, it will response an error message telling that there isn't enough information to proceed. 

In order to run the application, we need to have our GO environment setup properly, to know how to install please go to: 
https://golang.org/doc/install

You will need to install the following packages:
"github.com/gorilla/mux"
"gopkg.in/yaml.v2"
They're automatically downloaded during go build command, but also can be manually downloaded with the following commands:
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2


To start the API, execute the following command:
go run locator.go


You can modify /config/config.yml to change application port or listening ships location.

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
