# Chat Encriptado de Extremo a Extremo (E2E) en Golang

<!-- ![Logo](logo.png) -->
**Un chat seguro y descentralizado, donde tu privacidad es lo primero.** 

Este proyecto implementa una aplicación de chat encriptado de extremo a extremo (E2E) utilizando Go y WebSockets. A diferencia de las aplicaciones de chat tradicionales, no hay un servidor central que almacene o tenga acceso a tus mensajes. Solo tú y los destinatarios previstos pueden leerlos.

## Características Principales

*   **Encriptación de Extremo a Extremo (E2E):** Los mensajes se cifran en el dispositivo del remitente y solo se descifran en el dispositivo del destinatario. Nadie más, ni siquiera la propia aplicación, puede leerlos.
*   **Seguridad:** Se utilizan algoritmos criptográficos robustos:
    *   RSA para el intercambio seguro de claves simétricas.
    *   AES para el cifrado de los mensajes.
*   **Comunicación en tiempo real:** Se utilizan WebSockets para una comunicación fluida e instantánea.
*   **Manejo seguro de claves:** Las claves privadas nunca se transmiten ni se almacenan de forma insegura.

## Arquitectura

El proyecto se divide en los siguientes paquetes principales:



*   `cmd/server`: Contiene el punto de entrada principal de la aplicación (`main.go`).
*   `config`: Gestiona la configuración de la aplicación, incluyendo las claves criptográficas (`config.go`). 
*   `internal/websocket`: Maneja la comunicación WebSocket.
    *   `client.go`: Gestiona las conexiones WebSocket individuales.
    *   `message.go`: Define las estructuras de los mensajes y payloads.
*   `pkg/chat`: Contiene la lógica del chat.
    *   `encryption.go`: Funciones de encriptación específicas del chat.
    *   `message.go`: Estructuras de mensajes del chat.
    *   `room.go`: Lógica de las salas de chat.
    *   `user.go`: Define la estructura `User`.
*   `pkg/crypto`: Contiene la implementación del cifrado.
    *   `aes.go`: Funciones para cifrado AES.
    *   `rsa.go`: Funciones para cifrado RSA.
*   `logger`: Contiene la lógica de logging de la aplicación.

## Descripción de los paquetes

*   **`cmd/`**: Contiene el punto de entrada principal de la aplicación (`main.go`). Este archivo se encarga de parsear los flags `-server` y `-client` para iniciar la aplicación en el modo correspondiente. La lógica principal de la aplicación reside en los paquetes dentro de `internal`.

*   **`config/`**: Contiene los archivos de configuración de la aplicación, como la configuración del servidor, las claves criptográficas, etc.

*   **`internal/`**: Contiene el código fuente interno de la aplicación. Este código *no debe ser importado por proyectos externos*.

    *   **`view/`**: Contiene la lógica de la interfaz de usuario, utilizando la librería Bubble Tea. El archivo `chat.go` define el modelo, la lógica de actualización (`Update`) y la vista (`View`) del chat.
    
    *   **`model/`**: Define las estructuras de datos que se utilizan para la comunicación a través de WebSockets.
        *   `chat.go`: Contiene los modelos relacionados con el estado del chat en la interfaz de usuario (`ModelChat`, `IncomingMessage`, etc.). `IncomingMessage` se utiliza para comunicar mensajes desde la capa de websocket a la capa de UI.
        *   `message.go`: Contiene los modelos que representan los mensajes que se intercambian a través del WebSocket (`TextMessagePayload`, `WebsocketMessage`, etc.). `TextMessagePayload` es la estructura de datos que se envia por el websocket. `WebsocketMessage` es un envoltorio que contiene el tipo de mensaje y el payload.
    
        *   `user.go`: Contiene el modelo que representa al estrcutura basica de un usuario logeado en la plataforma.


    *   **`websocket/`**: Contiene toda la lógica relacionada con la comunicación mediante WebSockets.

        *   **`client/`**: Lógica específica del cliente WebSocket, incluyendo:
            *   Establecimiento de la conexión WebSocket.
            *   Manejo de mensajes enviados y recibidos desde el cliente.
            *   Integración con la interfaz de usuario (en `internal/ui`).

        *   **`server/`**: Lógica específica del servidor WebSocket, incluyendo:
            *   Gestión de las conexiones de los clientes.
            *   Broadcast de mensajes a los clientes conectados.
            *   Manejo de eventos de conexión y desconexión.

        *   **`handler/`**: Manejadores para los diferentes eventos del WebSocket. Estos manejadores contienen la lógica para procesar los mensajes recibidos y realizar las acciones correspondientes (ej. encriptar, desencriptar, etc).
        
        *   `connection.go`: Maneja la creación, gestión y cierre de las conexiones WebSocket, tanto del lado del cliente como del servidor.


*   **`pkg/`**: Contiene librerías reutilizables *dentro del proyecto*.

    *   **`crypto/`**: Contiene las funciones de cifrado.
        *   **`aes/`**: Implementación del cifrado AES.
        *   **`rsa/`**: Implementación del cifrado RSA.

    *   **`logger/`**: Contiene la lógica para el manejo de logs.

