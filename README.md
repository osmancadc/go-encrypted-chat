# Chat Encriptado de Extremo a Extremo (E2E) en Golang

<!-- ![Logo](logo.png) -->
**Un chat seguro y descentralizado, donde tu privacidad es lo primero.** 

Este proyecto implementa una aplicación de chat encriptado de extremo a extremo (E2E) utilizando Go y WebSockets. A diferencia de las aplicaciones de chat tradicionales, no hay un servidor central que almacene o tenga acceso a tus mensajes. Solo tú y los destinatarios previstos pueden leerlos.

## Características Principales

*   **Encriptación de Extremo a Extremo (E2E):** Los mensajes se cifran en el dispositivo del remitente y solo se descifran en el dispositivo del destinatario. Nadie más, ni siquiera la propia aplicación, puede leerlos.
*   **Descentralización:** No hay un servidor central. Cada cliente actúa como su propio servidor WebSocket, estableciendo conexiones directas con otros clientes.
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

![Diagrama de arquitectura](arquitectura.png)

## Flujos de funcionamiento

### Establecimiento de una conexión segura (Handshake)

El establecimiento de una conexión segura, o *handshake*, es crucial para el cifrado de extremo a extremo. Este proceso permite el intercambio seguro de la clave simétrica que se utilizará para cifrar los mensajes.

**Escenario:** El usuario A quiere iniciar una conversación con el usuario B (o crear un nuevo grupo).

**Caso 1: Intercambio inicial de claves públicas (A y B no se han comunicado antes)**

Este caso ocurre cuando A y B se conectan por primera vez.

1.  **Conexión WebSocket:** A y B establecen una conexión WebSocket entre sí.

2.  **Intercambio de claves públicas:**

    *   a. **A envía su clave pública a B:** A **envía** un mensaje especial de tipo `publicKeyExchange` que contiene su clave pública RSA a B.
    *   b. **B recibe la clave pública de A:** B **recibe** el mensaje y **guarda** la clave pública de A en su configuración local, asociándola con el ID de A.
    *   c. **B envía su clave pública a A:** B **envía** un mensaje de tipo `publicKeyExchange` que contiene su clave pública RSA a A.
    *   d. **A recibe la clave pública de B:** A **recibe** el mensaje y **guarda** la clave pública de B en su configuración local, asociándola con el ID de B.

3.  **Generación de clave simétrica (A):** A **genera** una clave simétrica aleatoria (AES).

4.  **Cifrado de clave simétrica (A):** A **cifra** la clave simétrica generada con la clave pública de B utilizando RSA.

5.  **Envío de clave simétrica cifrada (A -> B):** A **envía** la clave simétrica cifrada a B a través de WebSocket, dentro de un mensaje `inviteToGroup` (para la creación de grupos) u otro mensaje similar (para chats directos).

6.  **Recepción de clave simétrica cifrada (B):** B **recibe** la clave simétrica cifrada.

7.  **Descifrado de clave simétrica (B):** B **descifra** la clave simétrica recibida utilizando *su clave privada RSA*.

8.  **Almacenamiento de clave simétrica (B):** B **guarda** la clave simétrica descifrada en su configuración local, asociándola con el ID de la conversación o del grupo.

9.  **Comunicación cifrada:** A partir de este momento, A y B **pueden comunicarse de forma segura** utilizando la clave simétrica compartida para cifrar y descifrar los mensajes con AES.

**Caso 2: Claves públicas ya almacenadas (A y B ya se han comunicado antes)**

Este caso ocurre cuando A y B ya se han comunicado previamente y tienen las claves públicas del otro almacenadas en su configuración local.

1.  **Conexión WebSocket:** A y B establecen una conexión WebSocket entre sí.

2.  **Recuperación de claves públicas:**

    *   a. **A recupera la clave pública de B:** A **recupera** la clave pública de B de su configuración local.
    *   b. **B recupera la clave pública de A:** B **recupera** la clave pública de A de su configuración local.

3.  **Generación de clave simétrica (A):** A **genera** una clave simétrica aleatoria (AES).

4.  **Cifrado de clave simétrica (A):** A **cifra** la clave simétrica generada con la clave pública de B utilizando RSA.

5.  **Envío de clave simétrica cifrada (A -> B):** A **envía** la clave simétrica cifrada a B a través de WebSocket, dentro de un mensaje `inviteToGroup` (para la creación de grupos) u otro mensaje similar (para chats directos).

6.  **Recepción de clave simétrica cifrada (B):** B **recibe** la clave simétrica cifrada.

7.  **Descifrado de clave simétrica (B):** B **descifra** la clave simétrica recibida utilizando *su clave privada RSA*.

8.  **Almacenamiento de clave simétrica (B):** B **guarda** la clave simétrica descifrada en su configuración local, asociándola con el ID de la conversación o del grupo.

9.  **Comunicación cifrada:** A partir de este momento, A y B **pueden comunicarse de forma segura** utilizando la clave simétrica compartida para cifrar y descifrar los mensajes con AES.

### Envío de un mensaje

Este flujo describe el proceso de envío de un mensaje una vez que se ha establecido la conexión segura.

1.  **Escritura del mensaje (A):** El usuario A **escribe** un mensaje.
2.  **Obtención de la clave simétrica (Cliente de A):** El cliente de A **obtiene** la clave simétrica del grupo desde su configuración local.
3.  **Cifrado del mensaje (Cliente de A):** El cliente de A **cifra** el mensaje con la clave simétrica utilizando AES.
4.  **Envío del mensaje cifrado (Cliente de A -> Cliente de B):** El cliente de A **envía** el mensaje cifrado a través de WebSockets al cliente de B (y a otros miembros del grupo, si aplica).
5.  **Recepción del mensaje cifrado (Cliente de B):** El cliente de B **recibe** el mensaje cifrado.
6.  **Obtención de la clave simétrica (Cliente de B):** El cliente de B **obtiene** la clave simétrica del grupo desde su configuración local.
7.  **Descifrado del mensaje (Cliente de B):** El cliente de B **descifra** el mensaje con la clave simétrica utilizando AES.
8.  **Visualización del mensaje (B):** El mensaje descifrado se **muestra** al usuario B.

![Diagrama de mensaje](flujo.png)

## Cómo ejecutar la aplicación

#### Necesitarás ejecutar múltiples instancias en diferentes terminales para simular varios usuarios

1.  Clona el repositorio: `git clone https://github.com/osmancadc/go-encrypted-chat.git`
2.  Navega al directorio del proyecto: `cd go-ecrypted-chat`
3.  Construye y ejecuta la aplicación: `./scripts/start.sh`


## Contribuciones

Este proyecto está en constante evolución y siempre hay espacio para mejorar. Creo que la colaboración es la mejor forma de aprender y crecer juntos.

Si te interesa aprender sobre desarrollo en Go, WebSockets, criptografía o simplemente quieres contribuir a un proyecto open source, te invito a participar. Desde la corrección de errores hasta la implementación de nuevas funcionalidades, tu contribución es bienvenida. 

¡Juntos podemos hacer este proyecto aún mejor!

## Licencia

[MIT License](https://opensource.org/licenses/MIT)