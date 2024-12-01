
from handlers.ProtocolHandler import ProtocolHandler
from handlers.LogHandler import LogHandler
from handlers.ValidateHandler import ValidateHandler
from handlers.ProtocolSelector import ProtocolSelector
from handlers.ManageHandler import ManageHandler
from handlers.DbHandler import DbHandler
from GPSProtocolServer import GPSProtocolServer
from database.Database import create_pool
import asyncio
import dotenv
import os

dotenv.load_dotenv()

async def main():
    loop = asyncio.get_running_loop()
    HOST = os.environ.get("HOST")
    PORT = os.environ.get("PORT")
    
    #Conexion a la base de datos
    pool = await create_pool()

    # Crear la cadena de responsabilidad
    protocol_selector = ProtocolSelector()
    db_handler = DbHandler(pool)
    manage_handler = ManageHandler(db_handler)
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