import os
import asyncio
import yaml
from hypercorn.config import Config
from hypercorn.asyncio import serve
from fastapi import FastAPI

def GetPathByte(Folderpath) :
    size = 0

    for path, dirs, files in os.walk(Folderpath):
        try :
            for f in files:
                fp = os.path.join(path, f)
                size += os.path.getsize(fp)
        except exception as e :
            pass
    
    return { "used" : size }


def main() :

    with open('config/config.yaml', 'r') as f:
        APIconfig = yaml.load(f , Loader=yaml.FullLoader)

    app = FastAPI()

    @app.get("/")
    async def GetPath_Byte() :
        return GetPathByte(Folderpath=APIconfig["mounted_path"])

    config = Config()
    config.bind = [APIconfig["api_url"]]
    asyncio.run(serve(app=app, config=config))

if __name__ == "__main__" :
    main()
