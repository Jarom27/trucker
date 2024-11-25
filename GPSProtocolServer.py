from asyncio import Protocol,Transport,create_task
from Handler import Handler

class GPSProtocolServer(Protocol):
    def __init__(self, handler_chain: Handler):
        self.transport = None
        self.handler_chain = handler_chain
        self.addr = None
    def connection_made(self,transport: Transport):
        self.transport = transport
        self.addr = transport.get_extra_info("peername")
        print(f"Connecting with {self.addr}")
    
    def data_received(self, data:bytes):
        print(f"Receive data from {self.addr}: {data.hex()}")
        message = {"ip": self.addr[0], "data": data}  # Simular identificación del dispositivo
        create_task(self.process_message(message))
        

    def connection_lost(self,exec):
        if exec:
            print(f"Conexión perdida con {self.addr} debido a un error: {exc}")
        else:
            print(f"Conexión cerrada con {self.addr}")

    async def process_message(self,request):
        try:
            result = await self.handler_chain.handle(request)
            protocol = request["protocol"]
            print("Processing data")
            protocol.send_response = self.transport.write
            response = protocol.process_message(request["data"])

            if response == -1:
                self.transport.close()
                
        except Exception as e:
            print(f"Error happen with message: {e}")
        