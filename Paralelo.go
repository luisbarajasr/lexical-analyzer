/*
Luis Antonio Barajas Ramírez A01235589
Diego Alejandro Hernández Romero A01198079
Eduardo Antonio López Vicencio A01411926

Esta propuesta contiene el mismo algoritmo que el código secuencial únicamente cambia el orden en cómo analizamos los archivos.
El escáner utiliza un DFA para realizar el análisis léxico, el cual está representado por la matriz “MT”, que contiene transiciones de estado basadas en el carácter de entrada.
La matriz tiene 12 filas, que representan los estados del DFA, y 18 columnas, que representan los diferentes símbolos de entrada. Los estados se identifican con números del 0 al 11,
y los símbolos de entrada se identifican con etiquetas. Cada elemento de la matriz representa el estado al que se transiciona cuando se encuentra un determinado símbolo de entrada en
un estado específico. Seguido, el código define un mapa llamado “reservedWords” que almacena palabras reservadas en un lenguaje de programación de C++. Cada palabra reservada se asigna
como clave en el mapa con el valor “true”. Esto permite verificar fácilmente si una palabra dada es reservada utilizando el mapa.
La función “filter(c string) int” recibe un carácter “c” como entrada y devuelve un entero que representa una categoría o tipo específico de carácter. La función se utiliza para clasificar
los diferentes tipos de caracteres presentes en el código fuente.
La función “scaner(line string) string” lee caracter por caracter el contenido dentro del archivo que se está analizando. Este carácter se envía a la ya mencionada función “filter(c string)”
en donde se evalúa qué carácter es. El estado del DFA dependerá de que esté regresando la función “filter”. El estado tendrá que llegar a un estado de aceptación en donde el token debe ser mayor
de 100, y si no, seguirá iterando por los caracteres hasta llegar a un estado de aceptación. En esta misma función se estará creando el string de retorno que contendrá los tags de HTML.
La diferencia está en la función main() ya que cambia la manera y el orden en cómo analizamos ya que lo hacemos de una forma concurrente y paralela utilizando las “goroutines”. En la función
“go func(file string)” se ejecutan múltiples goroutines ya que se encuentra dentro de un for loop. Esto hace que cada goroutine procesa un archivo de manera independiente, lo que permite que las
tareas se ejecuten en paralelo. Se utilizan variables de tipo “sync.WaitGroup” para esperar a que un grupo de goroutines terminen su ejecución antes de avanzar en el programa. Se usan los metodos
“add()” para agregar una goroutine dentro del grupo, “done()” para hacer “informar” que ya termino un goroutine  y “wait()” para “bloquear” el programa y esperar a que terminen los goroutines.

Complejidad
En realidad, la complejidad sigue siendo la misma al código secuencial ya que por cada archivo “n” se tiene que iterar por sus “m” caracteres. Termina teniendo una complejidad:
O( número de archivos * número de caracteres dentro de archivo) = O(n*m) ó O(n2). Pensamos que al ser una complejidad cuadrática debería ser lenta o tardada, sin embargo, al utilizar concurrencia
y paralelismo, podemos ver el tiempo en el que se ejecuta es relativamente rápido a como se esperaría de su complejidad. Esto quiere decir que en casos de programas paralelos y concurrentes, en
términos de complejidad algorítmica la complejidad se comporta como secuencial pero en cuestión de tiempo de ejecución, esta disminuye a como se pensaría.
*/ 
package main

import (
	// "bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode"
	"strings"
	"time"
	"sync"
)

const (
	INT     = 100
	REAL    = 101
	VAR     = 102
	OP      = 103
	PUNTOS  = 104
	MULT    = 105
	DIV     = 106
	POT     = 107
	STRING  = 108
	ABIERTO = 109
	CERRADO = 110
	COMM    = 111
	ENT     = 112
	EXP     = 113
	END     = 114
	COMML   = 115
	ERR     = 200
)

//Este arreglo representa un DFA. 

