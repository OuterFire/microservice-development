FROM python:3-slim

WORKDIR /app

COPY /tests/ .

# Install python requirements
RUN pip install --no-cache-dir -r python_library/requirements.txt
