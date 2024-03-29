# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.17.1 as builder

ENV GO111MODULE on
ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/ensetservice
WORKDIR /src/ensetservice

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates
# RUN apk add --no-cache postgresql-client

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Google API
ENV GOOGLE_KEY=yourgoogleAPI.apps.googleusercontent.com
ENV GOOGLE_SECRET=yourgoogleAPI

ENV SESSION_SECRET=Dont_Forget_To_Generate_SESSION_SECRET_FOR_OUR_APPLICATION_v1
ENV DATABASE_URL=postgres://postgres:janah@pg_db:5432/ENSETService?sslmode=disable

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# COPY ./wait-for-postgres.sh . 
# Uncomment to run the migrations before running the binary:
# CMD /bin/wait-for-postgres.sh; /bin/app migrate; /bin/app
CMD /bin/app migrate; /bin/app
# CMD exec /bin/app
