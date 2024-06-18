FROM python:3.12.4

WORKDIR /app

COPY ./backend/requirements.txt /app/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /app/requirements.txt

CMD ["fastapi", "run", "main.py", "--port", "8081", "--reload"]