### Diagrama de componentes 
```mermaid
graph LR
    subgraph "cmd"
        main["main.go"]
    end

    subgraph "config"
        config["config.go"]
    end

    subgraph "internal"
        subgraph "view"
            chat["chat.go"]
        end
        subgraph "model"
            userModel["user.go"]
            encryptionModel["encryption.go"]
            roomModel["room.go"]
            messageModel["message.go"]
        end
        subgraph "websocket"
            subgraph "client"
                client["client.go"]
                handlerClient["handler.go"]
            end
            subgraph "server"
                server["server.go"]
                handlerServer["handler.go"]
            end
            connection["connection.go"]
        end
    end

    subgraph "pkg"
        subgraph "crypto"
            subgraph "aes"
                aes["aes.go"]
            end
            subgraph "rsa"
                rsa["rsa.go"]
            end
        end
        logger["logger.go"]
    end

    main --> config
    main --> logger
    main --> client
    main --> server

    client --> connection
    client --> chat
    client --> messageModel

    handlerClient --> messageModel
    handlerClient --> userModel
    handlerClient --> roomModel
    handlerClient --> encryptionModel
    handlerClient --> aes
    handlerClient --> rsa
    handlerClient --> connection

    server --> connection
    handlerServer --> messageModel
    handlerServer --> userModel
    handlerServer --> roomModel
    handlerServer --> encryptionModel
    handlerServer --> aes
    handlerServer --> rsa
        handlerServer --> connection

    chat --> messageModel
    chat --> userModel
    chat --> roomModel
    chat --> client

    style cmd fill:#ccf,stroke:#888,stroke-width:2px
    style config fill:#ccf,stroke:#888,stroke-width:2px
    style internal fill:#bbf,stroke:#666,stroke-width:2px
    style view fill:#dde,stroke:#777,stroke-width:2px
    style websocket fill:#aae,stroke:#444,stroke-width:2px
    style client fill:#ccf,stroke:#888,stroke-width:2px
    style server fill:#ccf,stroke:#888,stroke-width:2px
    style model fill:#dde,stroke:#777,stroke-width:2px
    style pkg fill:#99f,stroke:#222,stroke-width:2px
    style crypto fill:#bbf,stroke:#666,stroke-width:2px
    style aes fill:#ccf,stroke:#888,stroke-width:2px
    style rsa fill:#ccf,stroke:#888,stroke-width:2px
```

## Flujos de funcionamiento

### Establecimiento de una conexión segura (Handshake)

El establecimiento de una conexión segura, o *handshake*, es crucial para el cifrado de extremo a extremo. Este proceso permite el intercambio seguro de la clave simétrica que se utilizará para cifrar los mensajes. En una arquitectura con servidor central, el servidor actúa como intermediario para este intercambio.

**Escenario:** El usuario A quiere iniciar una conversación con el usuario B (o crear un nuevo grupo).

**Caso 1: Intercambio inicial de claves públicas (A y B no se han comunicado antes)**

Este caso ocurre cuando A y B se conectan por primera vez.

1.  **Conexión WebSocket (Cliente A y Cliente B):** Los clientes de A y B (en `internal/websocket/client/client.go`) establecen conexiones WebSocket *independientes* con el servidor (`internal/websocket/server/server.go`).

2.  **Intercambio de claves públicas (a través del servidor):**

    a.  **Envío de clave pública (Cliente A -> Servidor):** El cliente de A genera su par de claves RSA (en `pkg/crypto/rsa/rsa.go`) si no lo tiene ya, y envía un mensaje de tipo `publicKeyExchange` (definido en `internal/websocket/model/message.go`) que contiene su clave pública RSA al *servidor* a través de su conexión WebSocket. El handler del cliente (`internal/websocket/client/handler.go`) se encarga de formatear el mensaje y enviarlo.

    b.  **Recepción y retransmisión de clave pública (Handler Servidor):** El handler del servidor (`internal/websocket/server/handler.go`) recibe el mensaje de A, identifica al destinatario B, y retransmite la clave pública de A al cliente de B a través de la conexión WebSocket de B.

    c.  **Recepción y almacenamiento de clave pública (Handler Cliente B):** El handler del cliente de B (`internal/websocket/client/handler.go`) recibe la clave pública de A y la almacena en el modelo de usuario (`internal/model/user.go`), asociándola con el ID de A.

    d.  **Envío de clave pública (Cliente B -> Servidor):** Similar al paso a, pero B envía su clave pública al servidor.

    e.  **Recepción y retransmisión de clave pública (Handler Servidor):** El handler del servidor recibe el mensaje de B y lo retransmite a A.

    f.  **Recepción y almacenamiento de clave pública (Handler Cliente A):** Similar al paso c, pero A recibe y almacena la clave pública de B.

3.  **Generación de clave simétrica (Cliente A):** El cliente de A genera una clave simétrica aleatoria (AES) utilizando `pkg/crypto/aes/aes.go`.

