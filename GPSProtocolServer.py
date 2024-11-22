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
            message_type = request["message_type"]
            print(f"Result: {result}, {message_type}")
            protocol = request["protocol"]
            print("Executing commands")
            protocol.execute_message(request, self.transport)
                
        except Exception as e:
            print(f"Error happen with message: {e}")
        