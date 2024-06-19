import jwt
import os

import aiohttp
from typing import Union

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from cent import AsyncClient, Client, PublishRequest


app = FastAPI()

origins = [
    "http://localhost:3000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# SERVER API CLIENT https://github.com/centrifugal/pycent
# just a client https://github.com/centrifugal/centrifuge-python


class User(BaseModel):
    user_id: str

@app.post("/centrifugo/subscribe/")
def subscribe(user: User):
    hmac_secret = os.getenv("CENTRIFUGO_HMAC_SECRET")

    # NOTE sub is necessarily and stands for user_id/uid
    # another necessarily claim is for centrifugo hmac secret
    # can be provided with exp(iration) time
    # https://centrifugal.dev/docs/server/authentication

    encoded = jwt.encode({"sub": user.user_id}, hmac_secret, algorithm="HS256")

    print('\n POST /centrifugo/subscribe/ >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>')
    print(f"user_id: {user.user_id}")
    print(f"token: {encoded}")
    print('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n')

    return {"token": encoded}


# python client example btw
# https://github.com/centrifugal/centrifuge-python/blob/master/example.py

@app.post("/centrifugo/send-hello/")
async def send_hello():
    centrifugo_api_key = os.getenv("CENTRIFUGO_API_KEY")
    headers = {
        "Content-Type": "application/json",
        "X-API-Key": centrifugo_api_key,
    }
    payload = {
        "channel": "space",
        "data": {
            "hello": "from backend",
        }
    }

    async with aiohttp.ClientSession() as session:
        async with session.post("http://centrifugo:8000/api/publish", headers=headers, json=payload) as login_resp:
            print("Status login:", login_resp.status)
            print(await login_resp.json())

    return {"sent": "ok"}


class Message(BaseModel):
    payload: dict


@app.post("/centrifugo/send/")
async def send_payload(message: Message):
    centrifugo_api_key = os.getenv("CENTRIFUGO_API_KEY")
    client = AsyncClient("http://centrifugo:8000/api", centrifugo_api_key)
    request = PublishRequest(channel="space", data=message.payload)
    result = await client.publish(request)
    return {"sent": "ok"}


@app.post("/centrifugo/rpc/")
async def rpc_callback():
    # TODO centrifugo rpc callback route impl (not implemented now, backend streams service used for now)
    print('\n callback >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>')
    return {"status": "ok"}


@app.post("/from-streams/send/")
async def send_from_streams(message: Message):
    print('\n POST /from-streams/send/ >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>')
    print(f"payload from streams: {message.payload}")
    print('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n')

    if message.payload.get("x"):
        message.payload["x"] = message.payload["x"] * 2

    if message.payload.get("y"):
        message.payload["y"] = message.payload["y"] * 2

    centrifugo_api_key = os.getenv("CENTRIFUGO_API_KEY")
    client = AsyncClient("http://centrifugo:8000/api", centrifugo_api_key)
    request = PublishRequest(channel="space", data=message.payload)
    result = await client.publish(request)
    return {"sent": "ok"}