4.  **Cifrado de clave simétrica (Cliente A):** El cliente de A cifra la clave simétrica generada con la clave pública de B utilizando RSA (`pkg/crypto/rsa/rsa.go`).

5.  **Envío de clave simétrica cifrada (Cliente A -> Servidor -> Cliente B):** El cliente de A envía la clave simétrica cifrada al *servidor*, quien a su vez la retransmite al cliente de B. Esto se realiza dentro de un mensaje `inviteToGroup` (para la creación de grupos) u otro mensaje similar (para chats directos), manejado por el handler.

6.  **Recepción de clave simétrica cifrada (Handler Cliente B):** El handler del cliente de B recibe el mensaje del servidor.

7.  **Descifrado de clave simétrica (Cliente B):** El cliente de B descifra la clave simétrica recibida utilizando *su* clave privada RSA (`pkg/crypto/rsa/rsa.go`).

8.  **Almacenamiento de clave simétrica (Cliente B):** El cliente de B guarda la clave simétrica descifrada en el modelo de usuario (`internal/model/user.go`), asociándola con el ID de la conversación o del grupo.

9.  **Comunicación cifrada:** A partir de este momento, A y B pueden comunicarse de forma segura utilizando la clave simétrica compartida. Los mensajes cifrados se enviarán a través del servidor.

**Caso 2: Claves públicas ya almacenadas (A y B ya se han comunicado antes)**

Este caso ocurre cuando A y B ya se han comunicado previamente y tienen las claves públicas del otro almacenadas en su configuración local. El flujo es similar al caso 1, pero se omiten los pasos de intercambio de claves públicas (2a-2f).

1.  **Conexión WebSocket (Cliente A y Cliente B):** Los clientes de A y B establecen una conexión WebSocket con el servidor.

2.  **Recuperación de claves públicas (Cliente A y Cliente B):** Los clientes de A y B recuperan las claves públicas del otro desde el modelo de usuario (`internal/model/user.go`).

3.  **Generación, cifrado, envío, recepción, descifrado y almacenamiento de clave simétrica:** Los pasos 3-8 del Caso 1 aplican aquí.

4.  **Comunicación cifrada:** A partir de este momento, A y B pueden comunicarse de forma segura utilizando la clave simétrica compartida. Los mensajes cifrados se enviarán a través del servidor.

### Envío de un mensaje

Este flujo, en esencia, no cambia mucho con la introducción del servidor, solo que los mensajes ahora pasan por él:

1.  **Escritura del mensaje (Vista Chat):** El usuario A escribe un mensaje en la interfaz de usuario (`internal/view/chat.go`).

2.  **Obtención de la clave simétrica (Vista Chat):** La vista del chat (`internal/view/chat.go`) obtiene la clave simétrica del modelo de usuario (`internal/model/user.go`).

3.  **Cifrado del mensaje (Cliente A):** El cliente de A (`internal/websocket/client/client.go`) cifra el mensaje con la clave simétrica utilizando AES (`pkg/crypto/aes/aes.go`).

4.  **Envío del mensaje cifrado (Cliente A -> Servidor -> Cliente B):** El cliente de A envía el mensaje cifrado al *servidor*, quien a su vez lo retransmite al cliente de B. El handler del cliente (`internal/websocket/client/handler.go`) se encarga de formatear el mensaje para ser enviado.

5.  **Recepción del mensaje cifrado (Handler Cliente B):** El handler del cliente de B (`internal/websocket/client/handler.go`) recibe el mensaje cifrado *del servidor*.

6.  **Obtención de la clave simétrica (Cliente B):** El cliente de B obtiene la clave simétrica del modelo de usuario (`internal/model/user.go`).

7.  **Descifrado del mensaje (Cliente B):** El cliente de B descifra el mensaje con la clave simétrica utilizando AES (`pkg/crypto/aes/aes.go`).

8.  **Visualización del mensaje (Vista Chat):** El cliente de B envía el mensaje descifrado a la vista del chat (`internal/view/chat.go`) para que se muestre al usuario B.



## Cómo ejecutar la aplicación

#### Necesitarás ejecutar múltiples instancias en diferentes terminales para simular varios usuarios

1.  Clona el repositorio: `git clone https://github.com/osmancadc/go-encrypted-chat.git`
2.  Navega al directorio del proyecto: `cd go-ecrypted-chat`
3.  Construye y ejecuta la aplicación: `./scripts/start.sh [-server] [-client] [username]`


## Contribuciones

Este proyecto está en constante evolución y siempre hay espacio para mejorar. Creo que la colaboración es la mejor forma de aprender y crecer juntos.

Si te interesa aprender sobre desarrollo en Go, WebSockets, criptografía o simplemente quieres contribuir a un proyecto open source, te invito a participar. Desde la corrección de errores hasta la implementación de nuevas funcionalidades, tu contribución es bienvenida. 

¡Juntos podemos hacer este proyecto aún mejor!

## Licencia

[MIT License](https://opensource.org/licenses/MIT)