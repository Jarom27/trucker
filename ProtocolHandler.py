import Handler
from ProtocolSelector import ProtocolSelector
class ProtocolHandler(Handler.Handler):
    def __init__(self,selector:ProtocolSelector, next_handler=None):
        super().__init__(next_handler)
        self.selector = selector
        
    async def handle(self,request):
        print("Selecting protocol")
        protocol = self.selector.get_protocol(request["data"])
        response = protocol.process_message(request)
        request["response"] = response
        request["protocol"] = protocol
        
        return await super().handle(request)