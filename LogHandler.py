import aiofiles
from Handler import Handler
class LogHandler(Handler):
    async def handle(self, request):
        print(f"Registrando mensaje procesado: {request['response']}")
        # Guardar en logs o base de datos
        async with aiofiles.open("logs/connections.txt","a") as f:
            await f.write(f"\n{request}")
        return "Procesamiento completo"