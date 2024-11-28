class Handler:
    def __init__(self,next_handler=None):
        self.next_handler = next_handler

    async def handle(self, request):
        if self.next_handler:
            return await self.next_handler.handle(request)
        return None
        