var MT = [][]int{
//     0   ,  1    ,   2   ,  3    , 4     ,  5    ,  6    ,    7   ,    8   ,  9    , 10    ,   11  ,  12   , 13    ,   14  , 15  ,   16    , 17
//    dig  ,  +    , : ; . ,  #    ,  /    , pot   ,  "    ,    (   ,    )   ,    // ,  /n   , esp   ,    .  , _     ,   *   , var ,    $    , odd
	{1     , OP    , 11     , 6     , 6     , POT  ,  10   , ABIERTO, CERRADO,   7   , ENT   , 0     , 2     ,   5   ,  OP   , 5   ,    END  , 4   }, //Estado 0 = inicial 
	{1     , INT   , INT   , INT   , INT   , INT   , INT   , INT    , INT    , INT   , INT   , INT   , 2     ,   5   ,    8  , 4   ,   INT   , INT }, //Estado 1 = Enteros
	{3     , ERR   , ERR   , ERR   , ERR   , ERR   , ERR   , ERR    , ERR    , ERR   , ERR   , 4     , ERR   ,  ERR  ,   ERR , PUNTOS ,    ERR  , 4   }, //Estado 2 = primer float
	{3     , REAL  , REAL  , REAL  , REAL  , REAL  , REAL  , REAL   , REAL   , REAL  , REAL  , REAL  , 4     ,  REAL ,   8   , REAL,    REAL , REAL}, //Estado 3 = remaning float
	{ERR   , ERR   , ERR   , ERR   , ERR   , ERR   , ERR   , ERR    , ERR    , ERR   ,  ERR  , ERR   , ERR   ,  ERR  ,   ERR , ERR ,   ERR   , 4   }, //Estatdo 4 = Error
	{5     , VAR   , VAR   , VAR   , VAR   , VAR   , VAR   , VAR    , VAR    , VAR   , VAR   , VAR   , VAR   ,   5   ,   5   , 5   ,   VAR   , VAR }, //Estado 5 = variable
	{DIV   , 4     , 4     ,  4    , 7     , 4     , 4     , DIV    , DIV    ,   7   , DIV   , DIV   , DIV   ,   4   ,   8   , DIV  ,   DIV  , DIV }, // Estado 6 = division o guion
	{7     , 7     , 7     , COMM  , 7     , 7     , 7     , 7      ,  7     ,   7   , COMM  , 7     ,   7   ,   7   ,   8   , 7    ,  COMM  , COMM  }, //Estado 7 = comentario
	{  8   ,  8    ,   8   ,   8   ,   8   ,   8   ,   8   ,   8    ,   8    ,  8    ,   8   ,   8   ,   8   ,   8   ,   9   ,  8   ,   8    ,  8  }, //Estado 8= comentario largo 1 
	{  8   ,  8    ,   8   ,   8   , COMML ,   8   ,   8   ,   8    ,   8    ,  8    ,   8   ,   8   ,   8   ,   8   ,   8   ,  8   ,  8     ,  8  }, //Estado 9= comentario largo 2 
	{   10 ,  10   ,   10  ,  10   , 10    ,  10   ,STRING ,  10    ,  10    ,  10   , ERR   ,  10   , 10    ,   10  ,   10  ,  10  ,   10   , 10 }, //Estado 10 = strings 
	{PUNTOS, PUNTOS, PUNTOS, PUNTOS, PUNTOS, PUNTOS, PUNTOS, PUNTOS , PUNTOS , PUNTOS, PUNTOS, PUNTOS,   11  , PUNTOS, PUNTOS, PUNTOS, PUNTOS, PUNTOS}, //Estado 11 = PUNTOS
}

