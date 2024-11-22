
from ProtocolHandler import ProtocolHandler
from LogHandler import LogHandler
from ValidateHandler import ValidateHandler
from ProtocolSelector import ProtocolSelector
from GPSProtocolServer import GPSProtocolServer
import asyncio
import dotenv
import os

dotenv.load_dotenv()

async def main():
    loop = asyncio.get_running_loop()
    HOST = os.environ.get("HOST")
    PORT = os.environ.get("PORT")
    # Crear la cadena de responsabilidad
    protocol_selector = ProtocolSelector()
    log_handler = LogHandler()
    protocol_handler = ProtocolHandler(protocol_selector, log_handler)
    validate_handler = ValidateHandler(protocol_handler)
    
    # Crear servidor TCP con la cadena de responsabilidad
    server = await loop.create_server(lambda: GPSProtocolServer(validate_handler), HOST, PORT)
    print(f"Servidor corriendo en {HOST}:{PORT}")

    try:
        await server.serve_forever()
    except asyncio.CancelledError:
        print("Servidor detenido.")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("Servidor detenido.")