from asyncio import Protocol,Transport,create_task
from handlers.Handler import Handler
from Response_States import ResponseStates

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
            print(f"Conexión perdida con {self.addr} debido a un error: {exec}")
        else:
            print(f"Conexión cerrada con {self.addr}")

    async def process_message(self,request):
        try:
            #identify protocol chain
            print("Start process")
            result = await self.handler_chain.handle(request)
            print(f"Result: {result}")
            status = result["status"]
            if ResponseStates.CLOSE_CONNECTION == status:
                self.transport.close()
                print("Disconnect order was permormed with success")
            elif ResponseStates.SENT_RESPONSE == status:
                print("Send message")
                self.transport.write(result["message"])
                
        except Exception as e:
            print(f"Error happen with message: {e}")
        