//Mapa de variables reservadas. Si llegaran a faltar, unicamente seria ingresarlas
var reservedWords = map[string]bool{
	"asm":       true,
	"double":    true,
	"new":       true,
	"switch":    true,
	"auto":      true,
	"else":      true,
	"operator":  true,
	"endl":      true,
	"template":  true,
	"break":     true,
	"enum":      true,
	"private":   true,
	"this":      true,
	"case":      true,
	"extern":    true,
	"protected": true,
	"throw":     true,
	"catch":     true,
	"float":     true,
	"public":    true,
	"try":       true,
	"char":      true,
	"for":       true,
	"register":  true,
	"typedef":   true,
	"class":     true,
	"friend":    true,
	"return":    true,
	"union":     true,
	"const":     true,
	"goto":      true,
	"short":     true,
	"unsigned":  true,
	"continue":  true,
	"if":        true,
	"signed":    true,
	"virtual":   true,
	"default":   true,
	"inline":    true,
	"sizeof":    true,
	"void":      true,
	"delete":    true,
	"int":       true,
	"static":    true,
	"volatile":  true,
	"do":        true,
	"long":      true,
	"struct":    true,
	"while":     true,
	"cin":       true,
	"cout":      true,
	"include":   true,
	"iostream":  true,
	"namespace": true,
	"std":       true,
	"using":     true,
	"cstdlib":   true,
	"ctime":     true,
	"bool":      true,
}

//La funcion filter es constante ya que unicamente lo que hace es comparar una variable y regresa un entero. 
// No tieneNigun metodo que tome mas de O(1)
func filter(c string) int {
    if c == "0" || c == "1" || c == "2" || c == "3" || c == "4" || c == "5"|| c == "6" || c == "7" || c == "8" || c == "10" {
        return 0
    } else if c == "+" || c == "|" || c == "<" || c == ">" || c == "!" || c == "&" || c == "-" || c == "%" || c == "#" || c == "=" { //< > != == | | & + - / % *
        return 1
    } else if c == ";" || c == ":" || c == "," {
        return 2
    } else if c == "#" {
        return 3
    } else if c == "/" {
        return 4
    } else if c == "^" {
        return 5
    } else if c == "\"" || c == "\"\""   {
        return 6
    } else if c == "(" || c == "{" || c == "[" { 
        return 7
    } else if c == ")" || c == "}" || c == "]" {
        return 8
    } else if c == "\r" {
        return 10
    } else if c == " " {
        return 11
    } else if c == "." {
        return 12
    } else if c == "_" {
        return 13
    } else if c == "*"{
		return 14
	}else if strings.Contains("abcdefghijklmnopqrstuvwxyz", strings.ToLower(string(c))) {
        return 15
    } else if c == "$" {
        return 16
    } else {
        return 17
    }
}

/* 
Complejidad: O(m), en donde m = el numero de caracteres dentro del archivo o del string
*/
func scaner(line string) string{
	highlighted := "" //string de los tags HTML

	state := 0   // estado en el DFA
	lexeme := "" // string que genera el token
	tokens := []int{}

	read := true // mientras el state no sea ACCEPT ni ERROR
	lineIndex := 0 //indice necesario para iterar sobre el string (contenido del archivo)
	c := ""

	for { //O(n) ? 
		
		//Mientras no este en un estado de Aceptacion y su indice sea menor a la cantidad de carateres detro del archivo
		for state < 100 && lineIndex < len(line) { //O(n)
	
			if read {
			c = string(line[lineIndex])

			lineIndex += 1 

			} else {
				read = true
			}

			state = MT[state][filter(c)] //constante O(1),
			if state < 100 && state != 0 {
				lexeme += c
			}

		}
		
		if state == INT{
			read = false
			highlighted += "<span style=\"color: rgb(186,205,171)\">" + lexeme + "</span> "

		} else if state == REAL{
			
			read = false
			highlighted += "<span style=\"color: rgb(186,205,171)\">" + lexeme + "</span> "

		} else if state == OP{
			lexeme += c
			highlighted += "<span style=\"color: red\">" + lexeme + "</span> "

		}else if state == ABIERTO{
			lexeme += c
			highlighted += "<span style=\"color: rgb(204,118,2010)\">" + lexeme + "</span> "

		}else if state == CERRADO{
			lexeme += c
			highlighted += "<span style=\"color: rgb(204,118,2010)\">" + lexeme + "</span> "

		} else if state == STRING{
			lexeme += c
			highlighted += "<span style=\"color: rgb(1107,148,124)\">" + lexeme + "</span> "

		}else if state == DIV{
			read = false
			highlighted += "<span style=\"color: red\">" + lexeme + "</span> "

		}else if state == COMM{
			lexeme += c
			read = false
			highlighted += "<span style=\"color: rgb(116,152,103)\">" + lexeme + "</span> "
			
		}else if state == COMML{
			lexeme += c
			highlighted += "<span style=\"color: rgb(116,152,103)\">" + lexeme + "</span> "
			
		}else if state == PUNTOS{
			read = false
			highlighted += "<span style=\"color: black\">" + lexeme + "</span> "

		}else if state == MULT{
			lexeme += c

		}else if state == POT{
			lexeme += c

		}else if state == VAR{
			
			read = false
			if unicode.IsDigit(rune(lexeme[0])){
				
				state = ERR
			}else if reservedWords[lexeme]{ //Tiempo constante
				//aqui se revisar si la variable es palabra reservada o no, busca en el string "reserved words"

				highlighted += "<span style=\"color: blue\">" + lexeme + "</span> "

			}else{		
				highlighted += "<span style=\"color: #F1C376\">" + lexeme + "</span> " //si no es rservada es una variable
			}
			
		}else if state == EXP{
			
			if strings.ToLower(c) == "e"{

				highlighted += "<span style=\"color: #F1C376\">" + lexeme + "</span> "
				state = VAR
			}
		}else if state == ERR{
			
		}else if state == ENT{ 
		//Cuando el estado es ENT, se refiere a que hay un salto de linea, por lo que se le una el tag de HTML <br> para hacer un salto de linea
			highlighted += "<br>"

		}
		if state != 201{ 
			tokens = append(tokens, state)
		}

		if lineIndex == len(line){
			return highlighted
		}

		lexeme = ""
		state = 0
	}
}

