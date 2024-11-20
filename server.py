import asyncio
import dotenv
from datetime import datetime
from utilities.logs import register_connection
from utilities.extractors import identify_message_type,identify_device_id
from utilities.calculations import calculate_checksum
import os

if os.environ.get("ENV"):
    dotenv.load_dotenv()

HOST = os.environ.get("HOST")
PORT = os.environ.get("PORT")

async def client_connection(reader,writer):
    serial_number = 0
    while True:
        data = await reader.read(1024)
        addr = writer.get_extra_info("peername")
        message = ''.join([hex(byte) for byte in data])
        serial_number += 1
        try:
            await register_connection(f"{addr}: {message} : {datetime.now()}")
            message_type = identify_message_type(data)
            device_id = identify_device_id(data)
            if message_type == "0100":
                message_length = b"\x02"
                response = b"\x7e\x81\x00"+message_length+device_id
                response += bytes([serial_number])
                response += b"\x00\x00" #result
                response += calculate_checksum(response)
                response += b"\x7e"
                writer.write(response)
            elif message_type == "0003":
                break
        except FileNotFoundError:
            pass

    writer.close()
    await writer.wait_closed()

async def main():
    server = await asyncio.start_server(client_connection,HOST,PORT)
    
    async with server:
        await server.serve_forever()
        
if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("Bye...")
