FROM golang:latest

WORKDIR /app

# Installe Air correctement
RUN go install github.com/air-verse/air@latest

# Ajoute le dossier des binaires Go au PATH
ENV PATH="/go/bin:${PATH}"

CMD ["air"]