from handlers.Handler import Handler

class ValidateHandler(Handler):
    async def handle(self, request):
        print("Validando datos...")
        if len(request["data"]) == 0:
            raise ValueError("Mensaje vacío")
        return await super().handle(request)