
from handlers.ProtocolHandler import ProtocolHandler
from handlers.LogHandler import LogHandler
from handlers.ValidateHandler import ValidateHandler
from handlers.ProtocolSelector import ProtocolSelector
from handlers.ManageHandler import ManageHandler
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
    manage_handler = ManageHandler()
    log_handler = LogHandler(manage_handler)
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