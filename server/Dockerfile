FROM python:3

RUN pip install flask
RUN pip install pynacl

ADD . .
CMD ["python", "./server.py"]