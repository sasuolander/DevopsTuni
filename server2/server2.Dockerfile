FROM python:alpine3.16
WORKDIR /
COPY server2.py /server2.py
ENTRYPOINT ["python","/server2.py"]