/*
O(numero_de_archivos * caracteres_dentro_acrhivo)
*/
func main() {

	path := os.Args[1]

	var temp string =  "\\"
	
	files, err := filepath.Glob(filepath.Join(path, "*.txt")) //se abre la carpeta
	if err != nil {
		log.Fatal(err)
	}

	if path[len(path)-1] != temp[0]{
		path = path + string(temp[0])
	} 
	
	var wg sync.WaitGroup //Se inicia un grupo de espera para los Gorutines
	
	startTime := time.Now() // Empieza el cronometro

	//Este for loop creara multiples goroutines por cada file
	for _, file := range files { //Complejidad: O(n), en donde n = numero de archivos dentro de la carpeta
		

		wg.Add(1) //Se incrementa el tamaño del gurpo de espera en 1 

		go func(file string){ //Funcion anonima que crea un goroutine
			
			defer wg.Done() //Decrementa el tamaño del grupo de espera, significa que termina un goroutine

			content, err := os.ReadFile(file)
			if err != nil{
				log.Printf("Failed to read file %s: %s\n", file, err)
				return
			}
			
			outputDir := filepath.Dir(file) //Toma el path 
			baseName := filepath.Base(file) //Toma el nombre del file en el que esta
			outputFile := strings.TrimSuffix(baseName, filepath.Ext(baseName))+".html" //Crea el nombre del archivo HTML output
			outputPath := filepath.Join(outputDir, outputFile)  //Crea el path 
		
			outputFileHandle, err := os.Create(outputPath)
			if err != nil {
				log.Printf("Failed to create output file %s: %s\n", outputPath, err)
				return
			}
			defer outputFileHandle.Close() 
			
			htmlTags := scaner(string(content)) //Complejidad = O(m), en donde m = numero de caracteres dentro del archivo
			//Se guardan los html Tags que regresa la funcion scanner
	
			_, err = outputFileHandle.WriteString(htmlTags) //se crea el archivo 
			if err != nil {
				log.Printf("Failed to write HTML tags to file %s: %s\n", outputPath, err)
				return
			}
	
			fmt.Printf("Output file %s created successfully\n", outputPath) 
			
			}(file)
			
		}
		wg.Wait() 
		//Este llamado bloquea o "para temporalmente" el programa hasta que el grupo de goroutines llego a 0, es decir, 
		//Que todas las goroutines hayan terminado

		endTime := time.Now()             // Guarda el tiempo de finalizar
		elapsedTime := endTime.Sub(startTime) // Calcula el tiempo
		fmt.Printf("Program completed in %s\n", elapsedTime)
	

}