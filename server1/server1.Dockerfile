FROM python:alpine3.16
WORKDIR /
COPY server1.py /server1.py
ENTRYPOINT ["python","/server1.py"]
