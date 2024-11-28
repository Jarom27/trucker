from handlers.Handler import Handler
from protocols import IGPSProtocol

class ManageHandler(Handler):
    def __init__(self, next_handler=None):
        self.protocol = None
        super().__init__(next_handler)
        
    async def handle(self,request):
        self.protocol = request["protocol"]
        print("Managing protocol")
        print(self.protocol)
        response = self.protocol.process_message(request["data"])
        return response