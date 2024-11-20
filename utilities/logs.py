import aiofiles

async def register_connection(message):
    async with aiofiles.open("logs/connections.txt","a") as f:
        await f.write(f"\n{message}")

async def location_log(message):
    async with aiofiles.open("logs/locations.txt","a") as f:
        await f.write(f"\n{message}")