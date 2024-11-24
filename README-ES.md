# Desarrollo del Proyecto de Estación Meteorológica

## 1. Introducción
Este proyecto tiene como objetivo desarrollar un sistema de estaciones meteorológicas capaces de recopilar datos ambientales, procesarlos y transmitirlos a un servidor central para su análisis. Durante el proceso de desarrollo, se han implementado soluciones técnicas y mejoras para garantizar un sistema robusto y eficiente.

## 2. Descripción del Proyecto
El sistema consta de dos componentes principales:

### Estaciones Meteorológicas
Las estaciones están equipadas con sensores para medir variables climáticas como temperatura, humedad y presión. Utilizan una placa que integra un chip Arduino Mega y un ESP8266, descrito en [este enlace](https://www.luisllamas.es/arduino-mega-esp8266-en-un-unico-dispositivo/). Este hardware combina la capacidad de procesamiento del Arduino Mega con la conectividad WiFi del ESP8266, lo que permite la transmisión eficiente de datos al servidor.

El repositorio incluye versiones simplificadas del código para la conexión entre Arduino y el servidor, ubicadas en la carpeta `station/example/`. Estas versiones son útiles para pruebas iniciales y para comprender los fundamentos del sistema de comunicación.

### Servidor Central
El servidor está implementado en Go, utilizando una arquitectura de **workers** con una **job queue** para procesar tareas de forma concurrente. Esta estructura permite que múltiples estaciones se comuniquen simultáneamente sin afectar el rendimiento.

El servidor está basado en Docker y utiliza una **Network** privada para conectarse de manera segura con una base de datos MySQL. Este diseño asegura una gestión eficiente de los datos recibidos y procesados.

### Base de Datos
La base de datos está estructurada en tres tablas principales: `station`, `measurement` y `entry`. Estas tablas se crean automáticamente mediante un script bash que ejecuta código SQL ubicado en `server/docker/mysql/tables-init.sh`. La estructura es la siguiente:

- **Tabla `station`**:
  - Contiene información sobre las estaciones registradas.
  - Columnas:
    - `station_id`: Identificador único, auto-incremental.
    - `station_owner`: Nombre del propietario de la estación (requerido).
    - `latitude`: Latitud de la estación, validada para estar entre -90 y 90.
    - `longitude`: Longitud de la estación, validada para estar entre -180 y 180.
    - `password_hash`: Hash de la contraseña de la estación, requerido para autenticación.

- **Tabla `measurement`**:
  - Contiene mediciones específicas enviadas por las estaciones.
  - Columnas:
    - `measurement_id`: Identificador único, auto-incremental.
    - `random_num`: Número aleatorio generado por las estaciones durante pruebas iniciales.

- **Tabla `entry`**:
  - Registra entradas individuales enviadas por las estaciones, incluyendo ubicación, tiempo y mediciones asociadas.
  - Columnas:
    - `entry_id`: Identificador único, auto-incremental.
    - `station_id`: Relación con la estación que envió la entrada.
    - `latitude`: Latitud de la entrada, validada para estar entre -90 y 90.
    - `longitude`: Longitud de la entrada, validada para estar entre -180 y 180.
    - `measurement_id`: Relación con la medición registrada.
    - `entry_time`: Marca de tiempo indicando cuándo se recibió la entrada.
  - Relaciones:
    - `station_id` es una clave foránea que referencia a la tabla `station`.
    - `measurement_id` es una clave foránea que referencia a la tabla `measurement`.
  - **Razón para duplicar las coordenadas de la estación**: Esta duplicación asegura que los datos históricos no se vean afectados si una estación cambia de ubicación. Las entradas antiguas mantienen las coordenadas originales, evitando inconsistencias en los datos almacenados.

## 3. Proceso de Desarrollo

### Inicio del Proyecto
El desarrollo comenzó con dos enfoques principales:
1. **Configuración del Servidor**: Se construyó utilizando [Go](https://go.dev/learn/) y [Docker](https://docs.docker.com/get-started/), estructurado para manejar datos meteorológicos de forma eficiente. Se utilizó una **Network** privada para conectar el servidor y la base de datos [MySQL](https://dev.mysql.com/doc/).
2. **Estaciones Meteorológicas**: Se seleccionó una placa con un chip Arduino Mega y ESP8266 como base para la comunicación y recolección de datos.

### Desarrollo Inicial
Para validar la infraestructura inicial, se utilizó la herramienta [Packet Sender](https://packetsender.com/) para enviar datos simulados al servidor. Esto permitió probar la recepción, interpretación y almacenamiento de datos en condiciones controladas. Una vez que el servidor mostró estabilidad, las estaciones comenzaron a generar números aleatorios como datos simulados para realizar pruebas adicionales.

### Refinamiento y Expansión
El proyecto evolucionó con nuevas características:
1. **Estructuración del Código**: El servidor se organizó en módulos independientes para una mayor claridad y escalabilidad.
2. **Validación de Datos**: En lugar de encriptación, se priorizó una validación estricta de los datos transmitidos por las estaciones para garantizar su integridad.
3. **Interacciones Locales**: Las estaciones comenzaron a utilizar módulos SD para almacenar configuraciones de red y datos críticos.

## 4. Estructura Final del Proyecto

### Árbol del Proyecto
```
├── server 
│ ├── data-server 
│ ├── docker 
│ ├── safety 
│ ├── codec 
│ ├── tcp 
│ └── db-handle 
├── station 
│ ├── code 
│ │ ├── arduino 
│ │ ├── esp8266 
│ │ └── common 
│ ├── example 
│ │ ├── arduino-mega-station.ino 
│ │ └── esp8266-station.ino 
│ ├── libraries 
│ └── sd-content  
```
## 5. Cómo Usar

### Requisitos Previos
1. Instalar [Git](https://git-scm.com/).
2. Instalar [Docker](https://docs.docker.com/get-docker/).
3. Instalar [Arduino IDE](https://www.arduino.cc/en/software).
4. Clonar el repositorio:
```
git clone https://github.com/AmogusAzul/weather-station
cd weather-station
```

5. Configurar el servidor utilizando Docker:

```
cd server/docker docker-compose up
```
6. Cargar y subir el código a tu estación:
- Abre `station/code/arduino/arduino.ino` o `station/code/esp8266/esp8266.ino` en Arduino IDE.
- Ajusta la configuración de red y los tokens.
- Sube el código a la placa correspondiente.

## 6. Próximos Pasos

### Expansión del Servidor
- **Soporte HTTP**: Se integrará un sistema para visualizar gráficas y tablas directamente en un navegador.
- **Descarga de Datos**: Los usuarios podrán descargar informes en formatos como CSV o PNG.

### Desarrollo de las Estaciones
1. **Selección de Sensores**: Elegir sensores específicos para medir temperatura, humedad, presión y otras variables ambientales.
2. **Diseño de la Carcasa**: Crear una estructura resistente para proteger los componentes electrónicos de las estaciones en condiciones adversas.
3. **Optimización del Firmware**: Mejorar el rendimiento del código para manejar datos en tiempo real de manera eficiente.

### Pruebas y Validación
- Se realizarán pruebas de estrés para garantizar que el sistema pueda manejar múltiples estaciones conectadas simultáneamente.
- Se desarrollará una guía detallada para facilitar la implementación y el uso del sistema.

## 7. Conclusión
El proyecto de estación meteorológica ha avanzado significativamente, implementando una infraestructura técnica robusta y superando desafíos de diseño y funcionalidad. Los próximos pasos garantizarán un sistema eficiente para recopilar y analizar datos meteorológicos, proporcionando herramientas de visualización accesibles y útiles para los usuarios.

### Repositorio
El código para este proyecto está disponible en este [repositorio de GitHub](https://github.com/AmogusAzul/weather-station).
