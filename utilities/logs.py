import aiofiles

async def register_connection(message):
    async with aiofiles.open("logs/connections.txt","a") as f:
        await f.write(message)
