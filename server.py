import asyncio
from datetime import datetime
from utilities.logs import register_connection

HOST = "localhost"
PORT = 7700

async def client_connection(reader,writer):
    data = await reader.read(1024)
    addr = writer.get_extra_info("peername")
    message = ''.join([hex(byte) for byte in data])
    try:
        await register_connection(f"{addr}: {message} : {datetime.now()}")
    except FileNotFoundError:
        pass


    writer.close()
    await writer.wait_closed()

async def main():
    server = await asyncio.start_server(client_connection,HOST,PORT)
    
    async with server:
        await server.serve_forever()
        
if __name__ == "__main__":
    asyncio.run(main())
