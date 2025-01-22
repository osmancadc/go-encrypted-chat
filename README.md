# Chat Encriptado de Extremo a Extremo (E2E) en Golang

<!-- ![Logo](aquí_va_la_ruta_de_la_imagen.png) -->
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

*   `internal/websocket`: Maneja la comunicación WebSocket.
    *   `client.go`: Gestiona las conexiones WebSocket individuales.
    *   `message.go`: Define las estructuras de los mensajes y payloads.
    * `websocket.go`: contiene la configuración del servidor websocket.
*   `internal/chat`: Contiene la lógica del chat.
    *   `user.go`: Define la estructura `User`.
*   `internal/config`: Gestiona la configuración de la aplicación, incluyendo las claves criptográficas.
* `pkg/crypto`: Contiene la implementacion del cifrado RSA y AES.

![Diagrama de la arquitectura (opcional)](aquí_va_la_ruta_del_diagrama.png) ## Flujo de Funcionamiento (Ejemplo: Envío de un Mensaje)

1.  El usuario A escribe un mensaje.
2.  El cliente de A obtiene la clave simétrica del grupo desde la configuración.
3.  El cliente de A cifra el mensaje con la clave simétrica.
4.  El cliente de A envía el mensaje cifrado a través de WebSockets.
5.  El cliente del usuario B recibe el mensaje cifrado.
6.  El cliente de B obtiene la clave simétrica del grupo desde su configuración.
7.  El cliente de B descifra el mensaje.
8.  El mensaje descifrado se muestra al usuario B.

## Manejo de Claves

*   **Claves RSA (Asimétricas):**
    *   Se genera un par de claves RSA (pública y privada) para cada usuario al iniciar la aplicación.
    *   La clave *privada* se guarda de forma segura en la configuración del usuario.
    *   La clave *pública* se comparte con otros usuarios para el intercambio de claves simétricas.
*   **Claves AES (Simétricas):**
    *   Cada grupo de chat tiene una clave AES única.
    *   La clave AES se genera al crear el grupo y se distribuye de forma segura entre los miembros del grupo utilizando cifrado RSA.

## Cómo ejecutar la aplicación

1.  Clona el repositorio: `git clone <URL_del_repositorio>`
2.  Navega al directorio del proyecto: `cd <nombre_del_proyecto>`
3.  Construye la aplicación: `go build`
4.  Ejecuta la aplicación: `./<nombre_del_ejecutable>` (necesitarás ejecutar múltiples instancias en diferentes terminales para simular varios usuarios).

## Próximas mejoras (Opcional)

*   Persistencia de mensajes y claves (cifrado en disco).
*   Interfaz gráfica de usuario (GUI).
*   Manejo de errores más robusto.
*   Pruebas unitarias completas.

## Contribuciones

Las contribuciones son bienvenidas. Si encuentras un error o tienes una sugerencia, por favor abre un issue o un pull request.

## Licencia

[Añade la licencia que corresponda (por ejemplo, MIT)](aquí_va_el_enlace_a_la_licencia)