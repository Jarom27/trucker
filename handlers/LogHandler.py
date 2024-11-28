import aiofiles
from handlers.Handler import Handler
class LogHandler(Handler):
    def __init__(self,next_handler=None):
        self.next_handler = next_handler

    async def handle(self, request):
        print(f"Registrando mensaje procesado: {request}")
        # Guardar en logs o base de datos
        async with aiofiles.open("logs/connections.txt","a") as f:
            await f.write(f"\n{request}")
        return await super().handle(request)